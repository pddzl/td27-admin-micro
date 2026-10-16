## Context

The repo already ships a complete Kustomize stack (root `kustomization.yaml`: namespace, postgres/etcd StatefulSets, config Secrets, migrate Job, basis/gateway/web Deployments, Ingress). Deployments are manual. ArgoCD will pull from git instead of pushing.

## Decisions

### D1: Application targets the ROOT kustomization (path `.`)

Kustomize forbids configMapGenerator file references above the kustomization root, so the `td27-migrations` ConfigMap (built from `scripts/0*.sql`) can only be generated at the repo root. An Application pointing at `deploy/k8s` would reference a ConfigMap that is never rendered, and the migrate Job would fail. Targeting `.` gives ArgoCD exactly what `kubectl apply -k .` renders today (verified: 18 resources).

### D2: Migrate Job becomes an ArgoCD Sync hook

Kubernetes Jobs are immutable: re-applying an existing Job whose spec changed (any future edit to the Job template) hard-fails sync with `field is immutable`, and ArgoCD's prune/recreate cycles can trigger the same error. The idiomatic GitOps pattern is:

```yaml
annotations:
  argocd.argoproj.io/hook: Sync
  argocd.argoproj.io/hook-delete-policy: BeforeHookCreation
```

- The Job is no longer a tracked live resource; ArgoCD creates it during sync and deletes it at the start of the next sync (BeforeHookCreation), so recreation never hits immutability.
- Consequence: migrations re-run on every sync. This is already a documented convention ("scripts must be safe to re-apply").
- The annotations are inert under plain `kubectl apply -k .`, so the manual flow is unchanged.
- Alternative rejected: `syncOptions: [Replace=true]` on the Application — replaces ALL resources (delete+create), causing downtime for StatefulSets and Services on every drift correction.

### D3: No sync waves

Wave annotations would order resource application, but the stack already self-sequences via initContainers: the migrate Job waits for `pg_isready`, basis waits for etcd health, gateway readiness gates its Service, and ArgoCD health checks use the same probes. Waves would duplicate this and add coupling to ArgoCD specifics for no gain.

### D4: AppProject `td27` as a guardrail

Restricts `sourceRepos` to this GitHub repo, `destinations` to namespace `td27` on the in-cluster server, and cluster-scoped resource kinds to what the stack needs. Prevents a stray Application from deploying elsewhere.

### D5: Repo credentials

The repo is private (SSH-only for the owner). ArgoCD needs read access: documented path is an `argocd`-namespace `repo` secret (HTTPS token or deploy key) referenced by `application.yaml`'s source... kept OUT of the Application spec by using ArgoCD's standard credential resolution (repo URL match), so `application.yaml` itself stays credential-free and public-safe.

### D6: Automated sync policy

`automated: { prune: true, selfHeal: true }` + bounded retry. Prune keeps git authoritative (deleting a manifest from git deletes it from the cluster); selfHeal corrects out-of-band `kubectl` edits — matching the GitOps contract that the repo is the single source of truth.

## Risks / Trade-offs

- Auto-prune means an accidental deletion merged to the tracked branch is a deletion in the cluster. Mitigation: the tracked branch is `main`; direct pushes stay the existing team convention.
- Migration re-runs per sync rely on script idempotency (already enforced by convention; `001–003` are guarded/historical).
- Committed Secrets with dev values become cluster-pulling material — unchanged exposure from today's manual flow; External Secrets deferred (non-goal).
