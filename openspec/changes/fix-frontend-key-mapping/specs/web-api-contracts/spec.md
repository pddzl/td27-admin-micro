## ADDED Requirements

### Requirement: Frontend request payloads MUST use gateway contract keys

The web frontend MUST send JSON request fields whose names exactly match the
protobuf json tags the gateway decodes with `encoding/json`, and MUST call the
HTTP method and path the gateway registers. A mismatch MUST NOT silently drop
fields while reporting success.

#### Scenario: Menu boolean flags persist

- **WHEN** an admin edits a menu in 系统管理 → 菜单管理, sets 一直显示根路由 or
  KeepAlive to 是, and clicks 确认
- **THEN** the update request contains `keep_alive` / `always_show`
- **AND** the database row reflects the chosen values after the success message

#### Scenario: Admin modify-password reaches the backend

- **WHEN** an admin submits the 修改密码 dialog in 用户管理
- **THEN** the request targets `POST /user/modify-password` with
  `old_password` / `new_password` keys
- **AND** the backend responds with its validation result instead of HTTP 404

#### Scenario: Page size selection is honored

- **WHEN** any paginated list page (users, roles, dicts, APIs, crons, operation
  logs, buttons, service tokens) sends a page size larger than 10
- **THEN** the request maps `pageSize` to `page_size`
- **AND** the backend returns that number of rows per page

#### Scenario: Role buttons drawer loads

- **WHEN** the button-permission drawer opens for a role in 角色管理
- **THEN** button list and role permission queries use the registered
  `GET /button/list` and `GET /permission/by-role?id=` routes
- **AND** no 404/405 transport errors occur
