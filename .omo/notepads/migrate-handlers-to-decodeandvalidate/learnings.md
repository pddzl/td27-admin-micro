# Learnings

- Migrated `GetDeptDescendants` and `DeleteDept` in `dept.go` to use `pkg.DecodeAndValidate` + `pkg.WriteJson`
- `encoding/json` import retained — still referenced by other handlers in the same file
- Use `pkg.DecodeAndValidate(body, &req)` — returns error with user-facing message on decode/validation failure
- Use `pkg.WriteJson(w, status, pkg.Error(code, msg))` or `pkg.WriteJson(w, status, pkg.Success(data))` for all response writes
