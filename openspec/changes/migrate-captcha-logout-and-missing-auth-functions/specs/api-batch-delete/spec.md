## ADDED Requirements

### Requirement: Batch delete API entries
The system SHALL support deleting multiple API entries in a single request.

#### Scenario: Successful batch delete
- **WHEN** an authenticated user sends POST /api/apis/delete-by-ids with `{ids: [1, 2, 3]}`
- **AND** all specified API entries exist
- **THEN** all specified API entries SHALL be deleted
- **AND** the response SHALL return HTTP 200 with success

#### Scenario: Batch delete with non-existent IDs
- **WHEN** a user sends POST /api/apis/delete-by-ids with IDs that include non-existent entries
- **THEN** the system SHALL still delete any matching entries
- **AND** SHALL return success (idempotent behavior)
