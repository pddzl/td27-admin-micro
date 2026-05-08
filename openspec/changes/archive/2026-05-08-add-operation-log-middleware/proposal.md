## Why

The API gateway had no middleware to record operation/audit logs for mutation requests. The reference repo has an OperationRecord middleware that captures request params, response data, user info, and timing. This change adds the same capability.

## What Changes

- Add CreateOperationLog RPC to the operation_log gRPC service
- Create OperationRecord middleware in the API gateway
- Wire middleware on all mutation routes (POST/PUT/DELETE)

## Capabilities

### New Capabilities
- operation-log-recording: Audit logging for HTTP mutation operations

### Modified Capabilities
(N/A)

## Impact

- New RPC method on existing OperationLog service
- New middleware in api/basis
- No changes to existing business logic
