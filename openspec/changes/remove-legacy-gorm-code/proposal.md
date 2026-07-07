## Why

The GORM-to-sqlx conversion of the active repository layer is complete, but 8 dead files still reference `gorm.io/gorm`. These are legacy "authority" table models/repos that were replaced by the `sys_management_*` table models. They're compiled (same package) but never called.

Removing them:
- Eliminates the dependency on `gorm.io/gorm` from go.mod
- Removes dead code that could confuse future developers
- Verifies the sqlx conversion is complete

## What Changes

- Delete 8 dead files from `rpc/basis/internal/model/`
- Remove `gorm.io/gorm` from `go.mod`
- Run `go mod tidy` to clean up stale indirect deps
- Verify both binaries still build and vet clean

## Impact

No functional impact — these files are dead code (0 references across the codebase).
