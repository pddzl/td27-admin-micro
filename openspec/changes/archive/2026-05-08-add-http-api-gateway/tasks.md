## 1. Scaffold Gateway

- [x] 1.1 Create api/basis/api.go with rest.MustNewServer + graceful shutdown
- [x] 1.2 Create api/basis/etc/gateway.yaml config
- [x] 1.3 Create api/basis/internal/config/config.go extending RestConf
- [x] 1.4 Create api/basis/internal/svc/servicecontext.go with gRPC clients
- [x] 1.5 Create api/basis/types/response.go with standard JSON response format

## 2. Middleware & Auth

- [x] 2.1 Create api/basis/internal/middleware/jwt.go
- [x] 2.2 Implement public handlers: login, health

## 3. Authority Handlers

- [x] 3.1 Implement User handlers (CRUD + password + active + roles)
- [x] 3.2 Implement Role handlers (CRUD + permissions)
- [x] 3.3 Implement Menu handlers (CRUD + tree + user menus)
- [x] 3.4 Implement Dept handlers (CRUD + tree + descendants)
- [x] 3.5 Implement Dict handlers (CRUD + details)
- [x] 3.6 Implement API handlers (CRUD)
- [x] 3.7 Implement Button handlers (CRUD + user buttons)
- [x] 3.8 Implement Permission handlers (CRUD + check + reload)

## 4. Tool & Monitor Handlers

- [x] 4.1 Implement File handlers (upload + list + delete)
- [x] 4.2 Implement Cron handlers (CRUD + toggle + execute)
- [x] 4.3 Implement Cache handlers (get + set + delete + list)
- [x] 4.4 Implement ServiceToken handlers (CRUD + permissions + validate)
- [x] 4.5 Implement OperationLog handlers (list + cleanup)

## 5. Routes & Build

- [x] 5.1 Create api/basis/internal/router/routes.go wiring all handlers
- [x] 5.2 Run go build ./api/basis/... — verify compilation
