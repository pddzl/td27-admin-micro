---
name: e2e-testing
description: Drive end-to-end user journeys through browser, backend, and database for web admin systems. Use when asked to run E2E tests, verify a full user flow in the UI, or validate that pages show real data. Do not use for API-only regression (use api-regression-testing) or unit tests.
---

# E2E Testing

Execute real user journeys through the full stack — browser UI, backend
services, and database — and assert user-visible behavior. Complements the
API-level sibling skill `api-regression-testing`: run that first to prove the
API layer, then this skill to prove the product a user actually experiences.

## Ground rules

- Never guess. Ask for credentials, URLs, and whether browser automation is
  authorized before launching it.
- Every defect gets a bug-log entry (repro, root cause, fix, retest) and is
  never closed without a passing retest — same discipline as the sibling skill.
- Clean up rows created during a journey through the product's own delete UI
  or API, then verify removal in the database.
- One journey at a time; a failed step stops that journey only.

## Phase 0 — Parameterize

Collect before starting: frontend URL, backend API base, admin credentials, a
limited-permission test account, and database access. Confirm the stack is up
(health endpoints, login page renders). Reuse the route/page inventory from
`api-regression-testing` Phase 1 if it exists.

## Phase 1 — Flow inventory

Enumerate the journeys worth driving (catalog in
`references/flow-catalog.md`). Source them from the page list and route
inventory, and let the user pick priorities. Each journey gets: actor (which
account), steps, and the assertions at every layer.

**Coverage invariant**: every function in the inventory (route, page, or
feature) must map to at least one journey — record the mapping explicitly
(journey → functions covered). Anything unmapped is a gap: either fold it
into an existing journey or write a new one. "Every function is tested" is
then a checklist you can verify, not a feeling.

## Phase 2 — Session bootstrap

Two patterns; pick per journey:

- **Real login** — mandatory when login itself is the subject (wrong password,
  captcha rejection, lockout). For image captchas: fetch the captcha endpoint
  via curl, decode the image, read the code (image viewer/OCR), submit within
  its validity window.
- **Token injection** — when login is NOT the subject, inject a valid session
  token directly into whatever storage the app reads (find the key in the
  frontend auth store source, e.g. a localStorage token entry), then reload.
  Deterministic and fast; skips captcha flakiness.

Verify the session landed: the app renders the authenticated shell (not a
login redirect) before asserting anything else.

## Phase 3 — Execute journeys

For each journey step: act through the UI (click/fill — not API calls), then
assert at the layer where truth lives:

- UI assertion: element present/absent, text, computed state (a button can be
  in the DOM yet hidden — check `display`, not just existence).
- API assertion: the network request fired and returned the expected code.
- DB assertion: the row exists/changed/vanished.

## Phase 4 — Cross-layer truth

UI data must match the database. For any page displaying counts, sums, or
lists, compute the expected value via SQL (mind soft-delete filters and time
windows) and compare field by field. A plausible-looking number is not
evidence — this check has already caught "wrong dashboard statistics" that
were actually transport-layer failures showing empty defaults.

## Phase 5 — Bug loop

Record first, fix second, retest third — the exact repro plus one line of
evidence per layer. Re-run the affected journey after every fix; the pass is
green when one complete run of all journeys finishes with zero open bugs.

## Browser pitfalls (learned the hard way)

- Async races: a snapshot taken immediately after an action may predate the
  UI settling (a permission directive hides buttons a beat later). Re-check
  after a short wait or inspect computed styles before declaring PASS/FAIL.
- Teleported/floating components: collapsed submenu items render in a popup
  attached to `body`, not inside the sidebar DOM — search the whole document
  or hover to expand before concluding "empty".
- SPA navigation: hash routes may not trigger a full load; navigate explicitly
  and wait for the route's content, not just URL change.
- Token injection bypasses login flows by design — do not use it to test
  login, logout, or session-expiry journeys.
- A login redirect during an authenticated journey means a boot API call was
  rejected: check the router guard's boot calls (user info, menus, buttons)
  against the account's grants before suspecting the UI.

## Evidence format

One line per step: `<journey>.<step> <actor> → <observed> (layer: ui|api|db)`.
Keep screenshots/console/network excerpts when a step fails; they are the
repro.
