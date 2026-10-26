-- 010_drop_dead_user_menus_endpoint.sql
-- Retires the dead POST /menu/user-menus gateway endpoint (OBS-005): the
-- frontend never calls it (the sidebar is served by GET /menu/list) and the
-- handler returned an empty body for everyone because it never passed
-- RoleIds. The route and its HTTP handler are removed in code; this migration
-- retires the matching permission row and its role grants so the policy
-- table no longer carries a policy for a non-existent route.
-- Idempotent: re-running is a no-op.
UPDATE sys_management_permission
SET deleted_at = NOW(), updated_at = NOW()
WHERE deleted_at IS NULL
  AND domain = 'api'
  AND resource = '/menu/user-menus';

DELETE FROM sys_management_role_permissions
WHERE permission_id IN (
  SELECT id FROM sys_management_permission WHERE resource = '/menu/user-menus'
);
