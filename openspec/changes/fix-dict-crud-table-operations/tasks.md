## 1. Fix dict list handler

- [x] 1.1 Change `/dict/list` handler from URL query params to POST body

## 2. Fix dict create API response

- [x] 2.1 Add `RETURNING id` to dict INSERT query
- [x] 2.2 Return created dict data (DictResp) instead of SuccessResp from CreateDict
- [x] 2.3 Frontend: use `push(res.data)` after create

## 3. Fix dict delete

- [x] 3.1 Frontend: use `splice()` instead of `getTableData()`

## 4. Verify

- [x] 4.1 Create dict → table updates instantly without getTableData
- [x] 4.2 Delete dict → table updates instantly without getTableData
- [x] 4.3 Pagination works correctly with POST body
