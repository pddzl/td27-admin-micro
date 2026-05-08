## Context

The gateway had no operation audit trail. Admin operations (create user, delete role, etc.) were invisible.

## Goals / Non-Goals

**Goals:**
- Record operation logs for all mutation HTTP requests
- Capture: IP, method, path, status, user agent, params, response, timing, user info
- Non-blocking: log recording must not affect response time

**Non-Goals:**
- No changes to read-only endpoints (GET)
- No UI or reporting

## Decisions

1. gRPC CreateOperationLog RPC to communicate between gateway and service
2. Async goroutine for log writing (fire-and-forget)
3. Wraps only POST/PUT/DELETE routes (not GET)
4. Reuses existing OperationLogModel and repository

## Risks / Trade-offs

- Async goroutine may lose logs if process crashes before write completes
- Acceptable for audit logging (non-critical path)
