## Why

The project was migrated from the original Gin-based codebase (`Mine/td27-admin/server`) to this Go-zero rewrite (`td27-admin-micro`). However, a systematic comparison reveals that several functions and modules were not carried over:

1. **Captcha** — No captcha verification on login, leaving brute force unprotected
2. **Logout** — No endpoint to terminate sessions or invalidate tokens
3. **JWT token generation** — Login returns `"placeholder_token"` instead of a real JWT
4. **Multiple CRUD endpoints** — Batch operations, file download, and dashboard module are entirely missing

## What Changes

### Auth & Session
- **Captcha**: Add captcha generation and verification for the login flow
- **Logout**: Add a logout endpoint that invalidates the JWT token via server-side blocklist
- **Real JWT Token Generation**: Replace `"placeholder_token"` with proper JWT generation using existing `JWTManager` config
- **JWTManager Enhancement**: Add `GenerateToken()` and `ParseToken()` methods to the `JWTManager` struct in `rpc/basis`
- **Token Blocklist Service**: Add token invalidation mechanism (in-memory, designed for Redis swap-in)

### Batch Operations
- **API batch delete**: Add `POST /api/apis/delete-by-ids` for batch deleting API entries
- **Cron batch delete**: Add `POST /api/cron/delete-by-ids` for batch deleting cron entries
- **OperationLog delete**: Add `POST /api/operation-log/delete` (single) and `POST /api/operation-log/delete-by-ids` (batch)

### Data Retrieval
- **DictDetail flat list**: Add `POST /api/dict/detail/flat` for flat (non-paginated) dict detail listing by dict ID
- **File download**: Add `GET /api/file/download/:id` for downloading file binary content
- **Button BatchCheckPermission**: Add `POST /api/button/batch-check` for checking multiple button permissions at once

### Dashboard Module (Entirely New)
- **Dashboard statistics**: Add `GET /api/dashboard/statistics` with aggregate counts (users, roles, APIs, depts, etc.)
- **Dashboard recent operations**: Add `GET /api/dashboard/recent-operations` for latest audit log entries
- **Dashboard system info**: Add `GET /api/dashboard/system-info` for server runtime info (Go version, OS, uptime, etc.)

## Capabilities

### New Capabilities
- `captcha`: CAPTCHA image generation and verification for login protection
- `logout`: Token invalidation and session termination via logout endpoint
- `jwt-generation`: Proper JWT token generation with user claims replacing the placeholder token
- `token-blocklist`: Server-side token invalidation mechanism supporting in-memory and Redis backends
- `api-batch-delete`: Batch deletion of API entries
- `cron-batch-delete`: Batch deletion of cron entries
- `operation-log-delete`: Single and batch deletion of operation logs
- `dict-detail-flat`: Flat (non-paginated) dict detail listing by dict ID
- `file-download`: Binary file download endpoint
- `button-batch-check`: Batch permission checking for multiple button codes
- `dashboard`: Dashboard statistics, recent operations, and system information endpoints

### Modified Capabilities
- `user-auth`: Login flow now includes captcha verification and returns a real JWT token; logout endpoint added
- `api`: API management now includes batch delete
- `cron`: Cron management now includes batch delete
- `operation-log`: Operation log now includes single and batch delete endpoints
- `dict`: Dict detail now includes flat list endpoint
- `file`: File management now includes download endpoint
- `button`: Button management now includes batch check permission endpoint

## Impact

- **rpc/basis**:
  - Add `GenerateToken()` and `ParseToken()` methods to `JWTManager`
  - Update `UserLogic.Login()` to generate real JWT tokens
  - Add proto messages for dashboard service (DashboardStats, RecentOperation, SystemInfo)
  - Add proto messages for batch delete requests
  - Add file download streaming support
  - Add `DashboardServer`, `DashboardLogic`, `DashboardService`, `DashboardRepository`
- **api/gateway**:
  - Add public endpoint: `GET /api/captcha`
  - Add authenticated endpoints: `POST /api/logout`, batch deletes, file download, batch check, dashboard
  - Add token blocklist checking in JWT middleware
  - Add dashboard handler
- **Dependencies**:
  - Add `github.com/mojocn/base64Captcha` for captcha
  - `github.com/shirou/gopsutil` for dashboard system info (optional, can use Go stdlib)
- **No breaking changes** to existing endpoints
