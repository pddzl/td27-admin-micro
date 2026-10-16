# Tasks

## 1. ArgoCD manifests

- [x] 1.1 Create `deploy/argocd/project.yaml`: AppProject `td27` in the `argocd` namespace — sourceRepos `https://github.com/pddzl/td27-admin-micro.git`, destinations limited to namespace `td27` on server `https://kubernetes.default.svc`, namespace-scoped resource kinds only.
- [x] 1.2 Create `deploy/argocd/application.yaml`: Application `td27` in the `argocd` namespace — project `td27`, source repo `https://github.com/pddzl/td27-admin-micro.git` path `.` (root kustomization), destination namespace `td27` in-cluster, `syncPolicy.automated` with `prune: true`, `selfHeal: true`, and a bounded retry backoff.
- [x] 1.3 Add bootstrap instructions as header comments in `application.yaml`: install ArgoCD from the upstream manifest into `argocd`, apply `project.yaml` + `application.yaml`, and register private-repo credentials (deploy key or token) via an `argocd`-namespace `repo` secret (kept out of the Application spec).

## 2. Migrate Job → ArgoCD Sync hook

- [x] 2.1 Add `argocd.argoproj.io/hook: Sync` and `argocd.argoproj.io/hook-delete-policy: BeforeHookCreation` annotations to `deploy/k8s/migrate-job.yaml`; update its header comment to explain hook behavior under ArgoCD (recreated per sync, inert under plain kubectl).

## 3. Verification

- [x] 3.1 `kubectl kustomize .` renders unchanged except the new Job annotations (18 resources; hook annotations present).
- [x] 3.2 `openspec validate add-argocd-gitops-deployment` passes.
- [x] 3.3 `GONOSUMCHECK='*' go vet ./...` and `GONOSUMCHECK='*' go test ./...` still pass (no Go changes expected).
- [x] 3.4 Confirm Application YAML parses as valid Kubernetes resources (`kubectl apply --dry-run=client` or `kubeconform` if cluster is unavailable). — parsed with YAML safe-load (AppProject/td27, Application/td27); no cluster available for server-side validation.

## 4. Documentation

- [x] 4.1 Update AGENTS.md "Kubernetes deployment" section: add the GitOps apply path (install ArgoCD → apply project/application → auto-sync) and the hook behavior note for the migrate Job.
