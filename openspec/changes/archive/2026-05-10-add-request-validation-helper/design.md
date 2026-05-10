## Context

The gateway has 15 HTTP handler methods that independently decode JSON bodies using identical boilerplate:
```go
var req struct { ... }
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(pkg.Error(400, "invalid request body"))
    return
}
```

Some handlers use `pkg.WriteJson(w, ...)` for errors, others use raw `w.Header().Set + w.WriteHeader + json.NewEncoder(w).Encode(...)`. Only the login handler checks for empty required fields.

## Goals / Non-Goals

**Goals:**
- Centralize JSON decode + validation into a single `pkg.DecodeAndValidate(r, &req)` call
- Enable struct-tag validation (`validate:"required"`) on gateway request structs
- Standardize all 15 handlers to use the same error response format
- Keep inline structs where they're single-use (don't force named types)

**Non-Goals:**
- No behavior change for any endpoint
- No migration to named request types (can be done later)
- No validation at the gRPC service level (that's already handled)
- No changes to GET endpoints that use query/path parameters

## Decisions

### 1. Use `go-playground/validator` over custom validation

`go-playground/validator` is the de-facto Go validation library (11k+ GitHub stars), already used by Gin internally. Writing custom validation for every struct shape would recreate this wheel.

### 2. `DecodeAndValidate` signature

```go
// DecodeAndValidate decodes JSON body and validates struct tags.
// Returns a user-friendly error message on failure.
// Usage: pkg.DecodeAndValidate(r.Body, &req)
func DecodeAndValidate(body io.ReadCloser, req interface{}) error
```

Single return value (error) so handlers can do:
```go
var req struct { ... }
if err := pkg.DecodeAndValidate(r.Body, &req); err != nil {
    api.FailWithRequest(w, http.StatusBadRequest, err.Error())
    return
}
```

### 3. Validator is initialized once as a package-level singleton

```go
var validate = validator.New()
```

This is the standard pattern — `validator.New()` is cheap but not free, and reusing the instance allows caching validation rules.

### 4. Keep inline structs — don't force named types

The 15 request shapes are all used once. Moving them to a `types/` package adds boilerplate without benefit. The inline struct pattern is fine — the missing piece is centralized validation, not named types.

### 5. Standardize all handlers on `pkg.WriteJson`

Some handlers use `pkg.WriteJson(...)`, others use raw `w.Header().Set(...)`. Migrate all 15 to `pkg.WriteJson(...)` for consistency.

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| New dependency in go.mod | `go-playground/validator/v10` is stable, widely used, has no transitive dependency concerns. |
| Existing structs lack `validate` tags | Start with just the login handler (captcha fields) and other handlers that need required-field enforcement. Empty struct tags = no validation = backward compatible. |
| Flagged as "implementation" | This is purely a refactoring of existing behavior. No new features, no API changes. |
