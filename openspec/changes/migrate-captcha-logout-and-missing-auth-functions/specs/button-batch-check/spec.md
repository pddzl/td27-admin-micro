## ADDED Requirements

### Requirement: Batch check button permissions
The system SHALL support checking multiple button permissions in a single request.

#### Scenario: Successful batch check
- **WHEN** an authenticated user sends POST /api/button/batch-check with `{button_codes: ["btn1", "btn2", "btn3"]}`
- **THEN** the system SHALL check each button code against the user's role permissions
- **AND** return a map of `{button_code: boolean}` for each requested code

#### Scenario: Partial permissions
- **WHEN** a user has permission for some but not all requested buttons
- **THEN** the response SHALL include `true` for permitted buttons and `false` for denied buttons
