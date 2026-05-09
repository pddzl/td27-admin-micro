## Context

The project was rewritten from a Gin-based architecture (`Mine/td27-admin/server`) to Go-zero (`td27-admin-micro`). A full audit of the original server's router files vs the gateway's `handler.go` reveals missing functionality in several areas. The original had 50+ endpoints across sysManagement, sysMonitor, and sysTool modules. The rewrite covers most CRUD but is missing authentication features, batch operations, file download, and the entire dashboard module.

Key missing items from original:
1. Login uses `"placeholder_token"` (no real JWT generation)
2. No captcha (`POST /captcha` with `base64Captcha`)
3. No logout (`POST /logout` with token cache invalidation)
4. No batch delete for APIs, cron, or operation logs
5. No dict detail flat list endpoint
6. No file download endpoint
7. No button batch permission check
8. No dashboard module (statistics, recent ops, system info)

The original JWT service (`JwtService`) had full multi-login support with cache-backed token management (`AddToken`, `ValidateToken`, `RemoveToken`, `RemoveAllTokens`). This needs to be replicated in the rewrite.

## Goals / Non-Goals

**Goals:**
- Add all missing auth functions (captcha, logout, real JWT generation)
- Add batch delete endpoints for API, cron, operation log
- Add dict detail flat list endpoint
- Add file download endpoint
- Add button batch permission check endpoint
- Add dashboard module with statistics, recent operations, system info
- Follow existing go-zero patterns (proto → server → logic → service → repository)

**Non-Goals:**
- User self-registration (admin-only user creation already exists and is sufficient)
- OAuth2 or SSO
- Token refresh mechanism (can be added later; access token TTL is configurable)
- Prometheus/Jaeger/OpenTelemetry (excluded by design in the rewrite)

## Decisions

### 1. Captcha: base64Captcha, gateway-only
- **Library**: `github.com/mojocn/base64Captcha` (same as original)
- **Architecture**: Handled entirely at gateway layer (not via gRPC), matching the original's approach
- **Store**: In-memory store (matching original's `DefaultMemStore`), swappable to Redis
- **Config**: Add `Captcha` config block with `key-long`, `img-width`, `img-height`

### 2. JWT Token Generation: Add methods to existing JWTManager
- **Chosen**: Add `GenerateToken()` and `ParseToken()` to `JWTManager` in `rpc/basis/internal/svc/servicecontext.go`
- **Claims**: `userId` (uint64), `username` (string), `roleIds` ([]uint64), `bufferTime` (int64), `exp`, `iat`, `iss`
- **Library**: `github.com/golang-jwt/jwt/v4` (already used in gateway, add to rpc/basis)

### 3. Token Blocklist: In-memory + interface for Redis swap
- Same approach as original's `JwtService.RemoveToken()` — SHA-256 hash of token as identifier
- Original also had `ValidateToken()` which checked cache; the rewrite should check blocklist in JWT middleware
- Support multi-login concept: blocklist per-token, not per-user

### 4. Logout: Authenticated POST
- `POST /api/logout` with JWT auth, no operation record
- Token hash added to blocklist; client instructed to discard token

### 5. Batch Delete Endpoints
- Follow the existing delete pattern in the rewrite (POST, operation recorded)
- Reuse existing `common_pb.IdReq` or create `common_pb.IdsReq` for batch
- Proto: `rpc DeleteByIds(IdsReq) returns (SuccessResp)`

### 6. DictDetail Flat List
- The original had `POST dictDetail/flat` returning all details for a dict ID without pagination
- In the rewrite, this maps to adding a `FlatDictDetails` or similar method to the Dict service
- Proto: `rpc FlatDictDetails(FlatDictDetailsReq) returns (FlatDictDetailsResp)`

### 7. File Download
- Original: `GET file/download` returns raw file binary content
- Rewrite currently has only `GET /api/file/:id` returning file metadata
- Need to add a new endpoint `GET /api/file/download/:id` that streams the file content with proper Content-Type and Content-Disposition headers
- Use the stored file path from the database to read and serve the file

### 8. Button BatchCheckPermission
- Original: `POST button/batchCheck` accepts `{buttonCodes: []string}` returns map of `{code: bool}`
- Rewrite: Add `BatchCheckPermission` to the button proto and handler
- Proto: `rpc BatchCheckPermission(BatchCheckReq) returns (BatchCheckResp)`

### 9. Dashboard Module (Entirely New)
- New gRPC service in `sysMonitor` proto: `service Dashboard`
- Three methods:
  - `GetStatistics(Empty) returns (DashboardStatsResp)` — aggregate counts from multiple tables
  - `GetRecentOperations(RecentOpsReq) returns (RecentOpsResp)` — last N operation logs
  - `GetSystemInfo(Empty) returns (SystemInfoResp)` — runtime info (Go version, OS, uptime, CPU cores, etc.)
- Follow existing layer pattern: `server → logic → service → repository`
- System info can be gathered from Go stdlib (`runtime` package); no external dependency needed

## Risks / Trade-offs

- **[Risk] In-memory blocklist lost on restart** → Acceptable for admin system. Users re-login after restart.
- **[Risk] Captcha codes lost on restart** → Same as original's `DefaultMemStore`. Users just get new captcha.
- **[Risk] Dashboard stats query performance** → All aggregate queries hit indexed primary keys. Acceptable for admin dashboard refresh rates.
- **[Risk] File download endpoint exposes file paths** → Files are served by reading from the configured upload directory using database-stored paths. Validate paths to prevent directory traversal.
- **[Trade-off] Batch delete uses POST not DELETE** — Following existing rewrite conventions (POST for all mutations).

## Migration Plan

### Phase 1: Core Auth
1. Add `GenerateToken()`/`ParseToken()` to `JWTManager` in rpc/basis
2. Update `Login()` to return real JWT with user claims
3. Add captcha handler at gateway (base64Captcha)
4. Add token blocklist at gateway
5. Add logout endpoint at gateway
6. Wire blocklist check into JWT middleware

### Phase 2: Missing Endpoints
7. Add `common_pb.IdsReq` proto for batch operations
8. Add API batch delete
9. Add cron batch delete
10. Add operation log single+batch delete
11. Add dict detail flat list
12. Add file download endpoint
13. Add button batch check permission

### Phase 3: Dashboard Module
14. Create dashboard proto (sysMonitor)
15. Implement dashboard server/logic/service/repository
16. Add dashboard HTTP handler at gateway
17. Wire routes

### Phase 4: Verification
18. Build both binaries, run vet, verify no regressions
19. Manual test: captcha → login → protected endpoint → logout → protected endpoint rejected
20. Manual test: batch operations, file download, dashboard

## Open Questions

- Dashboard system info: use `github.com/shirou/gopsutil` (adds dependency) or Go stdlib `runtime` only? (Decision: start with stdlib `runtime` + `net/http` for simplicity)
- File download path traversal: what's the safest way to validate? (Decision: clean path, reject paths with `..`, check against configured upload directory)
