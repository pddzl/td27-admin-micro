## Why

The gateway has 15 HTTP handlers that define request structs inline and repeat the same 6-line boilerplate for JSON decoding and error responses. There is no centralized validation — each handler manually checks empty fields with `if field == ""`. This creates inconsistency, visual noise, and misses validation for handlers that skip checks entirely.

## What Changes

- Add a `DecodeAndValidate` helper function to `pkg/` that centralizes JSON body decoding + struct validation
- Add `go-playground/validator` as a dependency for struct tag validation (`validate:"required"`)
- Standardize error response format across all handlers (use `pkg.WriteJson` consistently)
- Migrate all 15 inline structs to use the new helper
- No API contract changes — all existing endpoints behave identically

## Capabilities

### New Capabilities
- `request-validation`: Centralized HTTP request body decoding and validation with struct tags for the API gateway

### Modified Capabilities
None — internal refactoring, no spec-level behavior changes.

## Impact

- **New dependency**: `github.com/go-playground/validator/v10`
- **New file**: `pkg/validate.go` — `DecodeAndValidate` helper
- **Modified files**: 7 handler files with 15 endpoint methods (api, button, cron, dept, dict, login, operation_log)
- **Boilerplate removed**: ~90 lines of repetitive JSON decode blocks replaced with single-line calls
