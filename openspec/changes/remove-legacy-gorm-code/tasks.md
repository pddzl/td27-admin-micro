## 1. Delete dead legacy files

- [x] 1.1 Delete `rpc/basis/internal/model/sysManagement/user_model.go`
- [x] 1.2 Delete `rpc/basis/internal/model/sysManagement/role_model.go`
- [x] 1.3 Delete `rpc/basis/internal/model/sysManagement/menu_model.go`
- [x] 1.4 Delete `rpc/basis/internal/model/sysManagement/user_repository.go`
- [x] 1.5 Delete `rpc/basis/internal/model/sysManagement/role_repository.go`
- [x] 1.6 Delete `rpc/basis/internal/model/sysManagement/menu_repository.go`
- [x] 1.7 Delete `rpc/basis/internal/model/sysManagement/data_permission.go`
- [x] 1.8 Delete `rpc/basis/internal/model/common.go`
- [x] 1.9 Delete `rpc/basis/internal/model/sysManagement/menu_dto.go` (referenced deleted types)
- [x] 1.10 Delete `rpc/basis/internal/model/sysManagement/user_dto.go` (referenced deleted types)
- [x] 1.11 Delete `rpc/basis/internal/model/sysManagement/role_dto.go` (empty file)

## 2. Remove GORM dependency

- [x] 2.1 Remove `gorm.io/gorm` from `go.mod`
- [x] 2.2 Run `go mod tidy` to clean up indirect deps

## 3. Verify

- [x] 3.1 Run `go build ./rpc/basis/...`
- [x] 3.2 Run `go build ./api/gateway/...`
- [x] 3.3 Run `go vet ./...`
