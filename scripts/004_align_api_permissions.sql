-- 004_align_api_permissions.sql
-- Aligns sys_management_api + api-domain sys_management_permission rows with the
-- routes actually registered in api/gateway/internal/handler/handler.go.
-- Fixes BUG-002 (path drift: camelCase/underscore vs current kebab routes) and
-- BUG-003 (wrong action values). Idempotent: safe to re-apply.
-- See openspec/changes/rbac-regression-testing/bugs.md.

BEGIN;

-- 1) Rename stale paths to the current registered routes (+ method fixes)
UPDATE sys_management_api SET path='/api/apis/delete-by-ids', method='POST', updated_at=NOW() WHERE path='/api/deleteByIds';
UPDATE sys_management_permission SET resource='/api/apis/delete-by-ids', updated_at=NOW() WHERE domain='api' AND resource='/api/deleteByIds';
UPDATE sys_management_api SET path='/api/el-tree', method='POST', updated_at=NOW() WHERE path='/api/elTree';
UPDATE sys_management_permission SET resource='/api/el-tree', updated_at=NOW() WHERE domain='api' AND resource='/api/elTree';
UPDATE sys_management_api SET path='/cron/delete-by-ids', method='POST', updated_at=NOW() WHERE path='/cron/deleteByIds';
UPDATE sys_management_permission SET resource='/cron/delete-by-ids', updated_at=NOW() WHERE domain='api' AND resource='/cron/deleteByIds';
UPDATE sys_management_api SET path='/cron/toggle-status', method='POST', updated_at=NOW() WHERE path='/cron/switchOpen';
UPDATE sys_management_permission SET resource='/cron/toggle-status', updated_at=NOW() WHERE domain='api' AND resource='/cron/switchOpen';
UPDATE sys_management_api SET path='/dept/get-el-tree-depts', method='POST', updated_at=NOW() WHERE path='/dept/getElTreeDepts';
UPDATE sys_management_permission SET resource='/dept/get-el-tree-depts', updated_at=NOW() WHERE domain='api' AND resource='/dept/getElTreeDepts';
UPDATE sys_management_api SET path='/dict-detail/create', method='POST', updated_at=NOW() WHERE path='/dict_detail/create';
UPDATE sys_management_permission SET resource='/dict-detail/create', updated_at=NOW() WHERE domain='api' AND resource='/dictDetail/create';
UPDATE sys_management_api SET path='/dict-detail/delete', method='POST', updated_at=NOW() WHERE path='/dict_detail/delete';
UPDATE sys_management_permission SET resource='/dict-detail/delete', updated_at=NOW() WHERE domain='api' AND resource='/dictDetail/delete';
UPDATE sys_management_api SET path='/dict-detail/list', method='POST', updated_at=NOW() WHERE path='/dict_detail/list';
UPDATE sys_management_permission SET resource='/dict-detail/list', updated_at=NOW() WHERE domain='api' AND resource='/dictDetail/list';
UPDATE sys_management_api SET path='/dict-detail/update', method='POST', updated_at=NOW() WHERE path='/dict_detail/update';
UPDATE sys_management_permission SET resource='/dict-detail/update', updated_at=NOW() WHERE domain='api' AND resource='/dictDetail/update';
UPDATE sys_management_api SET path='/file/download/:id', method='GET', updated_at=NOW() WHERE path='/file/download';
UPDATE sys_management_permission SET resource='/file/download/:id', updated_at=NOW() WHERE domain='api' AND resource='/file/download';
UPDATE sys_management_api SET path='/menu/el-tree', method='POST', updated_at=NOW() WHERE path='/menu/elTree';
UPDATE sys_management_permission SET resource='/menu/el-tree', updated_at=NOW() WHERE domain='api' AND resource='/menu/elTree';
UPDATE sys_management_api SET path='/opl/delete-by-ids', method='POST', updated_at=NOW() WHERE path='/opl/deleteByIds';
UPDATE sys_management_permission SET resource='/opl/delete-by-ids', updated_at=NOW() WHERE domain='api' AND resource='/opl/deleteByIds';
UPDATE sys_management_api SET path='/permission/reload-policy', method='POST', updated_at=NOW() WHERE path='/role_permission/rebuild';
UPDATE sys_management_permission SET resource='/permission/reload-policy', updated_at=NOW() WHERE domain='api' AND resource='/role_permission/rebuild';
UPDATE sys_management_api SET path='/service-token/:id', method='GET', updated_at=NOW() WHERE path='/serviceToken/detail';
UPDATE sys_management_permission SET resource='/service-token/:id', updated_at=NOW() WHERE domain='api' AND resource='/serviceToken/detail';
UPDATE sys_management_api SET path='/user/get-user-info', method='GET', updated_at=NOW() WHERE path='/user/getUserInfo';
UPDATE sys_management_permission SET resource='/user/get-user-info', updated_at=NOW() WHERE domain='api' AND resource='/user/getUserInfo';
UPDATE sys_management_api SET path='/user/modify-password', method='POST', updated_at=NOW() WHERE path='/user/modifyPasswd';
UPDATE sys_management_permission SET resource='/user/modify-password', updated_at=NOW() WHERE domain='api' AND resource='/user/modifyPasswd';
UPDATE sys_management_api SET path='/user/switch-active', method='POST', updated_at=NOW() WHERE path='/user/switchActive';
UPDATE sys_management_permission SET resource='/user/switch-active', updated_at=NOW() WHERE domain='api' AND resource='/user/switchActive';

