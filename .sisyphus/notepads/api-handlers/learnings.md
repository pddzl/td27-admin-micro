# Learnings - API Handler Implementation

## go-zero v1.9.4 Path Parameter Handling
- `rest.PathParam()` does NOT exist in go-zero v1.9.4
- Use `pathvar.Vars(r)["paramName"]` from `github.com/zeromicro/go-zero/rest/pathvar` to extract path params
- Route pattern syntax uses `:paramName` (e.g., `/api/dept/:id`)

## Handler Pattern
- Each service gets a struct with `svcCtx`, a constructor, and methods following `(w http.ResponseWriter, r *http.Request)` signature
- Body parsing: `json.NewDecoder(r.Body).Decode(&req)`
- Path params: `pathvar.Vars(r)["key"]`
- Query params: `r.URL.Query().Get("key")`
- Response: set headers, write status, `json.NewEncoder(w).Encode(pkg.Success(resp))`
- JWT middleware wraps private route handlers: `jwtMiddleware.Handle(handler.Method)`

## Proto-generated Types
- Proto client imports use the full module path: `td27/rpc/basis/types/authority/<service>_pb`
- Common types (IdReq, Empty, PageReq, SuccessResp) are in `td27/rpc/basis/types/common_pb`
- Optional proto fields are pointers in Go (`*uint64`, `*string`, `*uint32`, `*bool`)
- Non-optional proto fields are value types
