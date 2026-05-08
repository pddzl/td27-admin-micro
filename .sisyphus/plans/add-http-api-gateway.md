# Add HTTP API Gateway

## What
Create a separate go-zero HTTP REST API gateway process in `api/basis/` that exposes all 13 gRPC services to web/mobile clients via REST/HTTP. The gateway discovers the gRPC service via etcd (`basis.rpc`).

## Architecture

```
┌────────────────────────────────────────────────────┐
│  api/basis (port 8888)     rpc/basis (port 8080)  │
│  ┌─────────────────┐      ┌──────────────────┐    │
│  │ Public Routes    │      │  gRPC Server     │    │
│  │ POST /api/login  │      │  basis.rpc (etcd)│    │
│  │ GET  /api/health │      │  13 services     │    │
│  └────────┬─────────┘      └────────┬─────────┘    │
│           │                        │                │
│  ┌────────▼─────────┐              │                │
│  │ Private Routes   │     gRPC call                 │
│  │ (JWT required)   │◄─────────────────────────────│
│  │ /api/users/*     │                              │
│  │ /api/roles/*     │                              │
│  │ /api/menus/*     │                              │
│  │ /api/permissions │                              │
│  │ ...              │                              │
│  └──────────────────┘                              │
│                                                    │
│  Middleware chain: JWT → CORS → OpLog              │
└────────────────────────────────────────────────────┘
```

## Directory Structure

```
api/basis/
├── api.go                    # Entry point
├── etc/
│   └── gateway.yaml          # Config
├── internal/
│   ├── config/
│   │   └── config.go         # Gateway config
│   ├── svc/
│   │   └── servicecontext.go # gRPC client connections
│   ├── middleware/
│   │   └── jwt.go            # JWT middleware
│   ├── handler/
│   │   ├── login.go          # Auth handlers (public)
│   │   ├── user.go           # User CRUD handlers
│   │   ├── role.go
│   │   ├── menu.go
│   │   ├── permission.go
│   │   ├── dept.go
│   │   ├── dict.go
│   │   ├── api.go
│   │   ├── button.go
│   │   ├── file.go
│   │   ├── cron.go
│   │   ├── cache.go
│   │   ├── service_token.go
│   │   └── operation_log.go
│   └── router/
│       └── routes.go         # All route definitions
└── types/                    # HTTP-specific request/response types
    └── response.go
```

## Tasks

### 1. Scaffold the gateway process
- Create `api/basis/api.go` entry point using `rest.MustNewServer` with graceful shutdown
- Create `api/basis/etc/gateway.yaml` config
- Create `api/basis/internal/config/config.go` extending `rest.RestConf` with JWT settings
- Create `api/basis/internal/svc/servicecontext.go` with gRPC client connections (discover via etcd)

### 2. JWT middleware + response types
- Create `api/basis/internal/middleware/jwt.go` — validates JWT from `Authorization` header, extracts user info into context
- Create `api/basis/types/response.go` — standard JSON response format `{code, data, msg}`

### 3. Auth handlers (public)
- `POST /api/login` — calls gRPC User.Login, returns JWT token
- `GET /api/health` — health check

### 4. User handlers (private)
- `GET /api/user/:id` → User.GetUserInfo
- `POST /api/user/list` → User.ListUser
- `POST /api/user/create` → User.CreateUser
- `PUT /api/user/update` → User.UpdateUser
- `POST /api/user/delete` → User.DeleteUser
- `POST /api/user/password` → User.ModifyPassword
- `POST /api/user/active` → User.SwitchUserActive
- `POST /api/user/roles` → User.AssignRoles

### 5. Role handlers (private)
- `GET /api/role/:id`, `POST /api/role/list`, `POST /api/role/create`, etc.

### 6. Department, Menu, Dict, API, Button handlers (private)

### 7. Permission handlers (private)

### 8. Tool handlers (private) — File, Cron, Cache, ServiceToken

### 9. Monitor handler (private) — OperationLog

### 10. Router registration + build
- Wire all handlers in `routes.go`
- Build and verify the gateway compiles
