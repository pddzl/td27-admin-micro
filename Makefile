# td27-admin-micro — release workflow (mirrors ~/TrinaSolar/vmware-alarm-aggregation)
#
#   make release TAG=v0.1.0   # full release: import pre-built images, bump
#                             # tags, commit + push -> ArgoCD syncs
#   make secret               # one-time: create the out-of-band Secrets
#   make app                  # one-time: register the ArgoCD Application
#
# Run on the k8s node (or wherever BUILDER/CTR below reach it). Images are
# node-local (localhost/*, imagePullPolicy: Never) — there is no registry.
# No Dockerfiles live in this repo (deploy/docker/ removed): build the three
# images yourself, tagged $(IMG_BASIS)/$(IMG_GATEWAY)/$(IMG_WEB):$(TAG),
# BEFORE running make image-import (it fails with "no such image" otherwise).

GO      ?= go
GOCACHE ?= /tmp/td27-gocache

# Node-local images; kubelet never pulls them from a registry.
IMG_BASIS   ?= localhost/td27-basis
IMG_GATEWAY ?= localhost/td27-gateway
IMG_WEB     ?= localhost/td27-web
TAG         ?= v0.1.0

# Image builder CLI: podman (daemonless) or docker — both work identically.
# Either way the image is invisible to kubelet until imported, so the
# save -> ctr import step in image-import is always required.
BUILDER   ?= podman
# Command to load images into the k8s node's container runtime. Examples:
#   containerd (kubeadm): make image-import
#   k3s:                  make image-import CTR="sudo k3s ctr"
#   RKE2:                 make image-import CTR="sudo /var/lib/rancher/rke2/bin/ctr -n k8s.io --address /run/k3s/containerd/containerd.sock"
CTR       ?= sudo ctr -n k8s.io
ARGOCD_NS ?= argocd
BRANCH    ?= main
WORKLOAD_NS ?= td27

IMAGES = $(IMG_BASIS) $(IMG_GATEWAY) $(IMG_WEB)

.PHONY: build run clean image-import release tag-guard secret reconfig restart app

# --- Local development ------------------------------------------------------

build:
	GOMODCACHE=$(GOCACHE) $(GO) build -o /dev/null ./rpc/basis
	GOMODCACHE=$(GOCACHE) $(GO) build -o /dev/null ./api/gateway

run:
	GOMODCACHE=$(GOCACHE) $(GO) run rpc/basis/basis.go -f rpc/basis/etc/basis.yaml

clean:
	rm -f basis gateway

# --- Images: import pre-built images into the node's containerd --------------
#
# No in-repo Dockerfiles: build the three node-local images yourself first
# (any builder, any context), tagged:
#   localhost/td27-basis:$(TAG)  localhost/td27-gateway:$(TAG)  localhost/td27-web:$(TAG)

# Load the images into the k8s node's container runtime so pods can use them
# with imagePullPolicy: Never.
image-import:
	@for img in $(IMAGES); do \
		tar="/tmp/$$(echo $$img | tr '/' '-').$(TAG).tar"; \
		echo "==> exporting $$img:$(TAG) -> $$tar"; \
		$(BUILDER) save -o "$$tar" $$img:$(TAG); \
		echo "==> importing into containerd"; \
		$(CTR) images import "$$tar"; \
	done

# Full code release in ONE command — bundles the only manual steps of the
# GitOps flow, then lets ArgoCD do the rest automatically:
#
#   make release TAG=v0.1.1
#
# TAG is required and must not be 'latest': an unchanged tag means no manifest
# diff, so ArgoCD would sync nothing and pods would keep the old binary.
release: tag-guard image-import
	sed -i 's/newTag: .*/newTag: $(TAG)/' deploy/k8s/kustomization.yaml
	git add deploy/k8s/kustomization.yaml
	git commit -m "release $(TAG)"
	git push origin $(BRANCH)
	@echo "Pushed $(TAG). ArgoCD will sync automatically, or force it now: argocd app sync td27"

tag-guard:
	@test -n "$(TAG)" || { echo 'ERROR: make release TAG=vX.Y.Z (a real version tag)'; exit 1; }
	@test "$(TAG)" != "latest" || { echo 'ERROR: make release TAG=vX.Y.Z (a real version tag)'; exit 1; }

# --- Out-of-band Secrets (NOT managed by ArgoCD — never committed) -----------
#
# Holds plaintext credentials: postgres password, JWT signing key. The git
# sources live in deploy/k8s-configs/ (gitignored). ArgoCD does not declare
# these Secrets, so with prune: false it never touches them.
secret:
	kubectl create namespace $(WORKLOAD_NS) --dry-run=client -o yaml | kubectl apply -f -
	kubectl -n $(WORKLOAD_NS) create secret generic td27-db-credentials \
		--from-env-file=deploy/k8s-configs/db.env --dry-run=client -o yaml | kubectl apply -f -
	kubectl -n $(WORKLOAD_NS) create secret generic basis-config \
		--from-file=basis.yaml=deploy/k8s-configs/basis.yaml --dry-run=client -o yaml | kubectl apply -f -
	kubectl -n $(WORKLOAD_NS) create secret generic gateway-config \
		--from-file=gateway.yaml=deploy/k8s-configs/gateway.yaml --dry-run=client -o yaml | kubectl apply -f -

# Config-only change: NO image rebuild — configs live in Secrets but the
# services read them once at startup, so restart after `make secret`.
reconfig: secret
	kubectl rollout restart deployment/basis deployment/gateway -n $(WORKLOAD_NS)
	kubectl rollout status deployment/basis -n $(WORKLOAD_NS)
	kubectl rollout status deployment/gateway -n $(WORKLOAD_NS)

restart:
	kubectl rollout restart deployment/basis deployment/gateway deployment/web -n $(WORKLOAD_NS)
	kubectl rollout status deployment/basis -n $(WORKLOAD_NS)
	kubectl rollout status deployment/gateway -n $(WORKLOAD_NS)

# One-time: register the Application with ArgoCD (after `argocd repo add`;
# see deploy/argocd/application.yaml header for the full bootstrap).
app:
	kubectl apply -n $(ARGOCD_NS) -f deploy/argocd/project.yaml -f deploy/argocd/application.yaml
