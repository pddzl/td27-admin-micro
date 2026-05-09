## 1. JWT Token Generation (rpc/basis)

- [x] 1.1 Add `github.com/golang-jwt/jwt/v4` dependency to `rpc/basis/go.mod`
- [x] 1.2 Add `GenerateToken(userId uint64, username string, roleIds []uint64, bufferTime int64) (string, error)` method to `JWTManager` in `rpc/basis/internal/svc/servicecontext.go`
- [x] 1.3 Add `ParseToken(tokenStr string) (*jwt.MapClaims, error)` method to `JWTManager`
- [x] 1.4 Update `UserLogic.Login()` to load user roles and call `svcCtx.JWT.GenerateToken()` instead of returning `"placeholder_token"`
- [x] 1.5 Add captcha config struct to rpc/basis config (key-long, img-width, img-height) — or keep captcha config only at gateway
- [x] 1.6 Run `go mod tidy` and verify compilation

## 2. CAPTCHA Implementation (api/gateway)

- [x] 2.1 Add `github.com/mojocn/base64Captcha` dependency to `api/gateway/go.mod`
- [x] 2.2 Add `Captcha` config block to gateway config struct (`key-long`, `img-width`, `img-height`)
- [x] 2.3 Create `api/gateway/internal/handler/basis/sysManagement/captcha.go` with `CaptchaHandler` and `GenerateCaptcha` handler
- [x] 2.4 Create captcha store wrapper (in-memory, concurrent-safe) in gateway
- [x] 2.5 Register `GET /api/captcha` as a public route in `handler.go`
- [x] 2.6 Update `LoginHandler.Login()` to accept and verify `captcha_id` + `captcha_value` from request body before calling gRPC login
- [x] 2.7 Return HTTP 400 if captcha verification fails

## 3. Token Blocklist (api/gateway)

- [x] 3.1 Create `api/gateway/internal/middleware/blocklist.go` with `TokenBlocklist` interface and in-memory implementation
- [x] 3.2 Implement `IsBlocklisted(tokenID string) bool` and `AddToBlocklist(tokenID string, expiresAt time.Time) error`
- [x] 3.3 Implement cleanup goroutine for expired blocklist entries
- [x] 3.4 Initialize blocklist in `svc.ServiceContext` and expose it to middleware

## 4. Logout Endpoint (api/gateway)

- [x] 4.1 Create `api/gateway/internal/handler/basis/sysManagement/logout.go` with `LogoutHandler`
- [x] 4.2 Implement `Logout` handler: extract token from Authorization header, compute SHA-256 hash, add to blocklist
- [x] 4.3 Register `POST /api/logout` as a JWT-protected route (without operation record) in `handler.go`
- [x] 4.4 Verify logout returns HTTP 200 on success, HTTP 401 on missing/invalid token

## 5. Blocklist Integration in JWT Middleware

- [x] 5.1 Update `JwtMiddleware.Handle()` to compute SHA-256 token hash and check blocklist after parsing
- [x] 5.2 Return HTTP 401 with "token has been invalidated" if token is blocklisted

## 6. Common Proto: IdsReq for Batch Operations

- [x] 6.1 Add `IdsReq` message to `rpc/basis/proto/common/common.proto` with `repeated uint64 ids = 1`
- [x] 6.2 Regenerate proto Go types: run protoc for common.proto
- [x] 6.3 Add corresponding `IdsReq` Go type usage in repository/service layers

## 7. API Batch Delete

- [x] 7.1 Add `rpc DeleteByIds(common.IdsReq) returns (common.SuccessResp)` to API proto
- [x] 7.2 Implement `DeleteByIds` in API repository (batch delete query)
- [x] 7.3 Implement `DeleteByIds` in API service layer
- [x] 7.4 Implement `DeleteByIds` in API logic layer
- [x] 7.5 Add `DeleteByIds` server method in API server
- [x] 7.6 Add `POST /api/apis/delete-by-ids` handler in gateway API handler
- [x] 7.7 Register route in `handler.go` with operation record + JWT middleware
- [x] 7.8 Regenerate proto types

## 8. Cron Batch Delete

- [x] 8.1 Add `rpc DeleteByIds(common.IdsReq) returns (common.SuccessResp)` to Cron proto
- [x] 8.2 Implement `DeleteByIds` in cron repository (batch delete)
- [x] 8.3 Implement `DeleteByIds` in cron service (also stop cron schedulers for deleted entries)
- [x] 8.4 Implement `DeleteByIds` in cron logic
- [x] 8.5 Add `DeleteByIds` server method in cron server
- [x] 8.6 Add `POST /api/cron/delete-by-ids` handler in gateway cron handler
- [x] 8.7 Register route in `handler.go` with operation record + JWT middleware
- [x] 8.8 Regenerate proto types

## 9. OperationLog Single & Batch Delete

