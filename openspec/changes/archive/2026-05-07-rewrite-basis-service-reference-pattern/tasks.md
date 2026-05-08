## 1. Preparation & Setup

- [x] 1.1 Backup existing `rpc/basis` directory content to temporary location
- [x] 1.2 Add required dependencies to go.mod (Casbin, jwt-go, cron, etc.)
- [x] 1.3 Create new directory structure: `internal/service`, `internal/repository`, `proto/sysManagement`, `proto/tool`, `proto/monitor`

## 2. Core Infrastructure Implementation

- [x] 2.1 Update `internal/config/config.go` to include all reference repo configuration parameters
- [x] 2.2 Update `etc/basis.yaml` with new configuration values
- [x] 2.3 Implement DB initialization in `internal/initialization` (GORM with PostgreSQL, match reference repo settings)
- [x] 2.4 Initialize Casbin RBAC enforcer in service context
- [x] 2.5 Initialize JWT manager in service context
- [x] 2.6 Initialize cron scheduler in service context

## 3. Model Layer Implementation

- [x] 3.1 Migrate all sysManagement model definitions from reference repo to `internal/model/sysManagement`
- [x] 3.2 Migrate all sysTool model definitions from reference repo to `internal/model/tool`
- [x] 3.3 Migrate all sysMonitor model definitions from reference repo to `internal/model/monitor`
- [x] 3.4 Verify all GORM tags and model relationships match reference repo exactly

## 4. Repository Layer Implementation

- [x] 4.1 Implement repository interfaces for all sysManagement models (user, role, permission, menu, dept, dict, api, button, data permission)
- [x] 4.2 Implement repository interfaces for all tool models (cache, cron, file, service token)
- [x] 4.3 Implement repository interfaces for all monitor models (operation log)
- [x] 4.4 Verify all CRUD operations match reference repo implementation exactly

## 5. Service Layer Implementation

- [x] 5.1 Implement sysManagement module services (user, role, permission, menu, dept, dict, api, button, data permission)
- [x] 5.2 Implement tool module services (cache, cron, file, service token)
- [x] 5.3 Implement monitor module services (operation log, dashboard metrics)
- [x] 5.3 Implement Casbin RBAC integration in permission service
- [x] 5.4 Implement JWT authentication logic in user service
- [x] 5.5 Verify all business logic matches reference repo exactly

## 6. Protobuf & gRPC Layer Implementation

- [x] 6.1 Write protobuf definitions for all sysManagement module gRPC methods
- [x] 6.2 Write protobuf definitions for all tool module gRPC methods
- [x] 6.3 Write protobuf definitions for all monitor module gRPC methods
- [x] 6.3 Generate gRPC/pb Go code using goctl
- [x] 6.4 Implement gRPC server handlers in `internal/server` (controller layer) (ALL services completed)
- [x] 6.5 Update service entry point to register all gRPC servers (ALL 13 services registered)

## 7. Testing & Validation

- [x] 7.1 Run `go mod tidy` to resolve all dependencies
- [x] 7.2 Run `go build rpc/basis/basis.go` to verify no compilation errors
- [x] 7.3 Test service startup with config file
- [x] 7.4 Verify etcd service discovery works correctly
- [x] 7.5 Clean up backup files and old unused code

## 🎉 ALL 34 TASKS COMPLETE
