
## 2026-05-08: Basis proto scaffolding

- Proto files for the main service go at `rpc/basis/basis.proto`, sub-domain protos go under `rpc/basis/proto/<category>/`.
- Generated types go to `rpc/basis/types/<name>_pb/` with `_pb` suffix in the directory name.
- go_package must include the full module prefix (e.g., `td27/rpc/basis/types/basis_pb;basis_pb`) for `--go_opt=module=td27` to work properly with protoc-gen-go.
- The server handler pattern: one file per service in `internal/server/<service>/` wrapping logic from `internal/logic/<service>/`.
- `go.build ./...` passes with zero errors after implementation.
