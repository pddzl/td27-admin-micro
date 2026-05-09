## ADDED Requirements

### Requirement: Token blocklist interface
The system SHALL provide an abstraction for token invalidation that supports multiple backend implementations.

#### Scenario: Blocklist interface
- **WHEN** a token invalidation backend is implemented
- **THEN** it SHALL implement: `IsBlocklisted(tokenID string) bool` and `AddToBlocklist(tokenID string, expiresAt time.Time) error`

### Requirement: Default in-memory blocklist
The system SHALL provide a default in-memory implementation of the token blocklist.

#### Scenario: In-memory blocklist storage
- **WHEN** a token is logged out
- **THEN** its identifier SHALL be stored in a concurrent-safe in-memory map
- **AND** the entry SHALL include the token's original expiry time

#### Scenario: Automatic cleanup of expired entries
- **WHEN** a blocklist entry's expiry time has passed
- **THEN** the blocklist SHALL consider it valid (not blocklisted) and may remove it on periodic cleanup
- **AND** the system SHALL run periodic cleanup of expired entries (configurable interval, default 1 hour)

### Requirement: Blocklist checking in JWT middleware
The JWT middleware SHALL check the token blocklist before considering a token valid.

#### Scenario: Blocklisted token rejected
- **WHEN** a client presents a token that is in the blocklist
- **THEN** the JWT middleware SHALL return HTTP 401 with "token has been invalidated"

#### Scenario: Non-blocklisted token allowed
- **WHEN** a client presents a valid token that is NOT in the blocklist
- **THEN** the JWT middleware SHALL process the request normally

#### Scenario: Expired blocklisted token
- **WHEN** a client presents a token whose claims show it is expired
- **THEN** the JWT middleware SHALL reject it as "invalid or expired token" regardless of blocklist status

### Requirement: Blocklist storage
The token identifier used for blocklisting SHALL be derived from the JWT.

#### Scenario: Token ID derivation
- **WHEN** a token is added to the blocklist
- **THEN** the identifier SHALL be a SHA-256 hash of the raw token string (to avoid storing the actual token)
