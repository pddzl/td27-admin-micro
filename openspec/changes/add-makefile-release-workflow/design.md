## Context

The k8s manifests already exist (root kustomization + ArgoCD Application). The target is a self-managed node like the one running `vmware-alarm-aggregation`: podman (or docker) present, no registry, containerd CRI, ArgoCD installed, CN network (docker.io blocked, goproxy needed). The reference project solves this with a Makefile + node-local images + out-of-band Secrets; this change ports that procedure to td27's three-image stack.

## Decisions

### D1: Node-local image names + `imagePullPolicy: Never`

Images become `localhost/td27-basis`, `localhost/td27-gateway`, `localhost/td27-web` and pods set `imagePullPolicy: Never` so kubelet never leaves the node. The kustomization `images:` block keeps only renames/tags — no registry. Same-tag redeployment requires a new tag (enforced by tag-guard), matching the reference contract: "an unchanged tag means no manifest diff, so ArgoCD would sync nothing".

### D2: One Makefile, three images

`image-import` loops over the three images (build if missing, save, `ctr import`). `CTR` defaults to `sudo ctr -n k8s.io` with documented overrides (`k3s ctr`, RKE2 path) copied from the reference. `release` bumps all three `newTag:` values in `deploy/k8s/kustomization.yaml` via one `sed` per image line, commits as `release <TAG>`, pushes.

### D3: CN mirrors as defaults, all overridable

Defaults mirror the reference: `MIRROR=docker.m.daocloud.io/library/` prefix for golang/alpine/nginx/node bases, `APK_MIRROR=mirrors.aliyun.com`, `GOPROXY=https://goproxy.cn,direct` inside the build stage, and `NPM_REGISTRY=https://registry.npmmirror.com` for the web image (pnpm `--registry`). Passing `MIRROR=` / `APK_MIRROR=` empties them for unrestricted networks.

### D4: Secrets out of git (house rule)

The reference keeps its credentials Secret uncommitted and ArgoCD-unmanaged. td27 has three: `td27-db-credentials` (postgres), `basis-config`, `gateway-config` (whole-file go-zero configs). They move to `make secret`, rendered from `deploy/k8s-configs/{db.env,basis.yaml,gateway.yaml}` (gitignored, seeded from the current committed content — k8s hostnames preserved). AppProject guardrails don't block this: ArgoCD simply never declares these Secrets, and with `prune: false` out-of-band resources are never pruned.

### D5: ArgoCD policy aligned to house style

`prune: false` (reference default — safer: deleting a manifest from git won't delete cluster resources), `selfHeal: true`, plus `ignoreDifferences` on `/spec/template/metadata/annotations/kubectl.kubernetes.io~1restartedAt` for Deployments so `make restart` doesn't show as drift. Bootstrap comments switch to the `argocd repo add` CLI style (matching repoURL string exactly), replacing the raw-secret instructions.

### D6: `deploy/argocd/` path is safe as-is

The reference puts its Application in `deploy/_meta/` because ArgoCD scans plain-manifest directories. Our Application targets path `.` with a kustomization.yaml, so ArgoCD builds kustomize output only — it never scans `deploy/argocd/`. No move needed; a comment documents this.

### D7: web image moves to `deploy/docker/Dockerfile.web`

Keeps all image builds under `deploy/docker/` with consistent mirror args (node + nginx bases, pnpm registry). `web/Dockerfile` is removed to avoid two diverging definitions; build context stays `web/`.

## Risks / Trade-offs

- `prune: false` means manifests deleted from git linger in the cluster — accepted house trade-off (manual cleanup when actually removing resources).
- `imagePullPolicy: Never` fails pods if the image wasn't imported on THAT node — acceptable for single-node; multi-node would need the registry path later.
- ArgoCD not managing secrets means a lost cluster needs `make secret` re-run before apps can start — documented in bootstrap comments.
- Ingress requires ingress-nginx on the target cluster; if absent, web is reachable via port-forward or a NodePort patch (follow-up if needed).
