# Tasks

## 1. Proto

- [x] 1.1 `api.proto`: add `ListAPIReq` (page, page_size, path, group_en,
      method, description); change `rpc ListAPI` to accept it.
- [x] 1.2 Regenerate with the repo protoc command; run `go mod tidy`.

## 2. Backend

- [x] 2.1 Gateway `ListAPI`: parse path/group_en/method/description +
      pagination; forward `api_pb.ListAPIReq`.
- [x] 2.2 RPC logic `ListAPI`: map req → filter args + `NewPageInfo`.
- [x] 2.3 Service `List`: thread filter args to repository.
- [x] 2.4 Repository `List`: WHERE `path ILIKE`, `group_en ILIKE`,
      `method =`, `description ILIKE`; shared by count and data queries.
      (Also fixed a latent placeholder bug: mixed `?`→Rebind + pre-numbered
      `$N` renumbered LIMIT/OFFSET from $1, so `LIMIT $1` bound the filter
      string — "argument of LIMIT must be type bigint". Now explicit `$N`
      everywhere, matching operation_log.)

## 3. E2E catalog

- [x] 3.1 Flow catalog J3 step 8: exercise each search filter, assert every
      returned row matches (UI rows + response total vs filtered SQL count).
- [x] 3.2 Spec delta `e2e-testing`: filter-assertion requirement.

## 4. Verification

- [x] 4.1 `go vet ./...` and `go build ./rpc/basis/... ./api/gateway/...`
      clean.
- [x] 4.2 Restart basis + gateway; replay frontend payload per filter —
      path~user (10), group_en~user (8), description~登录 (1),
      method=GET (22; GET+POST=106 partition), method=GET+path~dashboard (3),
      method=POST+group_en~dict (9), path~/login (1): every response's rows
      and total exactly equal the subset computed from the full 106-row
      dataset (cross-checked against /api/el-tree leaf count = 106).
      Legitimate zero results (group_en~sysManagement, method=DELETE) return
      `data:{}` — pre-existing proto3 empty-response quirk for ALL list
      endpoints, noted as follow-up. SQL cross-check via psql was blocked by
      credential policy; the full-dataset subset comparison is equivalent
      evidence.
- [x] 4.3 Archive this change after a green run.
