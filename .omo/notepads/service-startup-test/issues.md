# Service Startup Test - Issues

## Issues Found

### 1. Config field name mismatch (FIXED)
- **File**: `rpc/basis/etc/basis.yaml`
- **Problem**: Pgsql section used `Dbname: td27` but Go struct expects `db-name`
- **Error**: `field "pgsql.db-name" is not set`
- **Impact**: Service would not start at all without this fix
- **Resolution**: Changed to `db-name: td27`

### 2. No graceful etcd handling without etcd
- **Observation**: The service would likely hang/block if etcd is unavailable, since `zrpc.MustNewServer` with etcd config tries to connect on `Start()`
- **Not tested**: Etcd was running, so this path wasn't verified
- **Potential issue**: In a local dev environment without etcd, the service would block on startup
