# Tasks

## 1. Audit

- [x] 1.1 Audited live DB (2026-10-26): BOTH row sets existed live
      (`/opl/*` 3 api + 3 perm, `/operation-log/*` 4 api + 4 perm, all
      `deleted_at IS NULL`) — so 011 retires the alias rows instead of
      renaming.

## 2. Code

- [x] 2.1 Deleted the three `/opl/*` alias routes from handler.go.
- [x] 2.2 Switched the three URLs in web/src/api/sysMonitor/operationLog.ts
      to `/operation-log/*`.

## 3. Migration

- [x] 3.1 Wrote and applied `scripts/011_unify_operation_log_uri.sql`
      (idempotent): soft-deleted 3 alias perm rows + their role grants +
      3 alias api-registry rows; canonical rows untouched.

## 4. Verification

- [x] 4.1 `go vet ./api/gateway/...` and `go build ./api/gateway/...` pass.
- [x] 4.2 `grep` for `/opl/` returns no live-code references in web/src,
      api/, or web/tests.
- [x] 4.3 Live (2026-10-26, running gateway after restart): canonical
      `/operation-log/list` returns code 0 with rows (page_size honored);
      `POST /opl/list` returns 404; permission rows migrated earlier, no
      403. Cleanup/delete routes share the same handler wiring verified by
      build + code review (destructive calls not executed live).
- [x] 4.4 Archived once 4.1–4.3 were green.
