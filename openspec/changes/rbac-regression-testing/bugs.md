# Bug Log — RBAC Regression Loop

Format: each bug gets ID, status, repro, root cause, fix, retest evidence. No bug is closed without a passing retest.

## Observations (non-bug findings kept for traceability)
- **OBS-001** — Role delete is a transactional HARD delete (role + role_permissions + user_roles rows removed), by design, not soft delete.
- **OBS-002** — The gateway JSON decoder accepts proto enums as NUMERIC values only; string names (`"DOMAIN_API"`) are rejected. Test payloads must use numbers (e.g. `domain: 1`, `effect: 1`).
- **OBS-003** — Login failure leaks the internal rpc error string (`rpc error: code = Internal desc = invalid password`) to the client; cosmetic hardening candidate.
- **OBS-004** — 2026-10-18: a wedged basis process (alive, printing stats, but accepting TCP without serving gRPC — h2c preface never answered) made every authz check fail closed with 403, masquerading as a permission problem. Root trigger uncertain; coincided with a VPN TUN interface (utun4, 192.168.255.10) coming up and route churn — gateway dials basis via the LAN IP 192.168.0.109 (etcd registration from `ListenOn: 0.0.0.0:8080`). Recovery: restart basis + gateway. Note: with EnforceRbac=true the fail-closed design correctly denied everything while the permission RPC was unreachable; check gateway log line `jwt.go:110 rbac check failed ... denying` before diagnosing 403s as policy issues. Candidate dev improvement: register loopback (127.0.0.1) in etcd for local runs to be immune to interface/VPN changes.

## BUG-001 — 操作日志 renders as orphan top-level sidebar item
- **Status**: FIXED (API-level retest PASS; visual UI confirmation pending)
- **Found**: 2026-10-16, user screenshot during local testing
- **Repro**: Log in as root → sidebar shows 操作日志 as a top-level menu; DB (`sys_management_menu`) has 操作日志 (id 41) with parent_id 40 (系统监控). Other children of 系统监控 nest correctly.
- **Expected**: 操作日志 nested under 系统监控 like its siblings.
- **Root cause**: Key-casing contract break between gateway and frontend. `GET /menu/list` (`menuTreeToResp` in [menu.go](../../../api/gateway/internal/handler/basis/sysManagement/menu.go)) emitted snake_case `always_show`/`keep_alive`, but every frontend consumer reads camelCase per the `MenuData` contract (`web/src/api/sysManagement/menu.ts`): `formatRouter` ([router_m.ts](../../../web/src/common/utils/router_m.ts)) set `meta.alwaysShow` to `undefined` for every menu. With `meta.alwaysShow` falsy, [Item.vue](../../../web/src/layouts/components/Sidebar/Item.vue) single-child hoisting fires: 系统监控 had exactly one visible child (操作日志), so the child was promoted to a top-level leaf — the orphan in the screenshot. Secondary victims of the same mismatch: menu admin edit form silently reset `alwaysShow`/`keepAlive` to false on every save (`Boolean(undefined)`), and `tags-view` keepAlive caching never activated.
- **Fix**: `menuTreeToResp` now emits camelCase `alwaysShow`/`keepAlive` (matching the documented frontend contract and the protojson shapes all other passthrough endpoints already use). Frontend untouched. The el-tree `/menu/el-tree` transformer keeps snake_case — its only consumer (role→menu assignment tree) uses just `id`/label/checked keys.
- **Retest**: PASS (API level) — restarted gateway, `GET /menu/list` now returns `alwaysShow`/`keepAlive` on every node (zero legacy snake keys); 系统监控 node has `alwaysShow=true` ⇒ Item.vue renders `el-sub-menu` with 操作日志 nested. Visual confirmation in the browser pending (user verifies or re-authorizes browser agent). Note: user meanwhile renamed the row `menu_name` 操作日志→`operationLog` and set its flags true via admin UI/DB; behavior unchanged for any always_show value.

