## ADDED Requirements

### Requirement: Single delete operation log
The system SHALL support deleting a single operation log entry.

#### Scenario: Successful single delete
- **WHEN** an authenticated user sends POST /api/operation-log/delete with `{id: 123}`
- **AND** the log entry exists
- **THEN** the entry SHALL be deleted
- **AND** return HTTP 200 with success

#### Scenario: Delete non-existent log
- **WHEN** a user sends POST /api/operation-log/delete with a non-existent id
- **THEN** the system SHALL return success (idempotent)

### Requirement: Batch delete operation logs
The system SHALL support deleting multiple operation log entries in a single request.

#### Scenario: Successful batch delete
- **WHEN** an authenticated user sends POST /api/operation-log/delete-by-ids with `{ids: [1, 2, 3]}`
- **THEN** all specified log entries SHALL be deleted
- **AND** return HTTP 200 with success
