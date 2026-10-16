## Tasks

- [x] 1. Draft proposal/design/spec delta and validate the change
- [x] 2. Rewrite `deploy/k8s/datastores.yaml`: drop StatefulSets/PVCs/peer Service, add selector-less Services + manual Endpoints (placeholder IP)
- [x] 3. Update `AGENTS.md` k8s topology and first-sync cleanup note
- [x] 4. Verify: `openspec validate`, `kubectl kustomize .` renders without datastore workloads, YAML parses
