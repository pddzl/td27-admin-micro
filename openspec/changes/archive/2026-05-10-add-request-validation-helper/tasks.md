## 1. Add validation dependency and helper

- [x] 1.1 Add `github.com/go-playground/validator/v10` to go.mod
- [x] 1.2 Create `pkg/validate.go` with `DecodeAndValidate(body, req)` helper and package-level `validate` singleton

## 2. Migrate login handler

- [x] 2.1 Add `validate:"required"` tags to the login inline struct (Username, Password, CaptchaId, Captcha)
- [x] 2.2 Replace manual `json.Decode` + manual empty-field checks with `pkg.DecodeAndValidate(r.Body, &req)`
- [x] 2.3 Ensure `pkg.WriteJson` is used consistently

## 3. Migrate single-field ID handlers (5 handlers)

- [x] 3.1 Migrate `button.go DeleteButton` — `{ Id uint64 }`
- [x] 3.2 Migrate `api.go DeleteApi` — `{ Id uint64 }`
- [x] 3.3 Migrate `dept.go DeleteDept` — `{ Id uint64 }`
- [x] 3.4 Migrate `dict.go DeleteDict` and `DeleteDictDetail` — `{ Id uint64 }`
- [x] 3.5 Migrate `operation_log.go DeleteOperationLog` — `{ Id uint64 }`

## 4. Migrate multi-ID handlers (3 handlers)

- [x] 4.1 Migrate `api.go DeleteApiByIds` — `{ Ids []uint64 }`
- [x] 4.2 Migrate `cron.go DeleteCron` — `{ Ids []uint64 }`
- [x] 4.3 Migrate `operation_log.go DeleteOperationLogByIds` — `{ Ids []uint64 }`

## 5. Migrate special-field handlers (4 handlers)

- [x] 5.1 Migrate `api.go ListApiByGroup` — `{ GroupEn string }`
- [x] 5.2 Migrate `button.go GetButtonsByPagePath` — `{ PagePath string }`
- [x] 5.3 Migrate `button.go GetUserButtons` — `{ RoleIds []uint64 }`
- [x] 5.4 Migrate `button.go BatchCheckPermission` — `{ ButtonCodes, RoleIds }`
- [x] 5.5 Migrate `dept.go GetDeptDescendants` — `{ DeptId uint64 }`

## 6. Standardize error response format

- [x] 6.1 Audit all 15 handlers for raw `w.Header().Set + w.WriteHeader + json.NewEncoder(w).Encode(...)` patterns
- [x] 6.2 Replace any remaining raw error writes with `pkg.WriteJson(w, status, pkg.Error(...))`

## 7. Verify

- [x] 7.1 Run `go build ./...` — verify compilation
- [x] 7.2 Run `go vet ./...` — verify no new issues
- [x] 7.3 Run `go test ./...` — verify tests pass
