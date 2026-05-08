# Register & Build Remaining Services

## Problem
The remaining 4 services (Cron, Cache, ServiceToken, OperationLog) have implementations already created in the source tree, but they need:
1. Registration in `basis.go` entry point
2. Proto-generated type fix (files are at wrong path)
3. `go mod tidy` and `go build` verification

---

## Step 1: Fix Proto-Generated File Locations

The proto generation placed files at `rpc/basis/types/td27/rpc/basis/types/...` but they should be at `rpc/basis/types/...` directly.

Run:
```bash
cd /Users/paul/Mine/td27-admin-micro

# Create correct directories
mkdir -p rpc/basis/types/tool/cron
mkdir -p rpc/basis/types/tool/cache
mkdir -p rpc/basis/types/tool/service_token
mkdir -p rpc/basis/types/monitor/operation_log

# Copy (not symlink - use cp) the generated files to correct location
cp rpc/basis/types/td27/rpc/basis/types/tool/cron/*.go rpc/basis/types/tool/cron/
cp rpc/basis/types/td27/rpc/basis/types/tool/cache/*.go rpc/basis/types/tool/cache/
cp rpc/basis/types/td27/rpc/basis/types/tool/service_token/*.go rpc/basis/types/tool/service_token/
cp rpc/basis/types/td27/rpc/basis/types/monitor/operation_log/*.go rpc/basis/types/monitor/operation_log/

# Remove the wrong nested tree
rm -rf rpc/basis/types/td27
```

---

## Step 2: Register Remaining Services in basis.go

**File**: `rpc/basis/basis.go`

**Add to imports**:
```go
"td27/rpc/basis/internal/server/tool"
"td27/rpc/basis/internal/server/monitor"
"td27/rpc/basis/types/tool/cron"
"td27/rpc/basis/types/tool/cache"
"td27/rpc/basis/types/tool/service_token"
"td27/rpc/basis/types/monitor/operation_log"
```

(The `tool` server import is already there from file service — ensure it's present)

**Update gRPC registration block** (add these within the func):
```go
s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
    user_pb.RegisterUserServer(grpcServer, sysManagement.NewUserServer(ctx))
    role_pb.RegisterRoleServer(grpcServer, sysManagement.NewRoleServer(ctx))
    permission_pb.RegisterPermissionServer(grpcServer, sysManagement.NewPermissionServer(ctx))
    menu_pb.RegisterMenuServer(grpcServer, sysManagement.NewMenuServer(ctx))
    dept_pb.RegisterDeptServer(grpcServer, sysManagement.NewDeptServer(ctx))
    dict_pb.RegisterDictServer(grpcServer, sysManagement.NewDictServer(ctx))
    api_pb.RegisterAPIServer(grpcServer, sysManagement.NewAPIServer(ctx))
    button_pb.RegisterButtonServer(grpcServer, sysManagement.NewButtonServer(ctx))
    file_pb.RegisterFileServer(grpcServer, tool.NewFileServer(ctx))
    cron.RegisterCronServer(grpcServer, tool.NewCronServer(ctx))
    cache.RegisterCacheServer(grpcServer, tool.NewCacheServer(ctx))
    service_token.RegisterServiceTokenServer(grpcServer, tool.NewServiceTokenServer(ctx))
    operation_log.RegisterOperationLogServer(grpcServer, monitor.NewOperationLogServer(ctx))

    if c.Mode == service.DevMode || c.Mode == service.TestMode {
        reflection.Register(grpcServer)
    }
})
```

---

## Step 3: Build & Verify

```bash
# Ensure required dependencies are present (check go.sum for common package)
cd /Users/paul/Mine/td27-admin-micro
go mod tidy
go build rpc/basis/basis.go
```

Address any compilation errors (likely missing imports or type mismatches in the new logic/server files).

---

## Step 4: Update Task Tracking

In `openspec/changes/rewrite-basis-service-reference-pattern/tasks.md`:
- [x] 6.4 Implement gRPC server handlers (all services done)
- [x] 6.5 Update service entry point (all registered)
- [ ] 7.1 Run `go mod tidy`
- [ ] 7.2 Run `go build rpc/basis/basis.go`
- [ ] 7.3 Test service startup
- [ ] 7.4 Verify etcd discovery
- [ ] 7.5 Clean up backup

---

## Files that need manual creation (if sub-agents didn't write them)

The sub-agents analyzed and wrote detailed code recommendations for these files. If the files don't exist at these paths, they need to be created with the code from the agent outputs:

- `rpc/basis/internal/logic/tool/cron.go`
- `rpc/basis/internal/server/tool/cron.go`
- `rpc/basis/internal/logic/tool/cache.go`
- `rpc/basis/internal/server/tool/cache.go`
- `rpc/basis/internal/logic/tool/service_token.go`
- `rpc/basis/internal/server/tool/service_token.go`
- `rpc/basis/internal/server/monitor/` (mkdir)
- `rpc/basis/internal/logic/monitor/operation_log.go`
- `rpc/basis/internal/server/monitor/operation_log.go`

Each follows the exact same pattern as `logic/tool/file.go` and `server/tool/file.go` respectively.
