# Harden E2E List Assertions Against Transport-Shape Bugs

## Why

The 2026-10-24 E2E run (J1–J10, commit `02a6cc35`) passed green while
`POST /api/opl/list` returned `total: 851` with an **empty** list — found by
hand two days later. Three procedural gaps caused the miss:

1. **J7 asserted render-only**: `#/sysMonitor/operationLog` "renders" was a
   PASS even with zero rows; an empty table produces no error.
2. **API journeys hand-built payloads**: the bug fires only with the
   frontend's real wire key `page_size` (snake_case, set by
   `fix-frontend-key-mapping`); hand-written camelCase `pageSize` decoded
   fine, so the mismatch was never exercised. Same failure class as BUG-001
   ("re-verify the CURRENT frontend contract").
3. **No list-response invariants existed**: nothing asserted `list` is
   present, page 1 is non-empty when `total > 0`, or `len(list) ≤ page_size`.

These are procedure gaps in the testing skills, not code gaps — fixing them
is how this bug family gets caught on every future run.

## What Changes

Docs-only changes to the two test skills:

- `skills/e2e-testing/SKILL.md`:
  - Ground rule: data-bearing pages MUST be asserted for actual content
    (visible rows), never render-only.
  - Ground rule: API assertions MUST replay the exact payload shape the
    frontend produces — derive from `web/src/api/*.ts` (or captured network
    traffic), never hand-built JSON.
  - A reusable **list-invariant assertion block** for every paginated list:
    `list` key present; page 1 non-empty when `total > 0`; `len(list) ≤
    page_size`; `total` matches the SQL count (Phase 4 cross-check).
  - Browser pitfall: an empty table while the DB holds rows is a defect, not
    a rendering quirk.
- `skills/e2e-testing/references/flow-catalog.md`:
  - J3: run the list-invariant block after create/delete refreshes.
  - J7: "page renders" → "page renders **with data** where the DB
    guarantees rows" (the operation log always has rows after any
    authenticated action).
  - New **J11 — Pagination journey**: page-size selector (10/20/50) changes
    the returned row count; page 2 shows different rows; `total` stays
    constant.
  - Coverage map: add list-invariant and pagination rows per paginated
    resource.
- `skills/api-regression-testing/SKILL.md`: same list-invariant block and
  payload-source-of-truth rule applied to the Phase 3 matrix for list
  endpoints.

## Impact

- No product, frontend, proto, or DB code changes — procedure/docs only.
- Effect: the /opl/list bug class (transport key drop + silently empty list)
  is caught by J7, J11, and the API matrix on any list endpoint.
- Next E2E/API regression run exercises the new rules; archive this change
  only after a green run that used them.
