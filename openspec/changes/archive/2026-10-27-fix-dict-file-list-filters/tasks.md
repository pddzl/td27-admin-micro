# Tasks

## 1. Proto

- [x] 1.1 `dict.proto`: add `ListDictReq` (page, page_size, cn_name,
      en_name); change `rpc ListDict` to accept it.
- [x] 1.2 `file.proto`: add `ListFileReq` (page, page_size, name); change
      `rpc ListFile` to accept it.
- [x] 1.3 Regenerate both; `go mod tidy`.

## 2. Backend

- [x] 2.1 Gateway `ListDict` / `ListFile`: parse filter keys; forward new
      req messages.
- [x] 2.2 RPC logic/service for dict and file: thread filter args.
- [x] 2.3 Repositories: WHERE `cn_name ILIKE`, `en_name ILIKE` (dict);
      `file_name ILIKE` (file) — count and data queries share the WHERE;
      explicit `$N` placeholders; `deleted_at IS NULL` preserved.

## 3. Frontend

- [x] 3.1 Dict page: pass `searchFormData.cn_name`/`en_name` into
      `dictListApi`; re-fetch after 重置.

## 4. Verification

- [x] 4.1 `go vet` + build both binaries clean.
- [x] 4.2 Live replay: dict cn_name~API (1), cn_name~HTTP (1),
      en_name~Methods (1), en_name~API (1), combo (0), nonexistent (0);
      file name~alpha (1), name~e2e-filter (2), name~.png (2), name~none
      (0), unfiltered (2) — every response's rows and total equal the
      subset of the full unfiltered dataset. (file table was empty, so two
      scratch files were uploaded through the API, used for the filter
      checks, and deleted afterwards — cleanup verified.)
- [x] 4.3 Browser check: dict page 查询 with 中文名称=API shows only
      API分组 (total 1); 重置 restores both rows (total 2); network shows
      `POST /dict/list` body containing `cn_name`.
- [x] 4.4 Archive this change after a green run.
