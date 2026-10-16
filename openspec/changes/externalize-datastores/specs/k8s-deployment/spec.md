## ADDED Requirements

### Requirement: External datastore dependencies
PostgreSQL and etcd SHALL run on the host, outside Kubernetes; the manifests SHALL expose them to in-cluster workloads only as selector-less Services backed by manual Endpoints in `deploy/k8s/datastores.yaml`.

#### Scenario: Stable in-cluster DNS names
- **WHEN** basis, gateway, or the migrate Job resolve `postgres:5432` or `etcd-client:2379`
- **THEN** traffic SHALL reach the host-local services through the Service + Endpoints shims, with no change to any pod spec or service config

#### Scenario: No datastore workloads in the manifest
- **WHEN** the root kustomization is rendered
- **THEN** it SHALL NOT contain any StatefulSet, PersistentVolumeClaim, or volumeClaimTemplate for postgres or etcd

#### Scenario: Single per-cluster edit point
- **WHEN** the host IP changes
- **THEN** only the two Endpoints `ip` values in `deploy/k8s/datastores.yaml` need editing, in git (ArgoCD selfHeal applies them; cluster edits are reverted)
