-- 007_grant_self_menu_button_read.sql
-- Two more "every authenticated user" endpoints, same class as BUG-007:
--   GET/POST /menu/user-menus  (read) — the sidebar's own menu tree
--   GET/POST /button/user-buttons (read) — the v-permission code preload;
--     when it 403s the frontend cache stays empty and every gated button
--     silently disappears (fail-closed UI).
-- Grants the existing read allow permissions for both routes to every
-- non-deleted role. Idempotent: ON CONFLICT DO NOTHING. Roles created AFTER
-- this migration need the same grants at assignment time (same caveat as
-- 006, see BUG-007/BUG-009 in openspec/changes/rbac-regression-testing/bugs.md).
--
-- Fixes BUG-009. See openspec/changes/rbac-regression-testing/bugs.md.

BEGIN;

INSERT INTO sys_management_role_permissions (role_id, permission_id, data_scope)
SELECT r.id,
       p.id,
       'all'
FROM sys_management_role r
CROSS JOIN sys_management_permission p
WHERE p.domain = 'api'
  AND p.resource IN ('/menu/user-menus', '/button/user-buttons')
  AND p.action = 'read'
  AND p.effect = 'allow'
  AND p.deleted_at IS NULL
  AND r.deleted_at IS NULL
ON CONFLICT (role_id, permission_id) DO NOTHING;

COMMIT;
