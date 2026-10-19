# Regression Test Matrix

Sections in execution order. Skip a section only when the project lacks the
feature — say so explicitly in the report. Evidence format per case:
`<case id> <method> <path> as <principal> → <HTTP status> <one-line note>`.

## P0 — Environment

- [ ] Migrations/schema at expected version (list applied migrations).
- [ ] All services healthy; note PIDs/ports for restart during the run.
- [ ] Enforcement flag/config actually enabled (not just present).
- [ ] Login smoke: correct credentials → token; wrong password → rejected;
      expired/garbage token → 401 on a protected route.
- [ ] Principals ready: superadmin + one limited-scope user. Capture user id,
      role ids, and a DB query to list their live grants.

## P1 — CRUD per resource

For each business resource (read its routes from the Phase 1 inventory):

- [ ] Create with minimal payload (all optionals omitted) → success, DB row
      correct, defaults applied. Watch for nil-pointer panics on optional
      proto fields — the classic case.
- [ ] Create duplicate/uniqueness-violating payload → clean business error, no
      SQL error text leaked.
- [ ] Get / list / pagination bounds (page 0, negative, huge).
- [ ] Update partial payload → only intended fields change.
- [ ] Delete → verify persistence semantics (soft vs hard) and that links/joins
      are cleaned or filtered by `deleted_at IS NULL` everywhere.
- [ ] Delete test rows via the product's own APIs; verify DB counts reach zero.

## P2 — API enforcement

With the limited-scope principal:

- [ ] Granted read route → 200.
- [ ] Granted mutation route → 200/2xx.
- [ ] Non-granted read → 403 with a proper message (not a stack/SQL leak).
- [ ] Non-granted mutation → 403.
- [ ] No token → 401; garbage token → 401.
- [ ] Deny-wins (if deny effects exist): role holds allow AND deny on the same
      resource → denied.
- [ ] Live revocation: revoke a grant server-side → next call 403 WITHOUT
      re-login; re-grant → 200 again on the same token.
- [ ] Sweep: hit every route in the inventory with the superadmin principal —
      zero 403s expected (validation/rpc-arg errors are fine; they prove the
      route was reached). Any 403 here is inventory/policy drift.

## P3 — Button / field-level permissions (if present)

- [ ] Login response (or dedicated endpoint) preloads the principal's button/
      action codes for their roles.
- [ ] Frontend hides controls the principal lacks (browser check; needs user
      cooperation or browser automation approval).
- [ ] Backend behavior for missing button permission documented and verified —
      note if backend enforcement is intentionally only api-domain.

## P4 — JIT elevation / temporary grants (if present)

- [ ] User requests elevation → appears in an approver's queue.
- [ ] Approve → permission active immediately, no re-login.
- [ ] Self-approval forbidden (approver ≠ requester), enforced server-side
      (requester identity from token, never the body).
- [ ] Expiry passes → permissions revert automatically on the next check.
- [ ] Revoke before expiry → reverts immediately.

## P5 — Auth lifecycle

- [ ] Password change (self-service) → old token blocked on next use (401),
      i.e. blocklist/session invalidation works without waiting for expiry.
- [ ] Re-login with new password → 200; old password → rejected.
- [ ] Admin-forced operations (if present) behave the same.

## Drift checks (static, cheap, run early)

- [ ] Route inventory ↔ permission registry: no orphans, no ungrantable routes.
- [ ] Every registry row's action equals the method→action mapping output.
- [ ] Resource path spelling consistent across api table, permission rows, and
      registered routes (watch camelCase vs kebab/snake renames).

## P6 — Data correctness (stats/aggregates vs SQL ground truth)

For every dashboard/statistics/report endpoint, do NOT trust plausible-looking
numbers — cross-check each field against the database:

- [ ] Call the endpoint and list every numeric field.
- [ ] For each field, write the equivalent SQL directly against the tables the
      service queries (read the repository/SQL first to learn the exact
      semantics: soft-delete filter? time window? status filter?).
- [ ] Compare field by field; a mismatch is a bug unless the code's intent says
      otherwise (e.g. a deliberate `deleted_at IS NULL` filter — then record
      the semantic in the evidence line).
- [ ] After inserting and deleting a test row, the stat must move accordingly
      (proves the number is live, not cached at startup).

## Bug log template

```markdown
## BUG-<id> — <one-line title>
- **Status**: OPEN | FIXED | WONTFIX
- **Found**: <date>, <how discovered>
- **Repro**: <exact request/payload + observed result>
- **Expected**: <what should happen>
- **Root cause**: <file:line, mechanism>
- **Fix**: <what changed, migration id if any, idempotency note>
- **Retest**: PASS/FAIL — <evidence + date>
```
