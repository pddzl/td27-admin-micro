## ADDED Requirements

### Requirement: Dashboard statistics
The system SHALL provide aggregate statistics for the admin dashboard.

#### Scenario: Get dashboard statistics
- **WHEN** an authenticated user sends GET /api/dashboard/statistics
- **THEN** the system SHALL return aggregate counts including:
  - Total number of users
  - Total number of roles
  - Total number of APIs
  - Total number of departments
  - Total number of menus
  - Total number of dict entries
  - Total number of operation logs (recent period)
  - Total number of files

### Requirement: Dashboard recent operations
The system SHALL provide the most recent operation log entries.

#### Scenario: Get recent operations
- **WHEN** an authenticated user sends GET /api/dashboard/recent-operations
- **THEN** the system SHALL return the 10 most recent operation log entries
- **AND** each entry SHALL include: user name, method, path, status, timestamp

### Requirement: Dashboard system info
The system SHALL provide runtime system information.

#### Scenario: Get system info
- **WHEN** an authenticated user sends GET /api/dashboard/system-info
- **THEN** the system SHALL return:
  - Go version
  - Operating system
  - Architecture
  - Number of CPU cores
  - Server uptime (process start time)
  - Goroutine count
  - Memory stats (allocated, total allocated, garbage collections)
