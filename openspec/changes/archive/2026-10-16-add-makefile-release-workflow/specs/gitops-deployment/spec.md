## MODIFIED Requirements

### Requirement: GitOps-managed deployment
The system SHALL be deployable through ArgoCD from this git repository, such that the cluster state converges to the ROOT kustomization of the tracked branch.

#### Scenario: Application watches the repository
- **WHEN** the ArgoCD Application `td27` is applied to the `argocd` namespace
- **THEN** it SHALL point at `git@github.com:pddzl/td27-admin-micro.git`, path `.` (root kustomization), destination namespace `td27` on the in-cluster cluster

#### Scenario: Automated sync
- **WHEN** a commit lands on the tracked branch
- **THEN** ArgoCD SHALL sync automatically without manual `kubectl apply`
- **AND** resources deleted from git SHALL NOT be pruned from the cluster (prune: false — manual cleanup is deliberate)
- **AND** out-of-band cluster changes SHALL be reverted to the git state (selfHeal: true)

#### Scenario: Rollout restart does not create drift
- **WHEN** `make restart` stamps the `kubectl.kubernetes.io/restartedAt` annotation into a Deployment pod template
- **THEN** ArgoCD SHALL ignore that annotation and stay InSync

#### Scenario: Sync retry on transient failure
- **WHEN** a sync fails (e.g., image not yet present on the node)
- **THEN** ArgoCD SHALL retry with a bounded backoff instead of leaving the sync permanently failed
