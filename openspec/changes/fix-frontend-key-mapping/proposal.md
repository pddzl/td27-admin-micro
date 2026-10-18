# Fix Frontend Request Key Mapping Mismatches

## Why

The gateway decodes HTTP JSON into protobuf-generated Go structs with standard
`encoding/json` (not protojson), so request fields only decode when the frontend
key exactly matches the proto json tag (snake_case). Several frontend payloads
send camelCase keys that are silently dropped, while the API still returns
success. Verified live against the running stack:

- `POST /menu/update` with the exact frontend payload returns `code:0` but
  `keep_alive`/`always_show` remain unchanged in the DB (reported bug: choosing
  是 and clicking 确认 in the menu dialog does nothing).
- `POST /user/modify-password` exists but the frontend calls
  `/user/modify-passwd` → 404; even with the right URL, `oldPassword` /
  `newPassword` keys do not decode (`old_password` / `new_password`).
- All paginated lists send `pageSize` while `common_pb.PageReq` expects
  `page_size`; every handler falls back to 10, so the UI page-size selector
  (20/50/…) silently has no effect.
- The role-management button-permission drawer is broken at the transport layer:
  `POST /button/list` (backend: GET with `page`/`page_size` query params) and
  `POST /permission/by-role` (backend: GET with `?id=`) both fail.

## What Changes

Frontend-only fixes (web/):

- **menu/index.vue**: send `keep_alive` / `always_show` in create/update
  payloads; read `row.keep_alive` / `row.always_show` when prefilling the edit
  dialog; fix 父节点 column prop (`parentId` → `parent_id`) and 路由名称 column
  prop (`name` → `menu_name`).
- **user.ts / user/index.vue**: rename the modify-password API URL to
  `/user/modify-password`; send `old_password` / `new_password`.
- **List request mapping** (user, role, dict, api, cron, operationLog, button,
  serviceToken): map `pageSize` → `page_size` following the existing convention
  in `dictDetail.ts` / `user.ts`.
- **button.ts**: `listButtonApi` switches to `GET /button/list` with
  `page` / `page_size` query params to match the registered route.
- **permission.ts**: `getRolePermissionsApi` switches to
  `GET /permission/by-role?id=` to match the registered route.

## Impact

- Only `web/src` files change; no backend, proto, or DB changes.
- Restores: menu boolean toggles (reported bug), user modify-password dialog,
  page-size selection on all list pages, role button-permission drawer.
- Not addressed here (separate follow-up needed, functional changes): service
  token module URL/key mismatches (`/service-token/*` vs `/serviceToken/*`,
  `apiIds` vs `permission_ids`, detail endpoint shape), cron create/update
  (`extraParams` key + `method` enum type mismatch), role inheritance dialog
  calling unregistered `/role/setInheritance`, dead role APIs
  (`updateRoleMenu`, `getRoleMenus`, `roleTree`, `getInheritance`).
