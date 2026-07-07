# Finish Remaining Services Implementation

## TL;DR
> Finish implementing the remaining gRPC service handlers (Cron, Cache, ServiceToken, OperationLog), register them in `basis.go`, and verify the build compiles.

## Remaining Work

Implement 4 remaining service logic handlers + 4 server handlers + registration in main entry point.

## Tasks

### Task 1: Cron Service Implementation
**Files to create:**
- `rpc/basis/internal/logic/tool/cron.go` - CronLogic with full RPC methods (GetCron, ListCron, CreateCron, UpdateCron, ToggleStatus, ExecuteNow, DeleteCron)
- `rpc/basis/internal/server/tool/cron.go` - CronServer handler

### Task 2: Cache Service Implementation
**Files to create:**
- `rpc/basis/internal/logic/tool/cache.go` - CacheLogic with RPC methods (GetCache, SetCache, DeleteCache, ListCache, CleanupExpired)
- `rpc/basis/internal/server/tool/cache.go` - CacheServer handler

### Task 3: ServiceToken Service Implementation
**Files to create:**
- `rpc/basis/internal/logic/tool/service_token.go` - ServiceTokenLogic with RPC methods (Create, Get, List, Update, ToggleStatus, AssignPermissions, GetPermissions, Delete, Validate)
- `rpc/basis/internal/server/tool/service_token.go` - ServiceTokenServer handler

### Task 4: OperationLog Service Implementation
**Files to create:**
- `rpc/basis/internal/logic/monitor/operation_log.go` - OperationLogLogic (List, CleanupExpired)
- `rpc/basis/internal/server/monitor/operation_log.go` - OperationLogServer handler
- `rpc/basis/internal/server/monitor/` directory

### Task 5: Register All Services in Entry Point
**File to edit:**
- `rpc/basis/basis.go` - Add imports for `cron_pb`, `cache_pb`, `service_token_pb`, `operation_log_pb`, `toolserver`, `monitorserver` packages and register all 4 servers

