# e2e-testing Specification

## Purpose
TBD - created by archiving change fix-api-list-filters. Update Purpose after archive.

## Requirements

### Requirement: Search filters must be exercised and asserted

E2E journeys on pages with search/filter forms MUST drive each filter
through the UI (or replay the exact frontend payload at the API layer) and
assert that every returned row matches the filter — visible rows, response
`total`, and the filtered SQL count must agree. A filter control that
changes nothing is a defect, even when the unfiltered list is correct.

#### Scenario: Dead filter is caught

- **WHEN** the tester selects a filter value on a list page (e.g.
  `method: GET` on the API page)
- **THEN** the run asserts every visible row matches the filter
- **AND** a filter that leaves the rows unchanged when filtered rows exist
  is recorded as a defect, not a rendering quirk
