## Why

Deployment is currently manual: build images, `kubectl apply -k .`. There is no drift correction, no deployment history, and no single place that says "the cluster should look like this commit". ArgoCD (GitOps) makes the git repo the source of truth: the cluster continuously converges to what is committed, every change is a commit, and rollbacks are `git revert`.

## What Changes

- Add `deploy/argocd/` with:
  - `project.yaml` — AppProject `td27` restricting sources to this repo and destinations to the `td27` namespace on the in-cluster server.
  - `application.yaml` — ArgoCD Application `td27` watching `https://github.com/pddzl/td27-admin-micro.git` at path `.` (the ROOT kustomization), `automated` sync with `prune` + `selfHeal`, and a sync retry policy.
- Convert `deploy/k8s/migrate-job.yaml` into an ArgoCD **Sync hook** (`argocd.argoproj.io/hook: Sync`, `hook-delete-policy: BeforeHookCreation`) so the immutable Job can be recreated on every sync instead of failing with "field is immutable".
- Document bootstrap (install ArgoCD, register private-repo credentials) in file header comments. No docs file.

## Impact

- Affected: `deploy/argocd/` (new), `deploy/k8s/migrate-job.yaml` (annotations only).
- No application code, config, or schema changes.
- Plain `kubectl apply -k .` keeps working unchanged (hook annotations are inert outside ArgoCD).
- After approval of this proposal, the cluster converges to whatever is on the tracked branch; committing to the branch becomes a deployment.

## Non-goals

- Installing ArgoCD itself (use the upstream install manifest into the `argocd` namespace; documented in bootstrap comments).
- Secret management (committed dev-value Secrets stay as-is; External Secrets/Sealed Secrets is a future change).
- Automated image tag bumping (argocd-image-updater or a CI step is a future change).
- Multi-environment overlays (dev/staging/prod) — single `td27` namespace for now.
