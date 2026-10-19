-- 008_grant_user_info_read.sql
-- Third "every authenticated user" endpoint (BUG-009 class, after 006/007):
-- GET /user/get-user-info (read) is called by the frontend router guard on
-- every boot (userStore.getInfo). Without the grant the guard's getInfo
-- rejects, the catch handler resets the token and bounces the user to login —
-- non-admin users could never stay logged in.
-- Grants the existing read allow permission to every non-deleted role.
-- Idempotent: ON CONFLICT DO NOTHING. Roles created AFTER this migration need
-- the same grant at assignment time (caveat as in 006/007).
--
-- Fixes BUG-009 (boot endpoints). See openspec/changes/rbac-regression-testing/bugs.md.

BEGIN;

INSERT INTO sys_management_role_permissions (role_id, permission_id, data_scope)
SELECT r.id,
       p.id,
       'all'
FROM sys_management_role r
CROSS JOIN sys_management_permission p
WHERE p.domain = 'api'
  AND p.resource = '/user/get-user-info'
  AND p.action = 'read'
  AND p.effect = 'allow'
  AND p.deleted_at IS NULL
  AND r.deleted_at IS NULL
ON CONFLICT (role_id, permission_id) DO NOTHING;

COMMIT;
