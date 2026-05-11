## Context

Four files currently combine dict and dict_detail logic. Model, repository, and service layers are already separated.

## Approach

### 1. Proto split
- `dict.proto`: Keep only Dict-related messages + `service Dict`. DictDetailResp stays here since DictResp references it.
- `dict_detail.proto`: New file with `service DictDetail` containing detail-specific RPCs (CreateDictDetail, UpdateDictDetail, DeleteDictDetail, FlatDictDetails, ListDictDetail).
- `DictDetailResp` stays in `dict.proto` (used by `DictResp.details`). DictDetail service proto imports from dict.proto.

### 2. Logic split
- `dict.go`: Keep DictLogic (GetDict, GetDictByENName, ListDict, CreateDict, UpdateDict, DeleteDict)
- `dict_detail.go`: New DictDetailLogic (CreateDictDetail, UpdateDictDetail, DeleteDictDetail, FlatDictDetails, ListDictDetail)

### 3. Server split
- `dict.go`: Keep DictServer
- `dict_detail.go`: New DictDetailServer

### 4. Handler split
- `dict.go`: Keep DictHandler with dict-only routes
- `dict_detail.go`: New DictDetailHandler with detail-only routes

### 5. Registration
- `basis.go`: Add `dict_detail_pb.RegisterDictDetailServer`
- Gateway `ServiceContext`: Add `DictDetailClient`
- `handler.go`: Add dict detail routes under existing dict section
