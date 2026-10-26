## ADDED Requirements

### Requirement: List filters must be forwarded end-to-end

Every list endpoint whose frontend sends filter keys MUST accept, forward,
and apply them: gateway handler MUST parse the exact wire keys the frontend
API module sends, the rpc request MUST carry them, and the repository MUST
apply them to both the count and data queries. A filter key the frontend
sends that the backend silently drops is a defect.

#### Scenario: Method filter narrows /api/list

- **WHEN** the API page sends `POST /api/list` with
  `{"method": "GET", "page": 1, "page_size": 10}` (exact frontend payload)
- **THEN** every returned row has `method == "GET"`
- **AND** `total` equals the SQL count of non-deleted rows with
  `method = 'GET'`

#### Scenario: Text filters partially match

- **WHEN** the API page sends `path`, `group_en`, or `description` filter
  values
- **THEN** returned rows match the filter (partial match for text inputs,
  exact match for the method dropdown)
- **AND** clearing the filters restores the unfiltered list and total
