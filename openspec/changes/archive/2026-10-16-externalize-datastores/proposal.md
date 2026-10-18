## Why

PostgreSQL will run on the host, outside Kubernetes (user decision, single-node cluster). Keeping etcd in-cluster while pg is host-local splits the stateful layer across two management domains for no benefit: etcd-as-StatefulSet is a recognized anti-pattern, buys zero HA on one node, and couples the service registry to cluster rebuilds. Moving etcd to the host as well makes everything ArgoCD manages stateless and rebuildable, and gives one dependency story and one backup surface (pg dumps + etcd snapshot).

## What Changes

- Remove the postgres + etcd StatefulSets, their volumeClaimTemplates, and the etcd peer Service from `deploy/k8s/datastores.yaml`.
- Keep Services `postgres:5432` and `etcd-client:2379` — now selector-less, backed by manual Endpoints pointing at the host IP (RFC-5737 documentation IP `203.0.113.10` as the single per-cluster edit point).
- No changes to basis/gateway/web/migrate manifests or service configs: same DNS names, same ports (`pg_isready -h postgres`, `etcdctl --endpoints=http://etcd-client:2379` keep working).
- Update the k8s topology in `AGENTS.md`.

## Impact

- Manifests only; no code or config changes.
- User must run postgres + etcd on the host (podman/systemd) and set the real host IP in `datastores.yaml` before first apply.
- One-time cleanup after first ArgoCD sync: `prune: false` leaves the previously deployed in-cluster StatefulSets running — delete them manually (`kubectl -n td27 delete sts postgres etcd`) and reclaim their PVCs.
