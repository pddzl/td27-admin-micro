## Decisions

### D1 — etcd also moves to the host (not just postgres)
- **Chosen**: host-local etcd alongside host-local pg.
- **Rejected: etcd stays in k8s** — etcd-as-StatefulSet is an anti-pattern (bootstrap coupling, stable-storage demands), no HA gain on a single node, and it would keep the service registry inside a layer we want to be able to wipe and re-sync freely.

### D2 — selector-less Service + manual Endpoints as the shim
- **Chosen**: keep Service names `postgres` / `etcd-client`, strip selectors, add manual Endpoints with the host IP.
- **Rejected: hostAliases on pod specs** — must be duplicated into basis, gateway, and migrate pod templates, loses Service semantics (no clusterIP, no port abstraction).
- **Rejected: ExternalName** — requires a DNS name for the host (bare LAN IPs unsupported), adds a DNS dependency for a single-node setup.
- **Rejected: EndpointSlice** — the modern API but more verbose for one static address; Endpoints remain fully supported.

### D3 — host IP lives in git (RFC-5737 placeholder)
- **Chosen**: the two Endpoints `ip:` values in `datastores.yaml`, defaulted to `203.0.113.10` (TEST-NET-3, documentation-only).
- **Rationale**: kustomize/kubectl have no env interpolation, and `selfHeal: true` means the cluster is reverted to git — so git must be the single source of truth. TEST-NET-3 guarantees that an unedited file fails fast instead of silently connecting to a real address.

### D4 — no pod spec or config changes
Service names and ports are unchanged, so `pg_isready -h postgres` (migrate Job), `etcdctl --endpoints=http://etcd-client:2379` (basis init container), and the `db-list` / etcd entries in `basis.yaml`/`gateway.yaml` need no edits.

## Non-goals

- No HA story for pg/etcd (single-node, dev-grade — same as before).
- No migration of existing in-cluster pg data to the host instance (user handles data move if any).
