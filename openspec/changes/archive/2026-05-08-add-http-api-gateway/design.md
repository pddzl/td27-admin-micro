## Context

The api/ directory currently exists but is empty. A new go-zero rest HTTP server will be created there as a separate process.

## Goals / Non-Goals

**Goals:**
- REST API for all 13 gRPC services
- JWT auth at HTTP layer
- CORS support
- Separate deployable binary

**Non-Goals:**
- No observability infrastructure (Prometheus/Jaeger/OTel)
- No changes to the existing rpc/basis service

## Decisions

1. Use go-zero rest framework (same as rpc/basis uses zrpc)
2. Manual HTTP handlers (not .api file generation) — cleaner mapping to existing gRPC client stubs
3. JWT middleware validates tokens at gateway, passes user context to handlers
4. Public routes: /api/login, /api/health — no auth required
5. Private routes: everything else — JWT required
6. gRPC client connects to basis.rpc via etcd

## Risks / Trade-offs

- Extra deployment: two processes instead of one
- Network hop: HTTP → gRPC adds latency
