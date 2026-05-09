## ADDED Requirements

### Requirement: Batch delete cron entries
The system SHALL support deleting multiple cron entries in a single request.

#### Scenario: Successful batch delete
- **WHEN** an authenticated user sends POST /api/cron/delete-by-ids with `{ids: [1, 2, 3]}`
- **AND** all specified cron entries exist
- **THEN** all specified cron entries SHALL be deleted
- **AND** their associated cron schedulers SHALL be stopped
- **AND** the response SHALL return HTTP 200 with success

#### Scenario: Batch delete stops cron schedulers
- **WHEN** cron entries are batch deleted
- **THEN** the system SHALL stop the running cron scheduler for each deleted entry
