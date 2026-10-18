# Tasks

## 1. Menu module (reported bug)

- [x] 1.1 `web/src/pages/sysManagement/menu/index.vue` — `operateAction`: rename payload keys `keepAlive` → `keep_alive`, `alwaysShow` → `always_show`
- [x] 1.2 `web/src/pages/sysManagement/menu/index.vue` — `editMenuDialog`: prefill from `row.keep_alive` / `row.always_show`
- [x] 1.3 `web/src/pages/sysManagement/menu/index.vue` — table columns: `prop="parentId"` → `prop="parent_id"`, `prop="name"` → `prop="menu_name"`

## 2. User modify-password

- [x] 2.1 `web/src/api/sysManagement/user.ts` — `modifyPasswdApi` URL `/user/modify-passwd` → `/user/modify-password`; rename interface keys to `old_password` / `new_password`
- [x] 2.2 `web/src/pages/sysManagement/user/index.vue` — `mpOperateAction` payload keys `oldPassword` / `newPassword` → `old_password` / `new_password`

## 3. Page-size mapping on list endpoints

- [x] 3.1 `web/src/api/sysManagement/user.ts` — `userListApi` maps `pageSize` → `page_size`
- [x] 3.2 `web/src/api/sysManagement/role.ts` — `roleListApi` maps `pageSize` → `page_size`
- [x] 3.3 `web/src/api/sysManagement/dict.ts` — `dictListApi` maps `pageSize` → `page_size`
- [x] 3.4 `web/src/api/sysManagement/api.ts` — `apiListApi` maps `pageSize` → `page_size`
- [x] 3.5 `web/src/api/sysTool/cron.ts` — `cronListApi` maps `pageSize` → `page_size`
- [x] 3.6 `web/src/api/sysMonitor/operationLog.ts` — `oplListApi` maps `pageSize` → `page_size`
- [x] 3.7 `web/src/api/sysTool/service_token.ts` — `listServiceTokenApi` maps `pageSize` → `page_size`
- [x] 3.8 `web/src/api/sysManagement/button.ts` — `listButtonApi` maps `pageSize` → `page_size`

## 4. Route/method alignment (role buttons drawer)

- [x] 4.1 `web/src/api/sysManagement/button.ts` — `listButtonApi`: `POST /button/list` → `GET /button/list` with `page` / `page_size` query params
- [x] 4.2 `web/src/api/sysManagement/permission.ts` — `getRolePermissionsApi`: `POST` body → `GET /permission/by-role?id=`

## 5. Verification

- [x] 5.1 Type-check passes (`pnpm vue-tsc` / build)
- [x] 5.2 Runtime: menu update with 是 persists in DB
- [x] 5.3 Runtime: menu edit dialog prefills 一直显示根路由 correctly
- [x] 5.4 Runtime: user modify-password no longer 404s (wrong old password rejected with backend message, not 404)
- [x] 5.5 Runtime: `/user/list`, `/role/list`, `/opl/list` honor `page_size`
