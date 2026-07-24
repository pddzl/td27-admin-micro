## Why

Dict create/delete currently call `getTableData()` which makes an unnecessary API call. The original pattern used `push()`/`splice()` for immediate table updates. `push()` is better UX — instant feedback without a loading state.

The issue was:
- `dict/list` handler reads pagination from URL query params, not POST body → list never loads → `tableData.value.list` is undefined → `push()` crashes

## What Changes

- Fix `dict/list` handler to read page/pageSize from POST body
- Fix dict create: backend should return the created dict so `res.data` is meaningful
- Fix dict delete: use `splice()` instead of `getTableData()`
- Add `RETURNING id` to dict INSERT so the created ID is available

## Impact

- No more unnecessary `getTableData()` calls after create/delete
- Table updates instantly without loading spinner
