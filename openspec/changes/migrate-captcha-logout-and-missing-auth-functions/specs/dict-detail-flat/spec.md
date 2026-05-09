## ADDED Requirements

### Requirement: Flat list of dict details
The system SHALL support retrieving all dict details for a given dict ID without pagination.

#### Scenario: Successful flat list
- **WHEN** an authenticated user sends POST /api/dict/detail/flat with `{dict_id: 5}`
- **AND** the dict exists
- **THEN** the system SHALL return all dict details for that dict
- **AND** the response SHALL NOT be paginated

#### Scenario: Flat list with non-existent dict
- **WHEN** a user sends POST /api/dict/detail/flat with a non-existent dict_id
- **THEN** the system SHALL return an empty list