-- 2) Recompute action for every api-domain permission (mirrors httpMethodToAction)
UPDATE sys_management_permission p
SET action = CASE
      WHEN a.method IN ('GET','HEAD','OPTIONS') THEN 'read'
      WHEN p.resource LIKE '%/list' OR p.resource LIKE '%/flat' OR p.resource LIKE '%/el-tree' OR p.resource LIKE '%/by-page' OR p.resource LIKE '%/get' OR p.resource LIKE '%/check' OR p.resource LIKE '%/validate' OR p.resource LIKE '%/batch-check' OR p.resource LIKE '%/user-buttons' OR p.resource LIKE '%/user-menus' OR p.resource LIKE '%/get-el-tree-depts' THEN 'read'
      WHEN p.resource LIKE '%/delete' OR p.resource LIKE '%/delete-by-ids' THEN 'delete'
      WHEN p.resource LIKE '%/create' THEN 'create'
      WHEN p.resource LIKE '%/execute' OR p.resource LIKE '%/cleanup' OR p.resource LIKE '%/reload-policy' THEN 'execute'
      ELSE 'update'
    END
FROM sys_management_api a
WHERE p.domain = 'api' AND p.resource = a.path AND p.deleted_at IS NULL;

-- 4) Ensure an allow-permission row exists for every api registry row (domain_id = api id)
INSERT INTO sys_management_permission (created_at, updated_at, name, resource, action, domain_id, domain, effect)
SELECT NOW(), NOW(), a.description, a.path, CASE
      WHEN a.method IN ('GET','HEAD','OPTIONS') THEN 'read'
      WHEN a.path LIKE '%/list' OR a.path LIKE '%/flat' OR a.path LIKE '%/el-tree' OR a.path LIKE '%/by-page' OR a.path LIKE '%/get' OR a.path LIKE '%/check' OR a.path LIKE '%/validate' OR a.path LIKE '%/batch-check' OR a.path LIKE '%/user-buttons' OR a.path LIKE '%/user-menus' OR a.path LIKE '%/get-el-tree-depts' THEN 'read'
      WHEN a.path LIKE '%/delete' OR a.path LIKE '%/delete-by-ids' THEN 'delete'
      WHEN a.path LIKE '%/create' THEN 'create'
      WHEN a.path LIKE '%/execute' OR a.path LIKE '%/cleanup' OR a.path LIKE '%/reload-policy' THEN 'execute'
      ELSE 'update'
    END, a.id, 'api', 'allow'
FROM sys_management_api a
WHERE a.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM sys_management_permission p
    WHERE p.domain = 'api' AND p.resource = a.path AND p.deleted_at IS NULL
  )
ON CONFLICT DO NOTHING;

-- 5) Grant every api-domain permission to the root role (role_id=1); other roles untouched
INSERT INTO sys_management_role_permissions (role_id, permission_id)
SELECT 1, p.id
FROM sys_management_permission p
WHERE p.domain = 'api' AND p.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM sys_management_role_permissions rp
    WHERE rp.role_id = 1 AND rp.permission_id = p.id
  );

COMMIT;
