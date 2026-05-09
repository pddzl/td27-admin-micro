## ADDED Requirements

### Requirement: JWT token generation on login
The system SHALL generate a real JWT token when a user logs in, replacing the current placeholder token.

#### Scenario: Successful login returns real JWT
- **WHEN** a user successfully authenticates with valid credentials
- **THEN** the system SHALL generate a JWT token signed with the configured `SigningKey`
- **AND** the token SHALL NOT be the string "placeholder_token"

#### Scenario: JWT token claims
- **WHEN** a JWT token is generated
- **THEN** it SHALL contain the following claims:
  - `userId` (uint64): The authenticated user's ID
  - `username` (string): The authenticated user's username
  - `roleIds` ([]uint64): The user's assigned role IDs
  - `exp` (int64): Token expiration time (Unix timestamp)
  - `iat` (int64): Token issued at time (Unix timestamp)
  - `iss` (string): Token issuer (from config)

#### Scenario: Token expiration
- **WHEN** a JWT token is generated
- **THEN** its `exp` claim SHALL be set to `current_time + expires_time` from config
- **AND** the gateway JWT middleware SHALL reject tokens past their expiration time

### Requirement: JWT signing key consistency
The JWT signing key SHALL be consistent across rpc/basis (token generation) and api/gateway (token validation).

#### Scenario: Matching signing keys
- **WHEN** rpc/basis generates a token with `JWT.SigningKey`
- **THEN** api/gateway SHALL validate it using `Auth.AccessSecret` (both currently configured as `pddzl`)
- **AND** if the keys differ, token validation SHALL fail

### Requirement: JWTManager.GenerateToken method
The JWTManager struct in rpc/basis SHALL have a `GenerateToken` method.

#### Scenario: GenerateToken method signature
- **WHEN** `GenerateToken` is called with `userId uint64`, `username string`, `roleIds []uint64`
- **THEN** it SHALL return a signed JWT string
- **AND** it SHALL return an error if token generation fails

#### Scenario: JWTManager uses configured values
- **WHEN** `GenerateToken` creates a token
- **THEN** it SHALL use `JWTManager.SigningKey` for HMAC signing
- **AND** it SHALL use `JWTManager.ExpiresTime` for token TTL
- **AND** it SHALL use `JWTManager.Issuer` for the `iss` claim
