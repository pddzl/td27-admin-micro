# release-workflow Specification

## Purpose
Define how td27 images move from source to the k8s node: images are built OUT of the repo per the contract in DEPLOY.md, imported node-locally into containerd (no registry), released via tag-guarded Makefile commits that ArgoCD applies, with Secrets always out-of-band.

## Requirements

### Requirement: Node-local image release flow
The system SHALL be releasable to the k8s node without a registry: the three `localhost/td27-{basis,gateway,web}` images SHALL be built by the operator (no Dockerfiles committed), imported into the node's containerd by the Makefile, have their tags bumped in the kustomization, and be pushed so ArgoCD syncs.

#### Scenario: Release with a version tag
- **WHEN** the user runs `make release TAG=v1.2.3` on the node with the three images already built and tagged `v1.2.3`
- **THEN** all three images SHALL be exported from the local builder store and imported into the node's container runtime
- **AND** the `newTag` values in `deploy/k8s/kustomization.yaml` SHALL be bumped to `v1.2.3`, committed, and pushed

#### Scenario: Tag guard
- **WHEN** `make release` is run without TAG or with `TAG=latest`
- **THEN** the target SHALL fail with an error instead of producing a no-op release

#### Scenario: Cluster converges to the released tag
- **WHEN** the tag bump commit lands on the tracked branch
- **THEN** ArgoCD SHALL sync and the pods SHALL roll onto the new images with `imagePullPolicy: Never`

### Requirement: Out-of-repo image builds
Image construction SHALL happen outside the repository (no committed Dockerfiles); the Makefile `image-import` target SHALL save and import pre-built `localhost/td27-*:$(TAG)` images, failing with "no such image" if any is missing.

#### Scenario: Import of pre-built images
- **WHEN** `make image-import TAG=v1.2.3` runs and the three images exist in the local builder store
- **THEN** each image SHALL be exported to a tar and imported into the node's container runtime via `CTR` (overridable for kubeadm/k3s/RKE2)

#### Scenario: Missing image fails fast
- **WHEN** an image with the requested tag does not exist in the local builder store
- **THEN** the import SHALL fail with a "no such image" error instead of silently releasing a stale tag

### Requirement: Secrets are out-of-band
The `td27-db-credentials`, `basis-config`, and `gateway-config` Secrets SHALL NOT be committed to git or managed by ArgoCD; they SHALL be created by `make secret` from gitignored files.

#### Scenario: Secret creation
- **WHEN** `make secret` runs
- **THEN** the three Secrets SHALL be created/updated in namespace `td27` from `deploy/k8s-configs/` files, which are gitignored

#### Scenario: ArgoCD ignores out-of-band secrets
- **WHEN** ArgoCD syncs the application
- **THEN** the three Secrets SHALL NOT be pruned or flagged as drift (prune: false; ArgoCD never declares them)
