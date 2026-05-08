# Add Operation Log Middleware

## What
Add a `CreateOperationLog` RPC to the gRPC service, then build an `OperationRecord` middleware in the gateway that captures request/response data and records audit logs for mutation operations.

## Changes

### gRPC side (rpc/basis)
1. Add `rpc CreateOperationLog(CreateOperationLogReq) returns (SuccessResp)` to `operation_log.proto`
2. Regenerate proto types
3. Implement `CreateOperationLog` in logic, service, server layers
4. Register in basis.go

### Gateway side (api/basis)
5. Create `middleware/operation_record.go` - captures request body, response, timing, user info, calls gRPC CreateOperationLog
6. Wire middleware on mutation routes in `handler/handler.go`

## Tasks

### 1. gRPC: Add CreateOperationLog RPC
- Add to proto/monitor/operation_log.proto
- Regenerate proto types
- Add logic/monitor/operation_log.go CreateOperationLog
- Add to existing OperationLogServer
- Register in basis.go

### 2. Gateway: Add OperationRecord middleware
- Create api/basis/internal/middleware/operation_record.go
- Wire on POST/PUT/DELETE routes in handler.go
