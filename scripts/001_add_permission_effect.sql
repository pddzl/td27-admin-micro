-- Migration 001: add effect column for deny-override RBAC support.
-- Existing rows default to 'allow' (previous behavior unchanged).

ALTER TABLE sys_management_permission
    ADD COLUMN IF NOT EXISTS effect TEXT NOT NULL DEFAULT 'allow';

ALTER TABLE sys_management_permission
    DROP CONSTRAINT IF EXISTS sys_management_permission_effect_check;

ALTER TABLE sys_management_permission
    ADD CONSTRAINT sys_management_permission_effect_check CHECK (effect IN ('allow', 'deny'));
