## Context

Current state: We have a basic go-zero gRPC authority service in `rpc/basis` with minimal functionality. We need to rewrite this service to incorporate all the mature, production-tested business logic from the `td27-admin/server` reference repository, while retaining go-zero framework conventions, gRPC interface, and etcd service discovery.

## Goals / Non-Goals

**Goals:**
- Migrate 100% of non-observability business logic from reference repository
- Implement strict layered architecture: gRPC handler (controller) → service → repository → model
- Maintain full compatibility with go-zero v1.9.4 framework features
- Retain gRPC interface with etcd service discovery for inter-service communication
- Exclude all observability components (OpenTelemetry, sysMonitor module) as requested
- Reuse existing PostgreSQL database schema from reference repository without modifications

**Non-Goals:**
- No HTTP API layer implementation (pure gRPC service only)
- No observability infrastructure (no Prometheus, no Jaeger, no OpenTelemetry instrumentation); keep operation logs and dashboards functionality
- No multi-service split (remain as single monolithic service in `rpc/basis`)
- No changes to external systems (etcd, MySQL database remain unchanged)
- No changes to authentication/authorization flow beyond what exists in reference repo

## Decisions

1. **Layered Architecture Pattern**
   - Decision: Adapt reference repo's controller/service/repository pattern to fit go-zero conventions
   - Rationale: Balances reference repo's proven structure with go-zero's recommended practices
   - Layers mapping:
     - `controller`: Go-zero gRPC handlers in `internal/server/` directory
     - `service`: Business logic layer in new `internal/service/` directory
     - `repository`: DB operation layer in new `internal/repository/` directory
     - `model`: DB entity definitions in `internal/model/` directory (reused from reference repo)

2. **Protobuf Definition Structure**
   - Decision: Split proto files by functional module matching reference repo structure
   - Rationale: Makes it easy to map reference repo functionality to gRPC methods, improves maintainability
   - Structure:
     - `proto/authority/`: All sysManagement module protos (user, role, permission, etc.)
     - `proto/tool/`: All sysTool module protos (cache, cron, file, service token)

3. **Shared Dependencies Management**
   - Decision: Initialize all shared dependencies once in `internal/svc/servicecontext.go`
   - Rationale: Follows go-zero best practices, avoids duplicate initialization, simplifies dependency injection
   - Dependencies included: GORM DB pool, Casbin enforcer, JWT manager, Cron scheduler

4. **Logging & Observability**
   - Decision: Use go-zero built-in logging exclusively, remove all Prometheus, Jaeger, and OpenTelemetry instrumentation code; retain full operation log and dashboard functionality from reference repo
   - Rationale: Aligns with requirement to exclude external observability infrastructure while keeping critical audit and monitoring features

5. **Service Discovery**
   - Decision: Use go-zero's native etcd service discovery with service key `basis.rpc`
   - Rationale: Maintains existing inter-service communication pattern, no changes needed for service consumers

6. **Database Integration**
   - Decision: Reuse reference repo's existing PostgreSQL schema and GORM model definitions with minimal modifications
   - Rationale: Avoids schema migration effort, maintains compatibility with existing data

7. **Configuration**
   - Decision: Adapt reference repo configuration parameters into go-zero's YAML config format
   - Rationale: Maintains go-zero's config loading pattern, no changes needed for config management processes

## Risks / Trade-offs

- **Risk**: Breaking change for existing service consumers → Mitigation: Generate and distribute new gRPC client stubs, communicate API changes to all consumers, provide transition period if required
- **Risk**: Go-zero framework incompatibility with reference repo code → Mitigation: Test each module incrementally during implementation, replace reference repo dependencies with go-zero compatible equivalents where needed
- **Risk**: Increased complexity from additional service/repository layers → Mitigation: Follow strict separation of concerns, keep layers thin and focused, no cross-layer calls allowed