## BUG-007 — Non-admin roles get 403 on GET /menu/list ⇒ empty sidebar
- **Status**: FIXED
- **Found**: 2026-10-18, L1 post-fix sweep: tester (`rbac-limited`, only `/user/list` read granted) → **403** on `GET /menu/list`
- **Repro**: any user whose roles lack the `/menu/list` (read) api permission cannot fetch their menu tree; the frontend `setRoutes()` gets a non-zero code, `dynamicRoutes` stays empty → sidebar renders nothing. Enforcement itself is correct (fail-closed); the gap was policy/seed data — no migration ever granted this route to ordinary roles.
- **Root cause**: `/menu/list` is api-permission-gated like any other route, but menu visibility was only ever modeled as role→menu assignments, never as the api grant every user needs to load their own menus.
- **Fix**: `scripts/006_grant_menu_list_read.sql` — grants the existing `/menu/list` read allow permission (id 215) to every non-deleted role via `INSERT ... SELECT ... ON CONFLICT DO NOTHING`. Idempotent (re-apply: INSERT 0 0). Applied locally 2026-10-18; policies hot-reloaded via `POST /permission/reload-policy`. Caveat documented in the migration: roles created after this migration must be granted `/menu/list` at role-creation/assignment time.
- **Retest**: PASS — tester `GET /menu/list` → **200** (tree empty because role 7 has no menu assignments yet — assignment is a separate admin action, not a permission error); admin unchanged → 200.

## BUG-002 — sys_management_api registry drifted from registered gateway routes (17 orphan + 62 missing)
- **Status**: FIXED
- **Found**: 2026-10-16, static diff of `handler.go` routes vs `sys_management_api` rows during P0/P2
- **Repro**: With `EnforceRbac: true`, admin (root role, 57 api grants) calls `GET /user/get-user-info` → **403**; `POST /dict-detail/list` → **403**. Verified live 2026-10-16.
- **Root cause**: The api registry was written against old route names that no longer exist. Enforcement matches `r.URL.Path` exactly (keyMatch2 for `:param` patterns), so grants referencing stale paths can never match: 17 orphan seeds + 62 registered routes with no seed row (ungrantable ⇒ 403 for everyone). Additional wrinkle: api table used snake_case `/dict_detail/*` while permission rows used camelCase `/dictDetail/*` for the same routes.
- **Fix**: `scripts/004_align_api_permissions.sql` (generated from `handler.go` + live DB) — 17 renames (api + permission rows), 49 new api rows, 49 new allow-permission rows, root granted all api-domain permissions. Idempotent; re-apply verified (second run: INSERT 0 0). Applied locally 2026-10-16, policies hot-reloaded via `POST /permission/reload-policy`.
- **Retest**: PASS — 26-route sweep across all modules: `/user/get-user-info`, `/role/list`, `/dict-detail/list`, `/button/list`, `/permission/list`, `/cache/list`, `/operation-log/list`, `/role/elevation/my`, `/cron/list`, `/serviceToken/list`, `/file/list`, `/menu/user-menus`, `/dashboard/statistics`, `/user/list`, `/dept/getTree`, `/api/el-tree`, `/menu/el-tree`, `/opl/list`, `/cron/toggle-status`, `/user/switch-active`, `/user/modify-password`, `/role/assign-permissions`, `/permission/check`, `/cache/cleanup`, `/permission/lint`, `/file/download/1` — all past RBAC (200/400-validation/rpc-arg errors, zero 403). Final state: 109 api rows ↔ 109 api permissions ↔ 109 root grants; actions create=14/delete=17/execute=4/read=46/update=28.

