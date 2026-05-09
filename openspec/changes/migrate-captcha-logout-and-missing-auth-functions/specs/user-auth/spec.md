## ADDED Requirements

### Requirement: Captcha-protected login flow
The login flow SHALL incorporate CAPTCHA verification.

#### Scenario: Login with captcha
- **WHEN** a client sends POST /login with `username`, `password`, `captcha_id`, and `captcha_value`
- **AND** captcha verification succeeds
- **AND** credentials are valid
- **THEN** the system SHALL return HTTP 200 with `token` (real JWT), `user` object

#### Scenario: Login without captcha fields
- **WHEN** a client sends POST /login without `captcha_id` and `captcha_value`
- **THEN** the system SHALL return HTTP 400 with an error message

### Requirement: Real JWT in login response
The login response SHALL contain a real JWT token instead of a placeholder.

#### Scenario: Login response includes real JWT
- **WHEN** login succeeds
- **THEN** the `token` field in the response SHALL be a real JWT signed with the configured signing key
- **AND** the JWT SHALL be parseable and validatable by the gateway's JWT middleware

### Requirement: Login request validation
The login endpoint SHALL validate all required fields.

#### Scenario: Missing username or password
- **WHEN** a client sends POST /login without `username` or `password`
- **THEN** the system SHALL return HTTP 400 with "invalid request body"

#### Scenario: Empty username or password
- **WHEN** a client sends POST /login with empty `username` or `password`
- **THEN** the system SHALL return HTTP 400 with validation error
