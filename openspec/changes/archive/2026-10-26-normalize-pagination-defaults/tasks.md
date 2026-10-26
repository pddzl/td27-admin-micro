# Tasks

## 1. Shared helper

- [x] 1.1 Add `NewPageInfo(page, pageSize int) *PageInfo` to
      `rpc/basis/internal/model/common/model.go` with an unexported
      `normalize()` helper (Page<=0 -> 1, PageSize<=0 -> 10).

## 2. Call-site replacement

- [x] 2.1 sysManagement: replace `PageInfo` literals (and inline default
      blocks) with `common.NewPageInfo(...)` in api.go, button.go, dict.go,
      permission.go, role.go, user.go; delete the dead clamp block in
      menu.go (ListMenu never builds a PageInfo).
- [x] 2.2 sysManagement: convert dept.go and dict_detail.go, which built
      PageInfo with no defaults (latent LIMIT 0), to the constructor.
- [x] 2.3 sysTool: replace literals in cache.go, cron.go, file.go,
      service_token.go with the constructor.
- [x] 2.4 sysMonitor: replace the literal in operation_log.go (and the
      inline defaults added earlier today) with the constructor.

## 3. Verification

- [x] 3.1 `GONOSUMCHECK='*' GONOSUMDB='*' GOFLAGS='-mod=mod' go vet ./rpc/basis/...`
- [x] 3.2 `GONOSUMCHECK='*' GONOSUMDB='*' GOFLAGS='-mod=mod' go test ./rpc/basis/...`
- [x] 3.3 Manual (2026-10-26, live gateway): `POST /api/opl/list` with `{}`
      returns 10 rows and `"total": 851`; explicit `page_size` honored
      (3 rows for page_size=3). List invariant holds.
- [x] 3.4 `grep -rn "PageInfo{" rpc/basis/` returns only model.go and no
      remaining inline default blocks in logic/.
