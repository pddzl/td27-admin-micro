## 1. Proto Split

- [ ] 1.1 Create `rpc/basis/proto/sysManagement/dict_detail.proto` with DictDetail service RPCs
- [ ] 1.2 Update `rpc/basis/proto/sysManagement/dict.proto` to remove detail-specific messages
- [ ] 1.3 Regenerate both proto files

## 2. Logic Split

- [ ] 2.1 Create `rpc/basis/internal/logic/sysManagement/dict_detail.go` with DictDetailLogic
- [ ] 2.2 Remove dict_detail methods from `logic/sysManagement/dict.go`

## 3. Server Split

- [ ] 3.1 Create `rpc/basis/internal/server/sysManagement/dict_detail.go` with DictDetailServer
- [ ] 3.2 Remove dict_detail methods from `server/sysManagement/dict.go`
- [ ] 3.3 Register DictDetailServer in `basis.go`

## 4. Gateway Split

- [ ] 4.1 Create `api/gateway/internal/handler/basis/sysManagement/dict_detail.go` with DictDetailHandler
- [ ] 4.2 Add DictDetailClient to gateway `svc/servicecontext.go`
- [ ] 4.3 Remove dict_detail handlers from `handler/basis/sysManagement/dict.go`
- [ ] 4.4 Register dict_detail routes in `handler.go`

## 5. Verify

- [ ] 5.1 Run `go build ./rpc/basis/...`
- [ ] 5.2 Run `go build ./api/gateway/...`
- [ ] 5.3 Run `go vet ./...`
