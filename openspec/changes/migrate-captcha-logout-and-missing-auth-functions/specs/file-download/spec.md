## ADDED Requirements

### Requirement: File download endpoint
The system SHALL support downloading uploaded file binary content.

#### Scenario: Successful file download
- **WHEN** an authenticated user sends GET /api/file/download/:id with a valid file ID
- **AND** the file exists on disk
- **THEN** the system SHALL return the file binary content
- **AND** set `Content-Type` based on the file's MIME type
- **AND** set `Content-Disposition: attachment; filename="..."` header
- **AND** return HTTP 200

#### Scenario: File not found on disk
- **WHEN** a user sends GET /api/file/download/:id with a valid file ID
- **BUT** the file does not exist on the configured upload path
- **THEN** the system SHALL return HTTP 404 with an error message

#### Scenario: File ID not found
- **WHEN** a user sends GET /api/file/download/:id with a non-existent file ID
- **THEN** the system SHALL return HTTP 404

#### Scenario: Path traversal prevention
- **WHEN** a file path is resolved from the database
- **THEN** the system SHALL clean and validate the path
- **AND** reject paths that escape the configured upload directory (contain `..`)
