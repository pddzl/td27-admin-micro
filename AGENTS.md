# AGENTS.md - td27-admin-micro
Compact guidance for AI agents working in this repository. Only includes repo-specific, non-obvious facts.

## Quick Overview
- Go 1.25.5 monorepo using go-zero v1.9.4
- Two processes: `rpc/basis` (gRPC server) and `api/gateway` (HTTP REST gateway)
- Module name: `td27`
- Tech stack: gRPC, go-zero rest, GORM, PostgreSQL, Casbin RBAC, JWT auth, bcrypt
- Three business modules: `sysManagement` (authority/ACL), `sysMonitor` (audit), `sysTool` (system tools)

## Project Structure
```
/
├── rpc/basis/            # gRPC service (port 8080, etcd key: basis.rpc)
│   ├── basis.go          # Entry point
│   ├── etc/basis.yaml    # Config
│   ├── proto/            # Protobuf definitions
│   │   ├── common/
│   │   ├── sysManagement/
│   │   ├── sysMonitor/
│   │   └── sysTool/
│   ├── types/            # Generated proto Go types (packages use _pb suffix)
│   └── internal/
│       ├── config/       # Go config struct
│       ├── svc/          # ServiceContext (dependency injection)
│       ├── service/      # Business logic layer
│       ├── repository/   # DB operations (GORM)
│       ├── logic/        # gRPC method implementations
│       ├── server/       # gRPC server handlers
│       ├── model/        # GORM model definitions
│       │   ├── sysManagement/ sysMonitor/ sysTool/ common/
│       └── initialization/  # DB connection setup
├── api/gateway/          # HTTP REST gateway (port 8888)
│   ├── gateway.go        # Entry point
│   ├── etc/gateway.yaml  # Config
│   └── internal/
│       ├── config/       # RestConf extension
│       ├── svc/          # gRPC client connections (via etcd)
│       ├── handler/      # HTTP handlers (calls gRPC clients)
│       ├── middleware/   # JWT auth + operation record
│       └── types/        # JSON response format
└── pkg/tool/             # Shared utilities (md5, time helpers)
```

## Architecture
```
Client → HTTP → api/gateway (port 8888) → gRPC via etcd basis.rpc → rpc/basis (port 8080) → PostgreSQL
```

All layers within `rpc/basis` follow: `server → logic → service → repository → model`.

## Key Commands
All commands run from repository root:

### Build/Run
```bash
# gRPC service
go run rpc/basis/basis.go -f rpc/basis/etc/basis.yaml

# HTTP gateway (depends on gRPC service via etcd)
go run api/gateway/gateway.go -f api/gateway/etc/gateway.yaml
```

### Protobuf Generation
Proto `go_package` follows the convention `td27/rpc/basis/types/<module>/<name>_pb;<name>_pb`.

```bash
# Generate a single proto file
protoc --go_out=. --go-grpc_out=. --go_opt=module=td27 --go-grpc_opt=module=td27 \
  -I ./rpc/basis/proto \
  ./rpc/basis/proto/<module>/<file>.proto

# Generate all protos at once
for f in rpc/basis/proto/common/common.proto \
  rpc/basis/proto/sysManagement/*.proto \
  rpc/basis/proto/sysMonitor/*.proto \
  rpc/basis/proto/sysTool/*.proto; do
  protoc --go_out=. --go-grpc_out=. --go_opt=module=td27 --go-grpc_opt=module=td27 \
    -I ./rpc/basis/proto "$f"
done
```

**Required after modifying proto files**: Regenerate + `go mod tidy` + restart.

### Testing
```bash
go test ./... -v                # all tests
go test ./rpc/basis/... -v     # gRPC service tests only
go fmt ./...
go vet ./...
```

### Proto Service List (14 total)
| Module | Services | Methods |
|--------|----------|---------|
| basis | Ping | 1 |
| sysManagement | User, Role, Permission, Menu, Dept, Dict, API, Button | 8 |
| sysMonitor | OperationLog | 1 |
| sysTool | File, Cron, Cache, ServiceToken | 4 |

## Framework Conventions & Gotchas

### go-zero Patterns
- `rpc/basis` uses `zrpc.MustNewServer(c.RpcServerConf, ...)` — auto-registers with etcd when `Etcd:` block is in config
- `api/gateway` uses `rest.MustNewServer(c.RestConf, ...)` — HTTP server
- Shared dependencies initialized once in `svc/servicecontext.go`
- Business logic in `service/` layer, gRPC method implementations in `logic/` layer
- `conf.MustLoad()` loads YAML config into struct via `mapstructure` tags

### Important Gotchas
1. **DB**: PostgreSQL via GORM (NOT MySQL). DSN uses `host= user= password= dbname= port=` format.
2. **gRPC errors**: Never compare directly — use `status.FromError(err)` and check gRPC status codes.
3. **gRPC clients**: Reuse instances to avoid goroutine leaks (generated clients create a new connection per call).
4. **etcd keys**: Convention `<service-name>.rpc` — must match between server config and client config exactly.
5. **Proto `_pb` suffix**: All proto packages use the `_pb` suffix convention (e.g., `user_pb`, `cron_pb`). The `go_package` in .proto files sets both the import path and package name.
6. **Proto generation**: Uses `protoc` directly (not goctl). Must pass `--go_opt=module=td27 --go-grpc_opt=module=td27` to strip the module prefix from output paths.
7. **Config field names**: Use `mapstructure` tags — YAML keys must match the `mapstructure` tag value (usually hyphenated: `db-name`, `signing-key`).
8. **HTTP handlers**: Gateway handlers call gRPC clients, not DB directly. Mutation endpoints (POST/PUT/DELETE) go through `operation_record` middleware for audit logging. Private routes use JWT middleware.

### Excluded Features
- No Prometheus, Jaeger, or OpenTelemetry (excluded by design — indirect deps from go-zero are acceptable)
- No HTTP API layer inside `rpc/basis` (pure gRPC)
- No CORS in gRPC service (handled by gateway)

## References
- go-zero docs: https://go-zero.dev
- goctl RPC: https://github.com/zeromicro/go-zero/blob/master/tools/goctl/rpc/README.md
