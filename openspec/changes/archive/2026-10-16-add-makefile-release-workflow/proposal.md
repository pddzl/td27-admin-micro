## Why

The current manifests assume registry-hosted images (`ghcr.io/pddzl/td27-admin-micro/*`) and fully git-managed config Secrets. The user deploys to a self-managed k8s node with NO registry: images are built manually with podman/docker and imported into the node's containerd (`ctr images import`). The release procedure must mirror the house style from `~/TrinaSolar/vmware-alarm-aggregation`: a Makefile-driven `make release TAG=vX` that builds + imports images, bumps the tag in git, pushes, and lets ArgoCD sync.

## What Changes

- **Root `Makefile`** (new), mirroring the reference project:
  - `build` / `run` / `clean` — local Go build/run (isolated GOCACHE).
  - `image-basis` / `image-gateway` / `image-web` — build the three node-local images (`localhost/td27-{basis,gateway,web}`) with CN-friendly defaults (daocloud mirror prefix, aliyun apk mirror, goproxy.cn, npmmirror) overridable via `MIRROR` / `APK_MIRROR` / `GOPROXY` / `NPM_REGISTRY`.
  - `image-import` — `podman/docker save` → `ctr images import` (CTR overridable for kubeadm/k3s/RKE2).
  - `release` — `tag-guard` + image-import + bump `newTag` for all three images in `deploy/k8s/kustomization.yaml` + commit + push → ArgoCD auto-syncs.
  - `secret` / `reconfig` / `restart` — out-of-band Secrets from `deploy/k8s-configs/` (gitignored) and rollout restart.
- **Dockerfiles**: add mirror build-args (`GO_IMAGE`, `RUN_IMAGE`, `APK_MIRROR`, `GOPROXY`, `NPM_REGISTRY`) to `deploy/docker/Dockerfile.{basis,gateway}` and `web/Dockerfile`-style args for the web image via a new `deploy/docker/Dockerfile.web` (moved from `web/Dockerfile` so all image builds share `deploy/docker/` conventions).
- **Manifests**: image names → `localhost/td27-*` with `imagePullPolicy: Never`; ArgoCD Application switched to `prune: false` + `selfHeal: true` + `ignoreDifferences` for `kubectl rollout restart`'s `restartedAt` annotation; bootstrap comments switched to `argocd repo add` (SSH) style.
- **Secrets out of git**: `deploy/k8s/app-configs.yaml` and the `td27-db-credentials` Secret embedded in `datastores.yaml` are removed from the kustomization. The three Secrets (`td27-db-credentials`, `basis-config`, `gateway-config`) are created out-of-band by `make secret` from `deploy/k8s-configs/` (gitignored, seeded once from the current content). ArgoCD deliberately does not manage them.

## Impact

- New: `Makefile`, `deploy/k8s-configs/` (gitignored), `deploy/docker/Dockerfile.web`.
- Modified: both backend Dockerfiles, `deploy/k8s/{datastores,basis,gateway,web}.yaml`, `deploy/k8s/kustomization.yaml`, `deploy/argocd/{application,project}.yaml` (comments), `.gitignore`, AGENTS.md.
- Deleted from git: `deploy/k8s/app-configs.yaml` (content moves to `deploy/k8s-configs/` + `make secret`).
- No application code changes. `kubectl apply -k .` still works after `make secret` has run once on the cluster.

## Non-goals

- Installing/operating ArgoCD itself, ingress-nginx, postgres, or etcd on the target node (already the cluster's concern; public images postgres:16-alpine / etcd may need the same ctr-import treatment on air-gapped nodes).
- Multi-environment overlays, image auto-bumping (argocd-image-updater).
- Deploying from this machine — the Makefile targets run on the node (or wherever podman + CTR point at it).
