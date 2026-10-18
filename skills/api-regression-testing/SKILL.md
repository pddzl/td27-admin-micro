---
name: api-regression-testing
description: Drive systematic API regression testing for backend admin systems. Use when asked to regression-test APIs, verify RBAC or permissions end to end, or inventory all registered endpoints. Do not use for unit tests or frontend-only work.
---

# API Regression Testing

Run a closed-loop regression pass over a backend admin system: build a complete
route/function inventory from source, diff it against the live permission
registry, execute an enforcement/auth/CRUD matrix, and log every defect with
repro, root cause, fix, and a passing retest before closing it.

## Ground rules

- Never guess. Ask the user for anything non-deterministic: credentials, base
  URL, whether to flip enforcement flags, whether an incident counts as a bug.
- Record everything in a bug log (default `bugs.md` in the change/test folder):
  one entry per defect with ID, status, found-date, repro, root cause, fix,
  retest evidence. A bug is never closed without a passing retest. Non-bug
  findings go in an Observations section, not silently dropped.
- Prefer curl/HTTP for evidence. Use a browser only when the subject is UI
  behavior, and ask before launching browser automation.
- Clean up test rows through the product's own delete APIs when possible; hard
  verify removal in the DB afterwards.

## Phase 0 — Parameterize

Collect before starting (ask when absent): base URL, auth token header name,
login flow (captcha? body shape?), a superadmin principal plus one
limited-scope test principal, DB access (dsn or exec command), the route
registration file(s), and where permissions/policies live (tables, files).
State the mapping from HTTP method/path to permission action if the system has
one — the matrix depends on it.

## Phase 1 — Inventory

1. Extract every registered route: run `scripts/extract_routes.py` against the
   handler/route registration source (see `--help`). It emits JSON with method,
   path, and middleware/auth wrapper per route, and can also list gRPC services
   from proto files.
2. Extract the permission registry (e.g. `SELECT resource, action, effect FROM
   permission_table WHERE deleted_at IS NULL`) and diff both directions:
   - registry rows with no registered route (orphan grants — never match), and
   - registered routes with no registry row (ungrantable — everyone gets 403).
3. Recompute each row's expected action from the method→action mapping and flag
   mismatches. Every drift item is a bug candidate; confirm one live (expect a
   403 on a route the registry should allow) before writing it up.

## Phase 2 — Environment gate

Verify migrations applied, services healthy (`/health` or equivalent), and the
enforcement flag actually on. Confirm auth end to end (login, token accepted,
expired token rejected 401). Secure the two principals; store tokens in files
with short lives and re-login when they expire instead of debugging phantom
401s.

## Phase 3 — Matrix

Execute `references/test-matrix.md` top to bottom, skipping sections the
project lacks (say so explicitly when skipping). Record one line of evidence
(per method + status + timestamp) per case. Do not batch failures — a red case
stops that section only; continue independent sections and loop back after
fixes.

## Phase 4 — Bug loop

For each defect: write the bug entry first, then fix, then retest the exact
repro, then re-run the full matrix (at minimum every section the fix could
touch). The pass is green only when one complete run finishes with zero open
bugs. Finish with a summary: bugs found/fixed, observations, leftover risks.

## Diagnostics

- Sudden blanket 403s with valid grants: check whether the permission/authz RPC
  is reachable. Fail-closed designs deny everything when it is down — that is
  enforcement working, not a policy bug. Confirm via server logs before
  touching policy data.
- Distinguish HTTP-layer mismatches (405 wrong method, 404 wrong path) from
  authz 403s; fix the test before suspecting the system.
- JSON bodies for proto-backed endpoints usually need numeric enum values, not
  string names.
- Long-lived sessions on laptops: VPN/route churn can wedge idle service
  processes. A process that accepts TCP but never answers is wedged — restart
  it instead of debugging the client.
