## Why

Currently dict and dict_detail logic is combined in single files across proto, logic, server, and gateway handler layers. The original `td27-admin/server` kept them in separate files. Splitting them improves:
- **Maintainability**: Each file has a single responsibility
- **Discoverability**: Easy to find dict_detail-specific code
- **Consistency**: Matches the original architecture and the project's own pattern (other modules have separate files per entity)

## What Changes

- **Proto**: Split `dict.proto` → `dict.proto` (Dict service) + `dict_detail.proto` (DictDetail service)
- **Logic**: Split `dict.go` → `dict.go` (DictLogic) + `dict_detail.go` (DictDetailLogic)
- **Server**: Split `dict.go` → `dict.go` (DictServer) + `dict_detail.go` (DictDetailServer)
- **Handler**: Split `dict.go` → `dict.go` (DictHandler) + `dict_detail.go` (DictDetailHandler)
- Register new `DictDetailServer` gRPC service in `basis.go`
- Add `DictDetailClient` to gateway `ServiceContext`
- Register new gateway handler routes

## Impact

- No functional changes — pure refactoring
- Proto regeneration required
- Routes and handlers remain the same, just renamed
