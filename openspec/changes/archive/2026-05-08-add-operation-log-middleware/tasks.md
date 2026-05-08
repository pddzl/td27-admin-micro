## 1. gRPC: Add CreateOperationLog RPC

- [x] 1.1 Update operation_log.proto with CreateOperationLogReq
- [x] 1.2 Regenerate proto types
- [x] 1.3 Implement logic CreateOperationLog
- [x] 1.4 Add server handler
- [x] 1.5 Register in basis.go

## 2. Gateway: Add OperationRecord middleware

- [x] 2.1 Create middlewate/operation_record.go
- [x] 2.2 Wire on mutation routes in handler.go
- [x] 2.3 Build verification
