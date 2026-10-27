# web-api-contracts Specification

## Purpose
TBD - created by archiving change unify-operation-log-uri. Update Purpose after archive.

## Requirements

### Requirement: One canonical URI per resource

Each gateway resource MUST be reachable under exactly one canonical URI
pattern (`/<resource>/<action>`); legacy aliases MUST NOT be registered.
The operation-log resource canonical base is `/operation-log`.

#### Scenario: Alias routes are gone

- **WHEN** a client calls `POST /opl/list` (or `/opl/delete`,
  `/opl/delete-by-ids`)
- **THEN** the gateway returns 404 instead of silently serving the resource
  under a second name

#### Scenario: Canonical routes serve the frontend

- **WHEN** the operation-log page lists rows or performs delete /
  delete-by-ids
- **THEN** the requests target `/operation-log/*` and succeed with no 403
  or 404 (permission rows migrated in step)

### Requirement: Permission rows must match registered URIs

Every registered route MUST have matching `sys_management_api` /
`sys_management_permission` (domain `'api'`) rows for its canonical path,
and no rows for unregistered aliases. Migrations renaming paths MUST be
idempotent and MUST NOT create duplicate rows.

#### Scenario: No stale /opl permission rows

- **WHEN** migration 011 has run twice
- **THEN** no `sys_management_api` or `sys_management_permission` row
  references `/opl/*`
- **AND** no duplicate `/operation-log/*` rows were created

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
