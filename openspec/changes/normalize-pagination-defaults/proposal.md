# Normalize Pagination Defaults in One Shared Method

## Why

Every list RPC in `rpc/basis` applies the same page/pageSize defaulting inline
(`page <= 0 -> 1`, `pageSize <= 0 -> 10`). The block is copy-pasted 12 times
across `internal/logic/sysManagement`, `internal/logic/sysTool`, and
`internal/logic/sysMonitor`:

- 11 logics receive flat `in.Page` / `in.PageSize` fields and clamp them on the
  protobuf request before building `common.PageInfo`.
- `OperationLogLogic.ListOperationLog` receives a nested `PageReq` message and
  clamps on the locally built `common.PageInfo` (added 2026-10-26 when the
  missing defaults caused `LIMIT 0` and an empty `list` at total 851).
- `MenuLogic.ListMenu` clamps values it never uses (it fetches all menus via
  `GetAll`), making its block pure dead code.

Drift risk: if the defaults ever change (e.g. pageSize 10 -> 20) or gain a cap
(max pageSize), 12 sites must be edited in sync; a missed site yields
inconsistent behavior. `common.PageInfo` is the single type all repositories
already consume, so it is the natural home for the invariant.

## What Changes

- Add a `NewPageInfo(page, pageSize int) *PageInfo` constructor to
  `rpc/basis/internal/model/common/model.go` that clamps `Page <= 0` to 1 and
  `PageSize <= 0` to 10.
- Convert every `common.PageInfo` literal construction (13 sites) to
  `common.NewPageInfo(...)`:
  - sysManagement: api, button, dept, dict, dict_detail, permission, role,
    user
  - sysTool: cache, cron, file, service_token
  - sysMonitor: operation_log (folds in the inline defaults added earlier
    today; the user-visible behavior of that fix is unchanged)
- Delete the dead clamp block in `MenuLogic.ListMenu` (it never builds a
  `PageInfo`; `GetAll` ignores pagination).
- Behavior change beyond pure refactor: `ListDept` and `ListDictDetail`
  previously built `PageInfo` with no defaults at all (same latent `LIMIT 0`
  bug family as the operation-log issue); they now get the shared defaults.

## Impact

- Affected code: 1 new constructor in `internal/model/common`, mechanical
  edits in 13 logic files plus dead-code removal in menu.go. No proto,
  gateway, frontend, DB, or config changes.
- Gateway handlers keep their existing inline defaults for missing/zero
  `page_size` (api.go, dict.go, button.go) — out of scope here; the RPC-level
  default already guarantees correct behavior.
- Verification: `go vet ./rpc/basis/...`, `go test ./rpc/basis/...`, plus a
  manual check that `POST /opl/list` with `{}` still returns 10 rows.
