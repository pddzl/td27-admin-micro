## ADDED Requirements

### Requirement: Logout endpoint
The system SHALL provide an authenticated endpoint to terminate user sessions.

#### Scenario: Successful logout
- **WHEN** an authenticated user sends POST /api/logout with a valid Bearer token
- **THEN** the system SHALL add the token to the blocklist
- **AND** return HTTP 200 with a success message

#### Scenario: Logout without token
- **WHEN** an unauthenticated client sends POST /api/logout without an Authorization header
- **THEN** the system SHALL return HTTP 401

#### Scenario: Logout with invalid token
- **WHEN** a client sends POST /api/logout with an expired or invalid Bearer token
- **THEN** the system SHALL return HTTP 401

### Requirement: Token invalidation on logout
The system SHALL prevent a logged-out token from being used for subsequent requests.

#### Scenario: Protected endpoint after logout
- **WHEN** a client that has logged out reuses the same token to access a protected endpoint
- **THEN** the system SHALL return HTTP 401

#### Scenario: Logout idempotency
- **WHEN** a client sends POST /api/logout with a token that has already been logged out
- **THEN** the system SHALL return HTTP 200 (idempotent — token is already invalid)

### Requirement: Logout does not require operation record
The logout endpoint SHALL NOT record an operation log entry.

#### Scenario: Logout bypasses operation log middleware
- **WHEN** a logout request is processed
- **THEN** it SHALL NOT be logged to the operation log (the token is being invalidated, making audit logging unreliable)
