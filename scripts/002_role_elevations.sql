-- Migration 002: role elevations for JIT (just-in-time) access.
-- A user requests temporary membership in a role; an approver grants it with
-- an expiry. Active elevations are merged into live role resolution, so
-- expired grants stop authorizing immediately (no token changes needed).

CREATE TABLE IF NOT EXISTS sys_management_role_elevations (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'revoked')),
    reason TEXT NOT NULL DEFAULT '',
    duration_minutes INT NOT NULL DEFAULT 60,
    expires_at TIMESTAMPTZ,
    requested_by BIGINT NOT NULL DEFAULT 0,
    decided_by BIGINT
);

CREATE INDEX IF NOT EXISTS idx_role_elevations_user_id ON sys_management_role_elevations (user_id);
CREATE INDEX IF NOT EXISTS idx_role_elevations_status ON sys_management_role_elevations (status);
