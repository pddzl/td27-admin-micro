# Service Startup Test - Learnings

## Test Date: 2026-05-07

### Summary
The basis service starts successfully and initializes all dependencies correctly.

### Key Findings

1. **Config field name mismatch (CRITICAL)**:
   - The YAML config (`basis.yaml`) used `Dbname` as the key under `Pgsql:`
   - The Go struct tag expects `db-name` (`mapstructure:"db-name"`)
   - This caused `conf.MustLoad` to fail with: `field "pgsql.db-name" is not set`
   - **Fix applied**: Changed `Dbname: td27` → `db-name: td27`
   - All other config sections already used correct hyphenated keys

2. **Configuration**:
   - Mode: `dev` (enables gRPC reflection)
   - ListenOn: `0.0.0.0:8080`
   - Both PostgreSQL and etcd run via Docker on the same host

3. **Docker dependencies exist**:
   - PostgreSQL 14.22 on port 5432
   - etcd v3.6.7 on ports 2379-2380
   - Both are running in Docker containers

### Successful Initialization Sequence
All 4 steps completed in <1 second:
1. PostgreSQL connection - success
2. Casbin enforcer - success
3. JWT manager - success
4. Cron scheduler - success

### gRPC Services Registered (14 total)
- **Authority**: User, Role, Permission, Menu, Dept, Dict, API, Button
- **Tool**: File, Cron, Cache, ServiceToken
- **Monitor**: OperationLog
- **Infra**: Health (SERVING), Reflection (v1 + v1alpha)

### Verified Service Health
- gRPC port 8080 confirmed listening via `lsof`
- `grpcurl` confirmed all services registered
- Health check returned `SERVING`
