## Context

Three goctl-generated scaffolding files (`pinglogic.go`, `basisserver.go`, `client/basis.go`) reference `td27/rpc/basis/types/basis` which was never generated. The entry point `basis.go` registers 13 service-specific servers but never the Basis/Ping server.

## Goals / Non-Goals

**Goals:**
- Remove all files referencing the missing `types/basis` package
- Fix `go build ./...` for the entire module
- Clean up empty directories

**Non-Goals:**
- No code changes beyond deletion
- No functional changes

## Decisions

Simple deletion of 4 files + cleanup of empty dirs. No dependencies involved.

## Risks / Trade-offs

None — the Basis/Ping RPC was never registered or called anywhere.
