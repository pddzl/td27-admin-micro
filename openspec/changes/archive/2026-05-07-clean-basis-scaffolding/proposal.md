## Why

Three goctl-generated scaffolding files reference `td27/rpc/basis/types/basis`, a package that was never generated and does not exist. These files block `go build ./...` and are entirely unused — the entry point `basis.go` never registers the Basis/Ping server. They are dead code from the initial goctl scaffolding.

## What Changes

- Delete `rpc/basis/internal/logic/basis/pinglogic.go`
- Delete `rpc/basis/internal/server/basis/basisserver.go`
- Delete `rpc/basis/client/basis/basis.go`
- Delete `rpc/basis/basis.proto` (its only RPC was Ping, never registered)
- Remove empty directories: `rpc/basis/internal/logic/basis/`, `rpc/basis/internal/server/basis/`, `rpc/basis/client/basis/`, `rpc/basis/client/`

## Capabilities

### New Capabilities
(No new capabilities — this is a cleanup)

### Modified Capabilities
(No capability changes — purely removing dead code)

## Impact

- Removes 3 Go files with broken imports (~80 lines total)
- Removes 1 unused proto file
- Cleans up 4 empty directories
- Fixes `go build ./...` for the entire module
- No functional impact — Basis/Ping was never registered or called
