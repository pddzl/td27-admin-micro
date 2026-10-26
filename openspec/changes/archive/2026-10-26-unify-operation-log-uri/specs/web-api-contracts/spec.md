## ADDED Requirements

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
