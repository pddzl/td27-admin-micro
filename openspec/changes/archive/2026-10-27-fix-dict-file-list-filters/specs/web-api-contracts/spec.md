## ADDED Requirements

### Requirement: Dict and file list filters must be forwarded end-to-end

The dict list (`POST /dict/list`) MUST accept, forward, and apply
`cn_name`/`en_name`; the file list (`POST /file/list`) MUST accept, forward,
and apply `name`. The dict page MUST include the current search-form values
in the request body, and both backends MUST apply the filters to the count
and data queries with partial matching.

#### Scenario: Dict Chinese-name filter narrows results

- **WHEN** the dict page sends `POST /dict/list` with
  `{"cn_name": "API", "page": 1, "page_size": 10}` (exact frontend payload)
- **THEN** every returned row's `cn_name` contains the filter value
- **AND** `total` equals the count of matching non-deleted rows

#### Scenario: File name filter narrows results

- **WHEN** the file page sends `name` to `POST /file/list`
- **THEN** every returned row's `file_name` contains the filter value
- **AND** clearing the filter restores the unfiltered list and total
