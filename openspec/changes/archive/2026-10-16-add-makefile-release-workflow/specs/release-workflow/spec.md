## ADDED Requirements

### Requirement: Node-local image release flow
The system SHALL be releasable to the k8s node without a registry: a Makefile SHALL build each image with podman/docker, import it into the node's containerd, bump the image tags in the kustomization, and push so ArgoCD syncs.

#### Scenario: Release with a version tag
- **WHEN** the user runs `make release TAG=v1.2.3` on the node
- **THEN** all three images (`localhost/td27-basis`, `localhost/td27-gateway`, `localhost/td27-web`) SHALL be built and imported into the node's container runtime
- **AND** the `newTag` values in `deploy/k8s/kustomization.yaml` SHALL be bumped to `v1.2.3`, committed, and pushed

#### Scenario: Tag guard
- **WHEN** `make release` is run without TAG or with `TAG=latest`
- **THEN** the target SHALL fail with an error instead of producing a no-op release

#### Scenario: Cluster converges to the released tag
- **WHEN** the tag bump commit lands on the tracked branch
- **THEN** ArgoCD SHALL sync and the pods SHALL roll onto the new images with `imagePullPolicy: Never`

### Requirement: CN-friendly image builds
Image builds SHALL default to CN-accessible mirrors (registry mirror prefix, apk mirror, Go module proxy, npm registry) and SHALL allow disabling them via Make variables for unrestricted networks.

#### Scenario: Default build on a CN network
- **WHEN** `make image-import` runs with default variables
- **THEN** base images SHALL be pulled via the mirror prefix and Go modules/npm packages via CN proxies

#### Scenario: Direct build on an unrestricted network
- **WHEN** `make image-import MIRROR= APK_MIRROR=` runs
- **THEN** builds SHALL use upstream registries without mirror prefixes

### Requirement: Secrets are out-of-band
The `td27-db-credentials`, `basis-config`, and `gateway-config` Secrets SHALL NOT be committed to git or managed by ArgoCD; they SHALL be created by `make secret` from gitignored files.

#### Scenario: Secret creation
- **WHEN** `make secret` runs
- **THEN** the three Secrets SHALL be created/updated in namespace `td27` from `deploy/k8s-configs/` files, which are gitignored

#### Scenario: ArgoCD ignores out-of-band secrets
- **WHEN** ArgoCD syncs the application
- **THEN** the three Secrets SHALL NOT be pruned or flagged as drift (prune: false; ArgoCD never declares them)
