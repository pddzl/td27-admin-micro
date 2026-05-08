## Why

The current service only exposes gRPC endpoints (port 8080, etcd key basis.rpc). Web and mobile clients need REST/HTTP. A separate API gateway process in api/ bridges this gap.

## What Changes

- Create a new go-zero rest service at api/basis/ as a separate process
- Implement REST/HTTP endpoints for all 13 gRPC services
- Add JWT auth middleware at the HTTP layer
- Add CORS support for browser clients
- Exclude observability (no Prometheus/Jaeger/OpenTelemetry)
- Gateway discovers the gRPC basis service via etcd service discovery

## Capabilities

### New Capabilities
- http-api-gateway: REST/HTTP translation layer for all gRPC services

### Modified Capabilities
(N/A)

## Impact

- New process in api/basis/ — separate binary, separate deployment
- Shares proto-generated types with rpc/basis
- Both processes register to etcd — gateway discovers basis.rpc automatically
- No changes to existing rpc/basis service
