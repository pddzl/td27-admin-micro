-- 006_grant_menu_list_read.sql
-- Every authenticated user needs GET /menu/list (read) to render their sidebar.
-- Without this grant the authz middleware rejects the call (403) and the user
-- sees no dynamic menus at all. Grants the existing '/menu/list' read allow
-- permission to every non-deleted role. Idempotent: ON CONFLICT DO NOTHING.
-- Note: roles created AFTER this migration still need this grant assigned at
-- role-creation/assignment time (documented caveat, see
-- openspec/changes/rbac-regression-testing/bugs.md BUG-007).
--
-- Fixes BUG-007. See openspec/changes/rbac-regression-testing/bugs.md.

BEGIN;

INSERT INTO sys_management_role_permissions (role_id, permission_id, data_scope)
SELECT r.id,
       p.id,
       'all'
FROM sys_management_role r
CROSS JOIN sys_management_permission p
WHERE p.domain = 'api'
  AND p.resource = '/menu/list'
  AND p.action = 'read'
  AND p.effect = 'allow'
  AND p.deleted_at IS NULL
  AND r.deleted_at IS NULL
ON CONFLICT (role_id, permission_id) DO NOTHING;

COMMIT;
