## Why

This change aligns the current go-zero gRPC authority service with the mature, production-tested business logic from the existing `td27-admin/server` reference repository. It retains all microservice benefits (gRPC interface, etcd service discovery, go-zero framework conventions) while eliminating duplicated effort by reusing proven authority system logic, excluding only observability components we don't need for this service.

## What Changes

- **BREAKING**: Complete rewrite of the `rpc/basis` monolithic gRPC service
- Migrate all non-observability business logic from `td27-admin/server` reference repository
- Implement strict controller → service → repository → model layered architecture pattern
- Exclude observability infrastructure: remove Prometheus, Jaeger, and OpenTelemetry initialization, use go-zero built-in logging only; retain full sysMonitor module functionality (operation logs, dashboards)
- Retain gRPC communication via etcd service discovery per go-zero conventions
- All existing gRPC APIs will be replaced with new APIs matching reference repository functionality

## Capabilities

### New Capabilities
- `user-management`: Complete user CRUD, authentication, password management
- `role-management`: Role CRUD, role permission assignment
- `permission-management`: Permission CRUD, API/resource access control
- `menu-management`: Menu tree management, user role menu assignment
- `department-management`: Department hierarchy management
- `dictionary-management`: System dictionary and dictionary detail management
- `api-management`: API endpoint management, permission mapping
- `button-management`: UI button permission management
- `data-permission`: Row-level data access control based on user roles
- `cache-management`: System cache viewing and clearing
- `cron-job-management`: Scheduled task CRUD, execution logging, manual trigger
- `file-upload`: File upload and management functionality
- `service-token-management`: Service-to-service authentication token management
- `operation-log-management`: System operation log recording and query functionality
- `dashboard`: System overview dashboard metrics calculation and query

### Modified Capabilities
(No existing capabilities modified - this is a complete rewrite of all functionality)

## Impact

- **Affected code**: Entire `rpc/basis` directory (all internal code, proto definitions, configuration)
- **APIs**: All existing gRPC APIs are replaced with new APIs matching reference functionality
- **Dependencies**: Add required dependencies from reference repo (Casbin, JWT, cron, etc.) that are compatible with go-zero framework
- **External systems**: No changes to etcd service discovery, PostgreSQL database connection remains the same
