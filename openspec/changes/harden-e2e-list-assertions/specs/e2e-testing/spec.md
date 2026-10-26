## ADDED Requirements

### Requirement: Data-bearing pages must assert content, not just rendering

E2E journeys that land on a page backed by a list MUST assert that visible
rows are present whenever seeded state guarantees rows (e.g. the operation
log always has entries after any authenticated action). A route that renders
without errors and without data MUST NOT pass.

#### Scenario: Operation log page shows rows

- **WHEN** an authenticated admin opens the operation log page after prior
  authenticated actions
- **THEN** the table shows at least one row
- **AND** the row count is consistent with the `total` from the page's list
  API response

### Requirement: Test payloads must match the frontend wire contract

API assertions inside any testing skill MUST send request payloads whose
keys, types, and endpoint shapes match what the frontend actually produces
(derived from `web/src/api/*.ts` or captured network traffic). Hand-built
payloads that diverge from the frontend contract MUST NOT be used as pass
evidence.

#### Scenario: Snake_case page_size replay

- **WHEN** a test calls a paginated list endpoint
- **THEN** the payload uses the frontend's wire keys (e.g. `page_size`, not
  `pageSize`)
- **AND** a key the frontend never sends (or the gateway silently drops)
  fails the run instead of passing

### Requirement: Paginated list responses must satisfy list invariants

Every paginated list response examined by the E2E or API regression skills
MUST satisfy: the `list` field is present; page 1 is non-empty when `total >
0`; `len(list) ≤ requested page_size`; and `total` matches the SQL count
computed in the cross-layer truth phase.

#### Scenario: Empty page one with non-zero total fails

- **WHEN** a list API returns `total: 851` and an empty `list` for page 1
- **THEN** the run records a defect instead of a PASS

### Requirement: Pagination controls must be exercised

An E2E run MUST include a pagination journey that changes the page size
through the UI selector and pages to a second page, asserting row counts and
distinctness.

#### Scenario: Page-size selector takes effect

- **WHEN** the tester selects 20 (or 50) rows per page on a list with
  `total` above that size
- **THEN** the response and the rendered table show exactly that many rows
- **AND** navigating to page 2 shows rows disjoint from page 1 while
  `total` remains unchanged