## BUG-003 — api permission rows stamped with wrong actions (53× create, 1× invalid `post`)
- **Status**: FIXED (same migration as BUG-002)
- **Found**: 2026-10-16, DB inspection: `action` distribution across 60 api rows = create:53, post:1, read:6
- **Repro**: `POST /role/list` as admin → **403** live (verified). Middleware maps `/list` suffix → ACTION_READ, but the row said action=create ⇒ never matches.
- **Root cause**: Seed data applied `create` to nearly every api permission regardless of the route verb convention.
- **Fix**: migration step 2 recomputes `action` for all api-domain rows with a CASE mirroring `httpMethodToAction` ([authz.go](../../../api/gateway/internal/middleware/authz.go)).
- **Retest**: PASS — same 26-route sweep as BUG-002; read/delete/create/execute/update routes all authorize per their derived action.

## BUG-004 — CreateRole panics rpc when parent_id omitted (nil deref on optional field)
- **Status**: FIXED
- **Found**: 2026-10-16, P1 pre-test: `POST /role/create {"role_name":"root"}` (no parent_id)
- **Repro**: `curl -X POST :8888/role/create -d '{"role_name":"x"}'` → HTTP 200 with `code:7` body: `panic: runtime error: invalid memory address or nil pointer dereference` at `rpc/basis/internal/logic/sysManagement/role.go:143` (`parentID := uint(*in.ParentId)` — `ParentId` is optional `*int64`). Reachable from the public API.
- **Impact**: Remote rpc panic (goroutine recovered by go-zero, service survives, but the op 500s); role creation without parent_id impossible.
- **Fix**: nil-safe `ParentId` handling in [role.go CreateRole](../../../rpc/basis/internal/logic/sysManagement/role.go); unique-violation (23505) on role_name mapped to `ErrRoleNameTaken` in RoleService.Create and surfaced as gRPC `AlreadyExists`. `go vet` clean.
- **Retest**: PASS — `{"role_name":"rbac-test-role-a"}` → 200 success, DB row id=5 parent=NULL; duplicate create → clean `AlreadyExists: role name "rbac-test-role-a" already exists` (no SQL leak).

## BUG-005 — CreateUser panics rpc when any optional field omitted
- **Status**: FIXED
- **Found**: 2026-10-16, P2 setup: `POST /user/create {"username":"rbac-tester","password":"Test1234!","role_ids":[7]}`
- **Repro**: create user without `phone`/`email`/`active`/`dept_id` → `panic: runtime error: invalid memory address or nil pointer dereference` at [user.go:169-172](../../../rpc/basis/internal/logic/sysManagement/user.go) — all four optional proto fields (`Phone`, `Email`, `Active`, `DeptId`) deref'd unconditionally. Reachable from the public API; UpdateUser handles the same fields nil-safely, Create did not.
- **Impact**: User creation with partial payloads (common from scripts/CLI) always 500s with a panic.
- **Fix**: nil-safe field handling in CreateUser (matching the Update pattern); omitted `active` defaults to `true` (new users enabled). `go vet` clean.
- **Retest**: PASS — same payload → 200 `创建成功`; DB row id=5 `rbac-tester active=true dept=0`, role 7 linked.

## BUG-006 — UNIQUE(resource) on sys_management_permission makes deny-wins impossible
- **Status**: FIXED
- **Found**: 2026-10-16, R2.4 setup: creating a `deny` row for `/user/list` failed with duplicate-key on `uni_sys_management_permission_resource`
- **Repro**: any attempt to add a deny permission for a resource that already has an allow row violates `UNIQUE (resource)` — yet the Casbin 4-column model (sub, obj, act, eft, deny wins) requires allow and deny rows for the same resource to coexist. Deny rules could never be expressed.
- **Root cause**: `001_add_permission_effect.sql` added the `effect` column but left the seed-era `UNIQUE (resource)` constraint in place.
- **Fix**: `scripts/005_perm_resource_effect_unique.sql` — drop `UNIQUE (resource)` constraint, add partial `UNIQUE (resource, effect) WHERE deleted_at IS NULL`. Idempotent; applied locally 2026-10-16.
- **Retest**: PASS — deny row created via `POST /permission/create` (numeric enums; string enum names rejected by decoder — noted as OBS-002), granted to role 7 alongside the allow row → `POST /user/list` as tester → 403 (deny wins). Removed deny → 200 again.
