# E2E Flow Catalog

Standard journeys for admin/backoffice systems. Pick what applies, adapt
names/codes per project. Actor legends: ADMIN = superadmin account, USER =
limited-permission account, ANY = both.

## J1 — Login journey (actor: ANY, real login required)

1. Open the app unauthenticated → redirected to login.
2. Wrong password → error shown, no session stored.
3. Expired/garbage captcha (if captcha exists) → rejected with a message.
4. Correct credentials + valid captcha → dashboard renders; session storage
   contains the token.
5. Reload the page → still authenticated (token persisted, guard passes).

## J2 — Logout / session invalidation (actor: ANY)

1. Trigger logout → session storage cleared, redirected to login.
2. Back-button or direct URL to an authenticated page → bounced to login.
3. If password change exists: change it in one session → the OTHER session's
   next API call gets 401 (blocklist) → forced re-login.

## J3 — CRUD journey per resource (actor: ADMIN)

1. Navigate to the resource page via menus (not URL) — proves menu wiring.
2. Create with minimal fields → success toast, row appears in table, DB row
   matches submitted values and defaults.
3. Edit the row → table and DB reflect the change.
4. Exercise one validation error (empty required field, duplicate name) →
   inline error, no DB write.
5. Delete → confirm dialog → row gone from table AND from the database
   (mind soft-delete: filtered out vs actually deleted).
6. Refresh the page → the deleted row stays gone (proves no stale cache).
7. Run the list invariant (SKILL.md Phase 3) on the resource's list endpoint
   after the create and delete refreshes — payload keys exactly as the
   frontend sends them.

## J4 — Permission-scoped experience (actors: ADMIN vs USER)

1. USER logs in → sidebar shows ONLY assigned menus (+ ancestor path).
2. USER opens an assigned page → action buttons (create/edit/delete) are
   hidden; read-only elements remain. Verify hidden means not-visible
   (computed display), not merely absent from a snapshot.
3. USER deep-links to an unassigned page → redirected or 403 view, never a
   rendered page with broken data.
4. ADMIN sees the same page fully equipped — proves gating is per-account,
   not broken for everyone.
5. While USER is logged in, ADMIN revokes a grant → USER's next relevant
   action reflects it without re-login (live enforcement).

## J5 — Dashboard / reports correctness (actor: ADMIN)

1. Note every displayed statistic.
2. Compute each via SQL with the service's exact semantics (soft-delete
   filter, time window).
3. Create a test row → the corresponding stat increments on reload; delete →
   decrements. Proves the numbers are live, not startup-cached.

## J6 — Profile journey (actor: USER)

1. Open profile → fields match the DB row for the account.
2. Change password (wrong old password → rejected; correct → success).
3. Session invalidation per J2.3; re-login with the new password works.

## J7 — Navigation integrity (actor: USER)

1. Every sidebar item the account can see opens a rendered page (no blank,
   no 404, no console error) **with its data present**: wherever the database
   guarantees rows (the operation log always has entries after any earlier
   journey), the page's table/list must show them — render-only is not a
   pass.
2. Collapsed groups expand (hover or click); their children are reachable —
   remember teleported popups render outside the sidebar container.
3. Browser refresh on a deep link re-lands on the same page (router guards
   re-run cleanly from a cold load).

## J8 — JIT elevation / temporary grants (actor: USER + ADMIN, if the system has it)

1. USER requests elevation to a role (reason + duration) → appears in the
   approver queue.
2. USER tries to self-approve → rejected server-side (identity from token,
   never the body).
3. ADMIN approves → USER's next action succeeds WITHOUT re-login (live
   effect); the elevated menu/button/page is visible.
4. Wait out the expiry (or shrink it in a test env) → USER's next action is
   denied again; UI reverts to the pre-elevation view.
5. ADMIN revokes mid-window → effect disappears immediately.

## J9 — Binary and file flows (actor: ADMIN, if applicable)

1. Upload a file through the UI → success toast, row in the file list, file
   present on disk/object store, DB row correct (size, name).
2. Download via the UI → the browser receives the file (assert download
   event and non-empty content, not just HTTP 200 — blob endpoints can
   return JSON errors with 200).
3. Delete → row gone; on-disk artifact handled per design (removed or
   orphaned-by-design — record which).

## J10 — Admin assignment journeys (actor: ADMIN, target USER)

The management UIs that change what OTHER users see — each is a two-account
journey: act as ADMIN, verify as USER.

1. Role create / role-permission assignment (tree-select) → USER's menus and
   buttons change accordingly on their next boot, without re-login where the
   system promises live enforcement.
2. User create + role assignment → the new account can log in and sees
   exactly the assigned scope.
3. User deactivate → their session stops passing authz (or login is
   rejected, per design).
4. Dict / dict-detail edits → the dependent dropdowns/pages reflect new
   entries after reload.

## J11 — Pagination journey (actor: ADMIN)

1. Open any paginated list whose `total` exceeds the default page size →
   first page shows the default row count (10) and the list invariant holds.
2. Change the page-size selector to 20, then 50 → the list API response and
   the rendered table show exactly that many rows (bounded by `total`).
   A selector that changes nothing means the request key is being dropped —
   verify the payload against the frontend API module.
3. Navigate to page 2 → rows are disjoint from page 1; `total` unchanged.
4. Repeat on at least one list per backend module (sysManagement, sysMonitor,
   sysTool) — pagination bugs are per-endpoint, not per-page.

## Coverage map template

```markdown
| Function (from inventory)      | Journey | Status |
|--------------------------------|---------|--------|
| POST /user/create + user page  | J3      | PASS   |
| JIT elevation approve          | J8.3    | PASS   |
| file download                  | J9.2    | FAIL (BUG-<id>) |
| list invariant — /user/list    | J3.7    | PASS   |
| page-size selector (user page) | J11.2   | PASS   |
```
Every inventory row appears exactly once here; gaps are visible by absence.

## Journey template

```markdown
## J<n> — <name>
- **Actor**: ADMIN | USER
- **Preconditions**: <data/accounts needed>
- **Steps**:
  1. <act> → expect <ui assertion> / <api assertion> / <db assertion>
- **Result**: PASS/FAIL, date, evidence links (screenshots, log lines)
```
