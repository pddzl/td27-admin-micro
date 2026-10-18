# Tasks

## 1. Dockerfiles

- [x] 1.1 `deploy/docker/Dockerfile.basis` + `Dockerfile.gateway`: add `ARG GO_IMAGE/RUN_IMAGE/APK_MIRROR/GOPROXY` with CN defaults (mirror prefix, goproxy.cn, aliyun apk sed).
- [x] 1.2 Create `deploy/docker/Dockerfile.web` (moved from `web/Dockerfile`) with `ARG NODE_IMAGE/RUN_IMAGE/NPM_REGISTRY` mirror args; delete `web/Dockerfile`.

## 2. Manifests

- [x] 2.1 `deploy/k8s/kustomization.yaml`: images → `localhost/td27-{basis,gateway,web}`, `newTag` placeholder `v0.1.0`.
- [x] 2.2 `deploy/k8s/{basis,gateway,web}.yaml`: `imagePullPolicy: Never`.
- [x] 2.3 Remove `deploy/k8s/app-configs.yaml` from resources; remove `td27-db-credentials` Secret from `datastores.yaml`; delete `app-configs.yaml`.
- [x] 2.4 `deploy/argocd/application.yaml`: `prune: false`, `ignoreDifferences` restartedAt, bootstrap comments → `argocd repo add` SSH style + `make secret` note.

## 3. Makefile + config files

- [x] 3.1 Root `Makefile`: build/run/clean, image-basis/image-gateway/image-web, image-import (CTR override), release (tag-guard + kustomization bump + commit + push), secret, reconfig, restart, app.
- [x] 3.2 Create `deploy/k8s-configs/{db.env,basis.yaml,gateway.yaml}` seeded from current committed values (k8s hostnames preserved); add `deploy/k8s-configs/` to `.gitignore`.

## 4. Verification

- [x] 4.1 `kubectl kustomize .` renders 15 resources (18 minus 3 Secrets) with localhost images + Never policy.
- [x] 4.2 `make -n release TAG=v0.1.0` dry-run shows the correct command sequence (tag-guard → 3 builds → export/import → sed bump → commit → push).
- [x] 4.3 `git status` shows `deploy/k8s-configs/` ignored (verified via `git check-ignore -v`).
- [x] 4.4 `openspec validate add-makefile-release-workflow` passes.

## 5. Documentation

- [x] 5.1 AGENTS.md: Kubernetes deployment section rewritten for the Makefile flow (node-local images, make release, out-of-band secrets).
