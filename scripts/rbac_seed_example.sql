-- =============================================================================
-- RBAC Seed Example (English)
-- =============================================================================
-- Purpose : Populate api-domain permissions so gateway per-request RBAC
--           enforcement can be enabled safely.
-- Enable  : set `EnforceRbac: true` in api/gateway/etc/gateway.yaml, then
--           restart the gateway. Policies are loaded by rpc/basis on startup
--           and refreshed periodically (Casbin.AutoLoadInterval) — or reload
--           on demand via POST /permission/reload-policy.
--
-- Action mapping (api/gateway/internal/middleware/authz.go):
--   GET                          -> read
--   POST .../list|/flat|/el-tree|/by-page|/get|/check|/validate|
--        /batch-check|/user-buttons|/user-menus|/get-el-tree-depts -> read
--   POST .../create              -> create
--   POST .../delete|/delete-by-ids -> delete
--   POST .../execute|/cleanup|/reload-policy -> execute
--   any other POST               -> update
--
-- Path patterns are matched with Casbin keyMatch2:
--   /user/*       matches /user/list, /user/create, ...
--   /api/dept/:id matches /api/dept/123
--   /*            matches every path
--   action 'all'  grants every action on the resource
-- =============================================================================

-- -----------------------------------------------------------------------------
-- 1. Example roles (adjust ids to your data)
-- -----------------------------------------------------------------------------
INSERT INTO sys_management_role (id, created_at, updated_at, role_name, parent_id)
VALUES
    (1, NOW(), NOW(), 'Super Admin', 0),
    (2, NOW(), NOW(), 'Auditor (Read-Only)', 0),
    (3, NOW(), NOW(), 'User Manager', 0)
ON CONFLICT (id) DO NOTHING;

-- -----------------------------------------------------------------------------
-- 2. API permissions
--    ids start at 9001 to avoid clashing with existing rows
--    effect defaults to 'allow'; a 'deny' rule overrides allow rules that
--    would otherwise grant the same request (deny-override)
-- -----------------------------------------------------------------------------
INSERT INTO sys_management_permission (id, created_at, updated_at, name, domain, resource, action, effect, domain_id)
VALUES
    -- Super Admin: everything
    (9001, NOW(), NOW(), 'Full access (wildcard)', 'api', '/*', 'all', 'allow', 0),

    -- Auditor: read-only browsing of every module
    (9101, NOW(), NOW(), 'Read all APIs', 'api', '/*', 'read', 'allow', 0),

    -- User Manager: read everything, manage users/roles/depts
    (9201, NOW(), NOW(), 'Read all APIs', 'api', '/*', 'read', 'allow', 0),
    (9202, NOW(), NOW(), 'User management',   'api', '/user/*', 'all', 'allow', 0),
    (9203, NOW(), NOW(), 'Role management',   'api', '/role/*', 'all', 'allow', 0),
    (9204, NOW(), NOW(), 'Dept management',   'api', '/dept/*', 'all', 'allow', 0),

    -- Deny example: User Manager may not delete users (overrides 9202)
    (9205, NOW(), NOW(), 'Deny user deletion', 'api', '/user/delete', 'delete', 'deny', 0)
ON CONFLICT (id) DO NOTHING;

-- -----------------------------------------------------------------------------
-- 3. Role -> permission bindings
-- -----------------------------------------------------------------------------
INSERT INTO sys_management_role_permissions (role_id, permission_id)
VALUES
    (1, 9001),                      -- Super Admin -> wildcard
    (2, 9101),                      -- Auditor     -> read everywhere
    (3, 9201), (3, 9202), (3, 9203), (3, 9204), (3, 9205)  -- User Manager
ON CONFLICT DO NOTHING;

-- -----------------------------------------------------------------------------
-- 4. Verify
-- -----------------------------------------------------------------------------
-- Check policies the way the Casbin adapter loads them:
--   SELECT CAST(rp.role_id AS TEXT) AS sub, p.resource, p.action, p.effect
--   FROM sys_management_role_permissions rp
--   JOIN sys_management_permission p ON rp.permission_id = p.id
--   WHERE p.domain = 'api';
--
-- Then:
--   1. Set EnforceRbac: true in api/gateway/etc/gateway.yaml
--   2. Restart the gateway (rpc/basis picks policies up on restart or via
--      Casbin.AutoLoadInterval)
--   3. Log in again — role ids are read from the JWT claims
--
-- Notes:
--   - Every role should be able to read /user/get-user-info and
--     /menu/user-menus, otherwise the frontend cannot bootstrap (both are
--     covered by the 'read /*' grants above).
--   - Role assignments are resolved live at check time; permission/policy
--     changes take effect without re-login (roles embedded in JWT are only
--     a fallback for legacy callers).
--   - Deny rules are subject-scoped and win over allow rules for the same
--     role (see permission 9205). Audit your policies with
--     POST /permission/lint — it reports dead allows, redundant rules and
--     unassigned api permissions.
--   - JIT elevation routes: /role/elevation/request and /role/elevation/my
--     are self-service (any authenticated user, no policy needed).
--     /role/elevation/decide, /revoke and /list are privileged — they fall
--     under the '/role/*' 'all' grants above, so read-only roles cannot
--     approve elevations. Who approves is also restricted server-side: a user
--     can never approve their own request.
--   - Anonymous endpoints (/login, /captcha, /health, /logout) are never
--     checked by the middleware.
-- -----------------------------------------------------------------------------
