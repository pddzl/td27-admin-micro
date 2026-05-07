## ADDED Requirements

All requirements for user management SHALL exactly match the implementation in `~/Mine/td27-admin/server/internal/model/sysManagement/user*.go` and related service/api code.

### Requirement: User CRUD operations
The system SHALL support full create, read, update, delete operations for user entities.

#### Scenario: Create new user
- **WHEN**: Valid user creation request is received via gRPC
- **THEN**: User is saved to database, password is hashed using bcrypt, success response is returned

#### Scenario: Get user by ID
- **WHEN**: Request to get user by ID is received
- **THEN**: User details are returned without sensitive information (password hash excluded)
