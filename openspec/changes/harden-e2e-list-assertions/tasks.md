# Tasks

## 1. e2e-testing skill

- [x] 1.1 SKILL.md ground rules: data-bearing pages assert visible content;
      API assertions replay frontend-derived payloads (`web/src/api/*.ts` or
      captured network), never hand-built JSON.
- [x] 1.2 SKILL.md: add the reusable list-invariant assertion block (list
      present; page 1 non-empty when total > 0; len(list) ≤ page_size;
      total vs SQL count via Phase 4).
- [x] 1.3 SKILL.md browser pitfalls: "empty table while DB has rows is a
      defect" entry.

## 2. Flow catalog

- [x] 2.1 J3: reference the list-invariant block after create/delete
      refresh steps.
- [x] 2.2 J7: upgrade "page renders" to "renders with data where the DB
      guarantees rows" (operation log example).
- [x] 2.3 Add J11 — Pagination journey (page-size 10/20/50, page 2
      distinctness, total stability).
- [x] 2.4 Update the coverage map template with list-invariant and
      pagination rows per paginated resource.

## 3. api-regression-testing skill

- [x] 3.1 Phase 3 matrix: apply the list-invariant block to list endpoints;
      add the payload-source-of-truth rule.

## 5. Interaction journeys (added 2026-10-26)

- [x] 5.1 SKILL.md ground rule: client-side interactions get a journey step
      that flips visible state both directions, twice; mount-time-only props
      called out (tree `default-expand-all`).
- [x] 5.2 SKILL.md browser pitfall: "toggle control exists but changes
      nothing is a defect".
- [x] 5.3 Flow catalog: J12 — Tree expand/collapse journey (dept tree),
      coverage-map row included.

## 4. Verification

- [x] 4.1 Docs self-check: all rules present, J11 defined, coverage map
      updated; no product code touched (`git diff --stat` shows only
      skills/ and openspec/).
- [ ] 4.2 Next E2E + API regression run adopts the new rules and passes
      green; archive this change afterward.
