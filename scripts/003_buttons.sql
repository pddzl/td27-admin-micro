-- Migration 003: button permission storage + seed data
-- Creates the missing sys_management_button table, seeds button codes for
-- high-value actions, wires them into the unified permission table and grants
-- them to the "root" role (id = 1).
--
-- Apply before restarting the services:  psql -d td27_micro -f scripts/003_buttons.sql

-- ---------------------------------------------------------------- table ----
CREATE TABLE IF NOT EXISTS sys_management_button (
    id          bigint PRIMARY KEY,
    created_at  timestamp,
    updated_at  timestamp,
    deleted_at  timestamp,
    button_code  varchar(100) NOT NULL,
    button_name  varchar(100) NOT NULL,
    description  text         NOT NULL DEFAULT '',
    page_path    varchar(200) NOT NULL DEFAULT ''
);

CREATE SEQUENCE IF NOT EXISTS sys_management_button_id_seq
    AS bigint START WITH 1 INCREMENT BY 1;
ALTER TABLE sys_management_button
    ALTER COLUMN id SET DEFAULT nextval('sys_management_button_id_seq');

CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_management_button_code
    ON sys_management_button (button_code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_management_button_deleted_at
    ON sys_management_button (deleted_at);
CREATE INDEX IF NOT EXISTS idx_sys_management_button_page_path
    ON sys_management_button (page_path);

-- ----------------------------------------------------------------- seed ----
-- Codes follow the <page>:<action> convention; page_path matches the router
-- path so the assignment UI and v-permission stay consistent.
INSERT INTO sys_management_button (button_code, button_name, description, page_path)
VALUES
    ('user:create',          '创建用户',   'Create new user',                    '/user'),
    ('user:update',          '编辑用户',   'Update user profile',                '/user'),
    ('user:delete',          '删除用户',   'Soft delete user',                   '/user'),
    ('user:reset-password',  '重置密码',   'Reset / modify user password',       '/user'),
    ('role:create',          '创建角色',   'Create role',                        '/role'),
    ('role:update',          '编辑角色',   'Update role',                        '/role'),
    ('role:delete',          '删除角色',   'Soft delete role',                   '/role'),
    ('role:assign-permissions', '分配权限', 'Assign menu/api/button permissions', '/role'),
    ('role:assign-users',    '分配用户',   'Assign users to role',               '/role'),
    ('cron:create',          '创建任务',   'Create cron job',                    '/cron'),
    ('cron:update',          '编辑任务',   'Update cron job',                    '/cron'),
    ('cron:execute',         '执行任务',   'Run cron job immediately',           '/cron'),
    ('cron:delete',          '删除任务',   'Delete cron job',                    '/cron'),
    ('cron:toggle',          '启停任务',   'Enable / disable cron job',          '/cron'),
    ('cache:get',            '读取缓存',   'Read cache value',                   '/cache'),
    ('cache:delete',         '删除缓存',   'Delete cache key',                   '/cache'),
    ('cache:cleanup',        '清理缓存',   'Bulk cache cleanup',                 '/cache'),
    ('token:create',         '创建令牌',   'Create service token',               '/service-token'),
    ('token:update',         '编辑令牌',   'Update service token',               '/service-token'),
    ('token:delete',         '删除令牌',   'Soft delete service token',          '/service-token'),
    ('token:assign-permissions', '分配权限', 'Assign permissions to service token', '/service-token')
ON CONFLICT DO NOTHING;

-- ------------------------------------------------------- permission rows ----
-- One unified-permission row per button (domain='button', resource=code).
-- Name embeds the unique button_code to stay unique against the global name index.
INSERT INTO sys_management_permission (name, domain, resource, action, effect, domain_id)
SELECT 'button: ' || b.button_code, 'button', b.button_code, 'all', 'allow', b.id
FROM sys_management_button b
WHERE b.deleted_at IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_management_permission p
      WHERE p.domain = 'button' AND p.domain_id = b.id AND p.deleted_at IS NULL
  );

-- ------------------------------------------------- grant to root role (1) ----
INSERT INTO sys_management_role_permissions (role_id, permission_id)
SELECT 1, p.id
FROM sys_management_permission p
WHERE p.domain = 'button' AND p.deleted_at IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM sys_management_role_permissions rp
      WHERE rp.role_id = 1 AND rp.permission_id = p.id
  );
