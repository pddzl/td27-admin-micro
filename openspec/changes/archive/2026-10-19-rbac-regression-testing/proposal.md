## Why

RBAC is the security core of the admin system (users are engineers managing pg-backed data), but it has never been regression-tested end-to-end. A regression loop is needed to prove: role/permission CRUD, API enforcement (401/403, deny-wins, live revocation), button permissions, JIT elevations, and the auth lifecycle all work together. Bugs found must be recorded with repro + fix + retest evidence (留痕).

## What Changes

- Add a structured RBAC regression test matrix (P0 environment → P5 auth lifecycle) executed against the local dev stack via curl (API) and browser (UI).
- Flip dev `EnforceRbac: false → true` in `api/gateway/etc/gateway.yaml` so API enforcement is actually exercised (approved by user).
- Record every bug in `bugs.md` with repro steps, root cause, fix, and retest result; loop until a full green pass.
- Bug #1 (already observed): 操作日志 renders as an orphan top-level sidebar item instead of nesting under 系统监控.

## Non-goals

- No new RBAC features; no spec/capability behavior changes.
- No changes to k8s manifests or release workflow.
- Frontend visual polish beyond what bugs require.