- [x] 9.1 Add `rpc Delete(common.IdReq) returns (common.SuccessResp)` and `rpc DeleteByIds(common.IdsReq) returns (common.SuccessResp)` to OperationLog proto
- [x] 9.2 Implement `Delete` and `DeleteByIds` in operation log repository
- [x] 9.3 Implement `Delete` and `DeleteByIds` in operation log service
- [x] 9.4 Implement `Delete` and `DeleteByIds` in operation log logic
- [x] 9.5 Add `Delete` and `DeleteByIds` server methods in operation log server
- [x] 9.6 Add `POST /api/operation-log/delete` and `POST /api/operation-log/delete-by-ids` handlers in gateway
- [x] 9.7 Register routes in `handler.go` with operation record + JWT middleware
- [x] 9.8 Regenerate proto types

## 10. DictDetail Flat List

- [x] 10.1 Add `FlatDictDetailsReq` (dict_id) and `FlatDictDetailsResp` (repeated DictDetailResp) to Dict proto
- [x] 10.2 Add `rpc FlatDictDetails(FlatDictDetailsReq) returns (FlatDictDetailsResp)` to Dict service
- [x] 10.3 Implement `FlatDictDetails` in dict detail repository (query by dict_id, no pagination)
- [x] 10.4 Implement `FlatDictDetails` in dict service
- [x] 10.5 Implement `FlatDictDetails` in dict logic
- [x] 10.6 Add `FlatDictDetails` server method in dict server
- [x] 10.7 Add `POST /api/dict/detail/flat` handler in gateway dict handler
- [x] 10.8 Register route in `handler.go` with JWT middleware
- [x] 10.9 Regenerate proto types

## 11. File Download

- [x] 11.1 Add `rpc DownloadFile(common.IdReq) returns (DownloadFileResp)` to File proto with `file_content` (bytes), `file_name`, `mime` fields
- [x] 11.2 Implement `DownloadFile` in file repository (find file record by ID)
- [x] 11.3 Implement `DownloadFile` in file service (read file binary from upload path, validate path safety)
- [x] 11.4 Implement `DownloadFile` in file logic
- [x] 11.5 Add `DownloadFile` server method in file server
- [x] 11.6 Add `GET /api/file/download/:id` handler in gateway file handler: stream binary with proper Content-Type and Content-Disposition headers
- [x] 11.7 Register route in `handler.go` with JWT middleware
- [x] 11.8 Regenerate proto types

## 12. Button BatchCheckPermission

- [x] 12.1 Add `BatchCheckPermissionReq` (repeated string button_codes, repeated uint64 role_ids) and `BatchCheckPermissionResp` (map<string, bool> results) to Button proto
- [x] 12.2 Add `rpc BatchCheckPermission(BatchCheckPermissionReq) returns (BatchCheckPermissionResp)` to Button service
- [x] 12.3 Implement `BatchCheckPermission` in button service (iterate codes, check each against Casbin or DB)
- [x] 12.4 Implement `BatchCheckPermission` in button logic
- [x] 12.5 Add `BatchCheckPermission` server method in button server
- [x] 12.6 Add `POST /api/button/batch-check` handler in gateway button handler
- [x] 12.7 Register route in `handler.go` with JWT middleware
- [x] 12.8 Regenerate proto types

## 13. Dashboard Module: Proto & RPC

- [x] 13.1 Create `rpc/basis/proto/sysMonitor/dashboard.proto`
- [x] 13.2 Generate proto types for dashboard
- [x] 13.3 Create dashboard repository with aggregate query methods
- [x] 13.4 Create dashboard service layer
- [x] 13.5 Create dashboard logic layer
- [x] 13.6 Create dashboard server and register in basis.go
- [x] 13.7 Register dashboard gRPC client in gateway `svc/servicecontext.go`

## 14. Dashboard Module: Gateway Handlers

- [x] 14.1 Create `api/gateway/internal/handler/basis/sysMonitor/dashboard.go` with `DashboardHandler`
- [x] 14.2 Implement `GetStatistics`, `GetRecentOperations`, `GetSystemInfo` handlers
- [x] 14.3 Register routes in `handler.go`

## 15. Route Registration & Final Wiring

- [x] 15.1 Add all new handler initializations in `handler.go`
- [x] 15.2 Wire blocklist into JWT middleware constructor
- [x] 15.3 Verify all routes compile

## 16. Build & Verify

- [x] 16.1 Run `go build ./rpc/basis/...` — verify compilation
- [x] 16.2 Run `go build ./api/gateway/...` — verify compilation
- [x] 16.3 Run `go vet ./...` — verify no issues
- [x] 16.4 Run `go fmt ./...` — verify formatting
- [ ] 16.5 Manual verification: captcha → login with real JWT → protected endpoint → logout → protected endpoint rejected
- [ ] 16.6 Manual verification: batch delete APIs → verify deletion
- [ ] 16.7 Manual verification: file download → verify binary content served
- [ ] 16.8 Manual verification: dashboard endpoints → verify stats, recent ops, system info
