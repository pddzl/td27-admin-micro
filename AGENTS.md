# AGENTS.md — td27-admin-micro

Compact guidance for AI agents. Only repo-specific, non-obvious facts.

## Quick Overview

- Go 1.26.4 monorepo using go-zero v1.9.4, migrated from Gin.
- Two processes: `rpc/basis` (gRPC, port 8080) and `api/gateway` (HTTP REST, port 8888).
- Module: `td27`.
- DB: **PostgreSQL** via **sqlx** (NOT GORM — was migrated from GORM). Uses `github.com/jmoiron/sqlx` with `pgx` driver.
- Auth: JWT (HMAC-SHA256, DB-persisted token blocklist), Casbin RBAC backed by the unified `sys_management_permission` table (read-only adapter, allow/deny `effect` column, deny wins), bcrypt cost 12 (transparent legacy-MD5 rehash on login), JIT role elevations, button permissions (`<page>:<action>` codes).
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

### Database migrations

- SQL migrations live in `scripts/` and must be applied in filename order with `psql` (no auto-migration):
  - `001_add_permission_effect.sql` — `effect` column on `sys_management_permission` (allow/deny)
  - `002_role_elevations.sql` — JIT elevation table
  - `003_buttons.sql` — button table + 21 seeded `<page>:<action>` codes + root grants
- Example roles for RBAC testing: `scripts/rbac_seed_example.sql` (not auto-applied).

### Verification

```bash
GONOSUMCHECK='*' go vet ./...
GONOSUMCHECK='*' go test ./...
```

Go unit tests exist for gateway middleware (authz, rate limiting), password hashing, Casbin models, and rpc services. Frontend tests live in `web/tests` (vitest, run from `web/`).

### Kubernetes deployment

- Entry point is the ROOT `kustomization.yaml`: run `kubectl apply -k .` from the repo root. `deploy/k8s/kustomization.yaml` alone omits the `td27-migrations` ConfigMap the migrate Job needs — kustomize forbids file references above its root, so the generator must live at the root.
- Build images first (build context noted per file): `deploy/docker/Dockerfile.basis` and `Dockerfile.gateway` use the repo root as context; the web image uses `web/Dockerfile` with `web/` as context. Default tags `ghcr.io/pddzl/td27-admin-micro/{basis,gateway,web}:latest`; retarget via the `images:` block in `deploy/k8s/kustomization.yaml`.
- Topology: ingress(nginx) → web:8500 → nginx proxies `/api/*` to gateway:8888 → gRPC basis:8080 → postgres/etcd StatefulSets. The gateway is intentionally not exposed; all public traffic enters through the web nginx.
- go-zero conf has no env interpolation: service configs ship as whole-file Secrets (`basis-config`, `gateway-config`) rewritten for k8s hostnames (`postgres`, `etcd-client`) with `Mode: prod`. For production, create these out-of-band instead of committing real values.
- The `td27-migrations` Job applies `scripts/0*.sql` in filename order via `psql -v ON_ERROR_STOP=1`. Jobs are immutable — after changing a script: `kubectl -n td27 delete job td27-migrations --ignore-not-found && kubectl apply -k .`, and scripts must be safe to re-apply.

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
- Public endpoints (health, captcha, login) have no per-route middleware.
- Global middlewares are attached via `server.Use` in `RegisterHandlers` (handler.go), gated by `RateLimit.Enabled`: per-IP token buckets (429 + `Retry-After`, stricter bucket on `/login`, `/captcha`, `/logout`) and a `MaxConcurrent` in-flight throttle (503). All routes are covered; `/health` is exempt. Buckets are in-memory (per instance) and trust `X-Forwarded-For` — only trust that header behind a proxy that sets it.

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
3. **Casbin adapter is read-only** — policies load from `sys_management_permission` (domains: menu/api/button/data; deny takes precedence). Adapter write ops return `errWriteUnsupported`; policy changes go through the Permission service/DB, never the enforcer. Authorization fails closed when the permission RPC is unavailable.
4. **Login requires captcha** — `POST /captcha` returns id + image, then `POST /login` with `captcha_id` + `captcha`.
5. **Config uses `mapstructure` tags** — YAML keys must match (hyphenated: `signing-key`, `db-name`).
6. **Proto import paths** use `sysManagement`, `sysMonitor`, `sysTool` (camelCase), NOT `monitor` or `tool`.
7. **Log encoding defaults to JSON** — add `Encoding: plain` in YAML for readable terminal output.
8. **Gateway blocks until basis.rpc appears in etcd** — `zrpc.MustNewClient` panics if not found.
9. **Permission checks resolve roles live from the DB** — gateway passes `user_id` to `CheckPermission`, which includes active JIT elevations; grants/revocations apply without re-login. Password changes revoke the caller's token via the DB-persisted blocklist (re-login required).
10. **Self-service endpoints bypass RBAC** — password modify and elevation request/my routes are exempt from the `update`/`create` permission check so ordinary users can't lock themselves out; `ModifyPassword` overrides `req.Id` from JWT claims (never trust body-supplied ids).
