## ADDED Requirements

### Requirement: Generate CAPTCHA image
The system SHALL provide an endpoint to generate CAPTCHA images for login protection.

#### Scenario: Successful CAPTCHA generation
- **WHEN** a client sends GET /api/captcha
- **THEN** the system returns a CAPTCHA image (base64-encoded) with a unique captcha ID

#### Scenario: CAPTCHA response format
- **WHEN** the system generates a CAPTCHA
- **THEN** the response SHALL include `captcha_id` (string, unique identifier) and `captcha_image` (string, base64-encoded PNG image)
- **AND** the `captcha_id` SHALL be passed back with the login request for verification

### Requirement: Verify CAPTCHA during login
The login endpoint SHALL require captcha verification before authenticating the user.

#### Scenario: Login with valid captcha
- **WHEN** a client sends POST /login with `username`, `password`, `captcha_id`, and `captcha_value`
- **AND** the captcha value matches the stored value for the given captcha ID
- **THEN** the system proceeds with user authentication

#### Scenario: Login with invalid captcha
- **WHEN** a client sends POST /login with an incorrect `captcha_value`
- **THEN** the system SHALL return HTTP 400 with an error message indicating captcha verification failed
- **AND** the system SHALL NOT attempt user authentication

#### Scenario: Login with expired captcha
- **WHEN** a client sends POST /login with a `captcha_id` that has expired or does not exist
- **THEN** the system SHALL return HTTP 400 with an error message indicating captcha expired or invalid

### Requirement: CAPTCHA lifecycle
CAPTCHA codes SHALL expire after a configurable duration.

#### Scenario: Captcha auto-expiry
- **WHEN** a CAPTCHA is generated
- **THEN** it SHALL be automatically invalidated after a configurable TTL (default: 5 minutes)
- **AND** the expired captcha SHALL be removed from the store

### Requirement: CAPTCHA store abstraction
The captcha store SHALL support pluggable backends.

#### Scenario: Default in-memory store
- **WHEN** no external store is configured
- **THEN** the system SHALL use an in-memory store for captcha codes
- **AND** the store SHALL be safe for concurrent access

#### Scenario: Pluggable store interface
- **WHEN** implementing a custom store
- **THEN** it SHALL implement the captcha store interface (Save, Get, Delete, Expire)
