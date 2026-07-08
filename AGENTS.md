# AGENTS.md — td27-admin-micro

Compact guidance for AI agents. Only repo-specific, non-obvious facts.

## Quick Overview

- Go 1.26.4 monorepo using go-zero v1.9.4, migrated from Gin.
- Two processes: `rpc/basis` (gRPC, port 8080) and `api/gateway` (HTTP REST, port 8888).
- Module: `td27`.
- DB: **PostgreSQL** via **sqlx** (NOT GORM — was migrated from GORM). Uses `github.com/jmoiron/sqlx` with `pgx` driver.
- Auth: JWT (HMAC-SHA256), Casbin RBAC (dummy adapter — NOT yet persisted to DB), bcrypt.
- Business modules: `sysManagement` (users, roles, permissions, menus, depts, dicts, APIs, buttons), `sysMonitor` (operation logs, dashboard), `sysTool` (files, cron, cache, service tokens).

## Architecture

```
Client → HTTP → api/gateway (8888) → gRPC via etcd basis.rpc → rpc/basis (8080) → PostgreSQL
```

Layers within `rpc/basis`: `server → logic → service → repository → model`. All repos use raw SQL/sqlx.

## Key Commands

Run from repository root. Build requires env vars due to mod cache permissions:

```bash
GONOSUMCHECK='*' GONOSUMDB='*' GOFLAGS='-mod=mod' go build ./rpc/basis/...
```

### Startup order

```bash
# 1. PostgreSQL + etcd must be running first
# 2. Start gRPC service
GONOSUMCHECK='*' go run rpc/basis/basis.go -f rpc/basis/etc/basis.yaml

# 3. Start HTTP gateway (depends on basis registering in etcd)
GONOSUMCHECK='*' go run api/gateway/gateway.go -f api/gateway/etc/gateway.yaml
```

### Proto generation

```bash
protoc --go_out=. --go-grpc_out=. --go_opt=module=td27 --go-grpc_opt=module=td27 \
  -I ./rpc/basis/proto ./rpc/basis/proto/<module>/<file>.proto
```

Proto `go_package` convention: `td27/rpc/basis/types/<module>/<name>_pb;<name>_pb`.

**Required after proto changes**: regenerate → `go mod tidy` → restart both processes.

### Verification

```bash
GONOSUMCHECK='*' go vet ./...
```

No tests exist yet.

## Repository Layer (sqlx)

All 15 repository files in `rpc/basis/internal/repository/` use raw SQL via sqlx.
- `GetContext` / `SelectContext` for queries.
- `NamedExecContext` for inserts/updates.
- `ExecContext` for deletes (soft delete: `UPDATE ... SET deleted_at=NOW()`).
- Every SELECT must include `AND deleted_at IS NULL` for soft-delete filtering.
- Use `COALESCE(created_at, NOW())` in SELECT to handle NULL timestamps from legacy data.

## HTTP Handler Patterns

All gateway handlers in `api/gateway/internal/handler/basis/` follow:
- `pkg/api.DecodeAndValidate(r.Body, &req)` for request parsing + validation.
- `pkg/api.FailWithRequest`, `api.FailWithMessage`, `api.OkWithData`, `api.OkWithDetailed` for responses.
- Inline `var req struct{ ... validate:"required" ... }` for request shapes.
- Mutations use `opRecordMiddleware.Handle(jwtMiddleware.Handle(handler))`.
- Reads use `jwtMiddleware.Handle(handler)` only.
- Public endpoints (health, captcha, login) have no middleware.

## Registered gRPC Services (16)

| Module | Services |
|--------|----------|
| basis | Ping |
| sysManagement | User, Role, Permission, Menu, Dept, Dict, DictDetail, API, Button |
| sysMonitor | OperationLog, Dashboard |
| sysTool | File, Cron, Cache, ServiceToken |

Non-obvious: Dict and DictDetail are **separate** services with separate protos, servers, and handlers.

## Important Gotchas

1. **Build requires** `GONOSUMCHECK='*' GONOSUMDB='*' GOFLAGS='-mod=mod'` — Go module cache is owned by root.
2. **DB is sqlx, NOT GORM** — original GORM code was removed entirely. No auto-migration.
3. **Casbin adapter is a dummy** — `// TODO: Implement PermissionAdapter`. RBAC policies not persisted.
4. **Login requires captcha** — `GET /api/captcha` → `POST /login` with `captcha_id` + `captcha`.
5. **Config uses `mapstructure` tags** — YAML keys must match (hyphenated: `signing-key`, `db-name`).
6. **Proto import paths** use `sysManagement`, `sysMonitor`, `sysTool` (camelCase), NOT `monitor` or `tool`.
7. **Log encoding defaults to JSON** — add `Encoding: plain` in YAML for readable terminal output.
8. **Gateway blocks until basis.rpc appears in etcd** — `zrpc.MustNewClient` panics if not found.
