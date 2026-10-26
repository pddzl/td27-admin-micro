-- 011_unify_operation_log_uri.sql
-- Retires the legacy /opl/* operation-log aliases (OBS: unify-operation-log-uri).
-- The gateway now registers only the canonical /operation-log/list|delete|
-- delete-by-ids|cleanup routes; the frontend calls the canonical paths. Both
-- row sets existed live, so the /opl/* rows are retired (not renamed —
-- renaming would duplicate the canonical rows):
--   - soft-delete the api-domain permission rows and their role grants
--   - soft-delete the sys_management_api registry rows (hidden from the
--     API management page)
-- Action values on the canonical rows are already correct per the suffix
-- mapping (/list → read, /delete|/delete-by-ids → delete, /cleanup → execute).
-- Idempotent: re-running is a no-op. After applying, reload the policy
-- (POST /permission/reload-policy) or restart the gateway.

BEGIN;

-- 1) Retire api-domain permissions for the alias paths
UPDATE sys_management_permission
SET deleted_at = NOW(), updated_at = NOW()
WHERE deleted_at IS NULL
  AND domain = 'api'
  AND resource IN ('/opl/list', '/opl/delete', '/opl/delete-by-ids');

-- 2) Drop role grants pointing at the alias permission rows
DELETE FROM sys_management_role_permissions
WHERE permission_id IN (
  SELECT id FROM sys_management_permission
  WHERE domain = 'api'
    AND resource IN ('/opl/list', '/opl/delete', '/opl/delete-by-ids')
);

-- 3) Retire the alias rows in the API registry
UPDATE sys_management_api
SET deleted_at = NOW(), updated_at = NOW()
WHERE deleted_at IS NULL
  AND path IN ('/opl/list', '/opl/delete', '/opl/delete-by-ids');

COMMIT;
