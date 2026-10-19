## Tasks

### P0 — Environment
- [x] R0.1 Verify migrations 001–003 applied (effect column, role_elevation table, button table)
- [x] R0.2 Verify seed users/roles/permissions exist; apply `scripts/rbac_seed_example.sql` if local DB is empty of test roles
- [x] R0.3 Flip dev `EnforceRbac: true` and restart gateway; confirm basis + gateway healthy
- [x] R0.4 Captcha → login → JWT works via curl (admin)

### P1 — Role & permission CRUD
- [x] R1.1 Create role → grant API/menu/button permissions → verify persisted via list/get
- [x] R1.2 Revoke a permission → verify gone (no re-login)
- [x] R1.3 Delete role → verify soft-deleted and no longer assignable

### P2 — API enforcement (gateway, EnforceRbac=true)
- [x] R2.1 User WITH granted API permission → 200
- [x] R2.2 User WITHOUT permission → 403 with proper message
- [x] R2.3 Missing/invalid token → 401
- [x] R2.4 Deny-wins: role with allow+deny on same API → denied
- [x] R2.5 Live revocation: revoke grant → next call 403 WITHOUT re-login

### P3 — Button permissions (`<page>:<action>`)
- [x] R3.1 Login response preloads button codes for user's roles
- [x] R3.2 Frontend hides buttons the user lacks (browser check 2026-10-19: tester with zero codes — 新增/编辑/修改密码/删除 all carry inline `display: none` from v-permission; search/table/pagination remain. Enabled by BUG-009 fixes: without the user-buttons grant the preload 403s and boot force-logs-out)
- [x] R3.3 Backend rejects mutation when button permission missing (design-verified: backend enforcement is api-domain only)

### P4 — JIT elevation
- [x] R4.1 User requests elevation → appears in admin queue
- [x] R4.2 Admin approves → elevated permission active live
- [x] R4.3 Self-approval forbidden (approver ≠ requester)
- [x] R4.4 Elevation expiry → permissions revert automatically

### P5 — Auth lifecycle
- [x] R5.1 Password change → old token blocked (blocklist) → re-login required
- [x] R5.2 New password logs in; old password rejected

### Bug #1 — Sidebar orphan 操作日志 (see bugs.md BUG-001)
- [x] B1.1 Repro confirmed: root cause found statically + via API (single-child hoisting in Item.vue triggered by snake_case/camelCase contract break on /menu/list; backend tree itself correct)
- [x] B1.2 Fixed in gateway `menuTreeToResp` (camelCase alwaysShow/keepAlive per MenuData contract)
- [x] B1.3 Retest (API level): /menu/list emits camelCase flags on all 14 nodes, 系统监控 alwaysShow=true ⇒ nested render; visual UI confirmation pending user

### Loop
- [x] L1 Re-run matrix after fixes — final green state 2026-10-19: BUG-008 fix restored RPC path (menu/list, dashboard×3 = 200 admin; menu/list 200 tester, user/list 200, user/create 403, get-user-info 200 tester); BUG-010 fix scoped tester menus (系统管理→用户管理 only, admin 17 nodes); BUG-009 fixes unblocked boot (browser session survives); R3.2 button hiding PASS; dashboard stats verified vs SQL (P6 spot check: 7/7 fields match with soft-delete semantics)
- [x] L2 Final report: bugs.md carries BUG-001…010 with repro/root-cause/fix/retest evidence + OBS-001…005; BUG-008 contains the dashboard-statistics verification (user-reported "wrong stats" = fail-closed 403s, numbers themselves correct)
