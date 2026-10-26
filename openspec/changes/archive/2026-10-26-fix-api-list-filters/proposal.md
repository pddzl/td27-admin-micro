# Fix Dead Search Filters on POST /api/list

## Why

The API management page's search form sends `{path, group_en, method,
description, page, page_size}` (`web/src/api/sysManagement/api.ts`), but the
gateway `ListAPI` handler (`api/gateway/internal/handler/basis/sysManagement/
api.go`) parses only `page`/`page_size` and forwards `common_pb.PageReq`,
which has no filter fields. The rpc signature is
`rpc ListAPI(basis.common.PageReq)` and the repository filters only by an
optional `group_en` that is always passed as `nil`. Result: all four search
filters (path, group_en, method, description) are silently dropped — the
user observed picking `method: GET` changes nothing. This is the same
"silently dropped request key" family as the `/opl/list` and `page_size`
bugs.

Not caught by E2E: no journey has ever driven the API page's search form —
J12's live run was also aborted before it started. Filter assertions did not
exist in the catalog; only list invariants and pagination were covered.

## What Changes

- **Proto** (`rpc/basis/proto/sysManagement/api.proto`): add
  `message ListAPIReq { uint32 page; uint32 page_size; string path; string
  group_en; string method; string description; }` and change
  `rpc ListAPI(ListAPIReq) returns (ListAPIResp)`. Regenerate + `go mod
  tidy`.
- **Gateway** (`api.go` ListAPI): parse the four filter keys alongside
  page/page_size and pass `api_pb.ListAPIReq`.
- **RPC logic** (`rpc/basis/internal/logic/sysManagement/api.go`): accept
  `ListAPIReq`, build a filter struct, keep `common.NewPageInfo` defaults.
- **Service** (`rpc/basis/internal/service/sysManagement/api.go`): thread
  the filter struct through `List`.
- **Repository** (`rpc/basis/internal/repository/sysManagement/api.go`):
  extend `List` WHERE clauses — partial match (`ILIKE '%x%'`) for the text
  inputs `path`, `group_en`, `description`; exact match for the dropdown
  `method`. Count query and data query share the same WHERE; keep
  `deleted_at IS NULL`.
- **E2E catalog** (`skills/e2e-testing/references/flow-catalog.md`): add
  J3 step 8 — exercise every search filter and assert every returned row
  matches the filter (rendered rows AND response `total` consistent with
  the filtered SQL count). Spec delta on `e2e-testing`.

## Impact

- Affected: proto (regenerated), gateway handler, rpc logic/service/repo,
  E2E skill docs. No DB schema change, no frontend change (the wire
  contract already sends the right keys — the backend must catch up).
- `orderKey`/`desc` table-sort fields exist in the frontend search form but
  are never sent by `apiListApi` — dead UI, tracked as a separate follow-up,
  NOT in scope here.
- Verification: replay the exact frontend payload (`method: "GET"`) against
  `/api/list`, assert every row has `method == "GET"`, `total` matches the
  filtered SQL count, and the list invariants hold; repeat per filter.
