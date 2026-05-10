## ADDED Requirements

### Requirement: Centralized JSON body decoding and validation
The gateway SHALL provide a `DecodeAndValidate` function that decodes a JSON request body into a struct and validates it using struct tags.

#### Scenario: Valid JSON body is decoded successfully
- **WHEN** the request body contains valid JSON matching the target struct
- **THEN** the function SHALL return nil
- **AND** the target struct SHALL be populated with the decoded data

#### Scenario: Malformed JSON body
- **WHEN** the request body contains invalid JSON
- **THEN** the function SHALL return an error describing the parse failure

#### Scenario: Missing required field
- **WHEN** the request body is valid JSON but `validate:"required"` field is missing or empty
- **THEN** the function SHALL return an error naming the missing field

### Requirement: Handlers use consistent error response format
All gateway handlers SHALL use `pkg.WriteJson(w, status, pkg.Error(...))` for error responses instead of raw `w.Header().Set` + `w.WriteHeader` + `json.NewEncoder(w).Encode(...)`.

#### Scenario: Decode error returns 400 with consistent format
- **WHEN** `DecodeAndValidate` returns an error
- **THEN** the handler SHALL respond with HTTP 400
- **AND** the response body SHALL use `pkg.Error(code, msg)` format

### Requirement: Backward compatible behavior
All existing endpoints SHALL behave identically after migration. No new fields, no removed fields, no changed HTTP status codes.

#### Scenario: Existing successful request still succeeds
- **WHEN** a previously working request is sent after migration
- **THEN** the response SHALL be identical in status code and body shape

#### Scenario: Previously unvalidated requests remain unvalidated
- **WHEN** an inline struct has no `validate` tags added
- **THEN** `DecodeAndValidate` SHALL skip validation for that struct
