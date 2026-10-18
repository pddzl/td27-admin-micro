-- 005_perm_resource_effect_unique.sql
-- Fixes BUG-006: UNIQUE (resource) on sys_management_permission prevented an
-- allow AND a deny row for the same resource, making deny-wins authorization
-- impossible (Casbin 4-column model with effect requires both to coexist).
-- Replaced by UNIQUE (resource, effect) among live rows.
-- Idempotent: safe to re-apply.
-- See openspec/changes/rbac-regression-testing/bugs.md.

BEGIN;

ALTER TABLE sys_management_permission
  DROP CONSTRAINT IF EXISTS uni_sys_management_permission_resource;

CREATE UNIQUE INDEX IF NOT EXISTS uni_sys_management_permission_resource_effect
  ON sys_management_permission (resource, effect)
  WHERE deleted_at IS NULL;

COMMIT;
