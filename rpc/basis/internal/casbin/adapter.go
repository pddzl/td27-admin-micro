package casbin

import (
	"errors"
	"fmt"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/jmoiron/sqlx"
)

// errWriteUnsupported fails fast on direct enforcer policy mutations: policies
// are owned by the unified permission table. In-memory updates (e.g.
// RebuildRolePolicies) must disable auto-save around the mutation.
var errWriteUnsupported = errors.New("casbin adapter is read-only: manage policies via the unified permission table")

// PostgresAdapter 基于统一权限表的 Casbin 适配器
// 仅实现 LoadPolicy，通过统一权限表管理策略，不直接操作 Casbin 格式
type PostgresAdapter struct {
	db *sqlx.DB
}

func NewPostgresAdapter(db *sqlx.DB) *PostgresAdapter {
	return &PostgresAdapter{db: db}
}

// LoadPolicy 从统一权限表加载 API 域的策略到 Casbin model
func (a *PostgresAdapter) LoadPolicy(mod model.Model) error {
	type policyRow struct {
		Sub      string `db:"sub"`
		Resource string `db:"resource"`
		Action   string `db:"action"`
		Effect   string `db:"effect"`
	}

	var rows []policyRow
	err := a.db.Select(&rows, `
		SELECT CAST(rp.role_id AS TEXT) AS sub, p.resource, p.action, p.effect
		FROM sys_management_role_permissions rp
		JOIN sys_management_permission p ON rp.permission_id = p.id
		WHERE p.domain = 'api'
	`)
	if err != nil {
		return fmt.Errorf("load policy failed: %w", err)
	}

	for _, r := range rows {
		effect := r.Effect
		if effect == "" {
			effect = "allow"
		}
		line := fmt.Sprintf("p, %s, %s, %s, %s", r.Sub, r.Resource, r.Action, effect)
		persist.LoadPolicyLine(line, mod)
	}

	return nil
}

func (a *PostgresAdapter) SavePolicy(mod model.Model) error {
	return errWriteUnsupported
}

func (a *PostgresAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errWriteUnsupported
}

func (a *PostgresAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return errWriteUnsupported
}

func (a *PostgresAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errWriteUnsupported
}

func (a *PostgresAdapter) LoadFilteredPolicy(mod model.Model, filter interface{}) error {
	return a.LoadPolicy(mod)
}

func (a *PostgresAdapter) IsFiltered() bool {
	return false
}

func (a *PostgresAdapter) AddPolicies(sec string, ptype string, rules [][]string) error {
	return errWriteUnsupported
}

func (a *PostgresAdapter) RemovePolicies(sec string, ptype string, rules [][]string) error {
	return errWriteUnsupported
}

func (a *PostgresAdapter) UpdatePolicy(sec string, ptype string, oldRule, newRule []string) error {
	return errWriteUnsupported
}

func (a *PostgresAdapter) UpdatePolicies(sec string, ptype string, oldRules, newRules [][]string) error {
	return errWriteUnsupported
}

func (a *PostgresAdapter) UpdateFilteredPolicies(sec string, ptype string, newRules [][]string, fieldIndex int, fieldValues ...string) error {
	return errWriteUnsupported
}
