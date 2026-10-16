## ADDED Requirements

### Requirement: GitOps-managed deployment
The system SHALL be deployable through ArgoCD from this git repository, such that the cluster state converges to the ROOT kustomization of the tracked branch.

#### Scenario: Application watches the repository
- **WHEN** the ArgoCD Application `td27` is applied to the `argocd` namespace
- **THEN** it SHALL point at `https://github.com/pddzl/td27-admin-micro.git`, path `.` (root kustomization), destination namespace `td27` on the in-cluster cluster

#### Scenario: Automated sync
- **WHEN** a commit lands on the tracked branch
- **THEN** ArgoCD SHALL sync automatically without manual `kubectl apply`
- **AND** resources removed from git SHALL be pruned from the cluster (prune: true)
- **AND** out-of-band cluster changes SHALL be reverted to the git state (selfHeal: true)

#### Scenario: Sync retry on transient failure
- **WHEN** a sync fails (e.g., image not yet present in the registry)
- **THEN** ArgoCD SHALL retry with a bounded backoff instead of leaving the sync permanently failed

### Requirement: Migration runs as a sync hook
The `td27-migrations` Job SHALL run as an ArgoCD Sync hook so it can be recreated on every sync without hitting Job immutability.

#### Scenario: Hook annotations present
- **WHEN** `deploy/k8s/migrate-job.yaml` is rendered
- **THEN** the Job SHALL carry `argocd.argoproj.io/hook: Sync` and `argocd.argoproj.io/hook-delete-policy: BeforeHookCreation`

#### Scenario: Re-sync recreates the Job
- **WHEN** a new sync occurs after a previous sync created the hook Job
- **THEN** ArgoCD SHALL delete the previous hook Job before creating a new one (no `field is immutable` failure)

#### Scenario: Manual deployment unaffected
- **WHEN** the stack is deployed with plain `kubectl apply -k .`
- **THEN** the hook annotations SHALL be ignored and the migration Job SHALL run exactly as before

### Requirement: Deployments constrained by AppProject
An AppProject `td27` SHALL restrict where and from what ArgoCD may deploy for this application.

#### Scenario: Source restriction
- **WHEN** an Application in project `td27` declares a source repository other than `https://github.com/pddzl/td27-admin-micro.git`
- **THEN** ArgoCD SHALL reject the sync

#### Scenario: Destination restriction
- **WHEN** an Application in project `td27` declares a destination outside namespace `td27` on the in-cluster server
- **THEN** ArgoCD SHALL reject the sync
