## Why

Three endpoints from the original `td27-admin/server` were not carried over to the go-zero rewrite:
- `dept/list` — flat filterable department list (for dropdowns/selectors)
- `dictDetail/list` — paginated dict detail list
- `api/elTree` — role-API permission tree for UI assignment

## What Changes

- **Dept list**: Add `GET /api/dept/list` endpoint returning paginated flat department list
- **DictDetail list**: Add `POST /api/dict/detail/list` endpoint returning paginated dict details
- **API elTree**: Add `POST /api/apis/tree` endpoint returning API tree with checked state for a role

## Impact

- Proto additions to dept, dict, and api services
- New handlers in gateway
- New routes registered
