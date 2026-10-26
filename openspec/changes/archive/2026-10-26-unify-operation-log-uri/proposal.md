# Unify Operation-Log URIs on /operation-log/*

## Why

The operation-log endpoints are registered twice on `main`
(`api/gateway/internal/handler/handler.go`): the canonical
`/operation-log/list|delete|delete-by-ids|cleanup` routes and legacy
`/opl/list|delete|delete-by-ids` aliases pointing at the same handlers. The
frontend still calls the `/opl/*` aliases. Two names for one resource is
inconsistent with every other module (users, roles, dicts all use their full
resource name), doubles the permission surface, and per repo convention dead
endpoints are deleted rather than kept as aliases.

The frontend is the only consumer of `/opl/*` (verified: no other references
in code or tests), so a clean rename is safe.

## What Changes

- **Gateway** (`handler.go`): delete the three `/opl/*` alias routes; keep
  the four `/operation-log/*` routes unchanged (middleware wiring identical).
- **Frontend** (`web/src/api/sysMonitor/operationLog.ts`): switch the three
  calls to `/operation-log/list`, `/operation-log/delete`,
  `/operation-log/delete-by-ids`.
- **Migration** (`scripts/011_unify_operation_log_uri.sql`, idempotent):
  rename `/opl/*` rows to `/operation-log/*` in `sys_management_api.path`
  and `sys_management_permission.resource` (domain `'api'`), following the
  `004_align_api_permissions.sql` pattern. Before writing it, audit the live
  rows: if `/operation-log/*` rows already coexist with `/opl/*` rows, the
  migration must also drop the now-duplicate `/opl/*` rows instead of
  creating duplicates.
- After applying: policy reload (`POST /permission/reload-policy`) or
  gateway+basis restart, since permission rows changed.

## Impact

- Endpoints `/opl/list|delete|delete-by-ids` disappear — breaking only for
  the bundled frontend, which is updated in the same change.
- URL suffix action mapping is unchanged (`/list` → read, `/delete` →
  delete, `/delete-by-ids` → delete, `/cleanup` → default update, same as
  today).
- Verification: `go vet`, `go build`, frontend grep for `/opl/` returns
  nothing outside historical docs; live check that the operation-log page
  lists rows, delete works, and no 403/404 appears (permission rows intact).
