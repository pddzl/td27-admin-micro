-- 009_align_file_delete_action.sql
-- The gateway's action mapping classifies paths ending in /delete as action
-- 'delete' (suffix-based heuristics), but the /file/delete permission row was
-- seeded with action 'read'. The policy therefore never matched the enforced
-- action and every POST /file/delete was denied (403) even for the root role.
-- Idempotent: only updates the row when it still carries the wrong action.
UPDATE sys_management_permission
SET action = 'delete', updated_at = NOW()
WHERE deleted_at IS NULL
  AND domain = 'api'
  AND resource = '/file/delete'
  AND action <> 'delete';
