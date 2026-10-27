# Fix Dead Search Filters on /dict/list and /file/list

## Why

Sweep of the remaining dead search filters (same family as
`fix-api-list-filters`). Three pages have search forms; only the API page
works end-to-end now:

- **dict** (字/典管理): the search form holds `cn_name`/`en_name`, but the
  page's `getTableData()` calls `dictListApi({page, pageSize})` — the
  values never enter the POST body (user-observed: "the parameter didn't
  appear in the post data"). The backend cannot apply them anyway:
  gateway `ListDict` parses only pagination, proto is
  `rpc ListDict(basis.common.PageReq)`, repository `List` has no WHERE
  filters.
- **file**: the page sends `name` (plus `orderKey`/`desc`), but gateway
  `ListFile` decodes `common_pb.PageReq` (dropping `name`) and repository
  `List` has no filters.

## What Changes

- **Frontend** (`web/src/pages/sysManagement/dict/index.vue`): pass
  `cn_name`/`en_name` from `searchFormData` to `dictListApi`; also re-fetch
  after 重置 (reset) so the table reflects the cleared filters.
- **Proto** (`rpc/basis/proto/sysManagement/dict.proto`): add
  `ListDictReq { uint32 page; uint32 page_size; string cn_name; string
  en_name; }`; change `rpc ListDict(ListDictReq)`. Regenerate + `go mod
  tidy`.
- **Proto** (`rpc/basis/proto/sysTool/file.proto`): add `ListFileReq
  { uint32 page; uint32 page_size; string name; }`; change
  `rpc ListFile(ListFileReq)`. Regenerate + `go mod tidy`.
- **Gateway**: `ListDict` / `ListFile` parse the filter keys and forward
  the new req messages.
- **RPC logic/service/repo** for dict and file: thread filters; repository
  applies `cn_name ILIKE`, `en_name ILIKE` (dict) and `file_name ILIKE`
  (file) to BOTH count and data queries, using explicit `$N` placeholders
  (no `?`+`Rebind` mixing — the placeholder bug fixed in
  `fix-api-list-filters`).

## Impact

- Affected: two protos (regenerated), two gateway handlers, rpc
  logic/service/repo for dict and file, one frontend page. No DB schema
  change. `file.ts` API module already sends the right keys — no frontend
  change needed for file.
- Out of scope (unchanged from fix-api-list-filters): `orderKey`/`desc`
  table sorting (dead UI on file page) — tracked separately.
- Verification: replay exact frontend payloads per filter against
  `/dict/list` and `/file/list`, assert rows/total equal the subset of the
  full unfiltered dataset (psql blocked by credential policy — same
  full-dataset subset method); empty result returns `data:{}` (pre-existing
  proto3 quirk). Frontend dict page: filters fire on 查询 and reset restores
  the full list.
