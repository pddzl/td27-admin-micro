package sysManagement

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
)

// RoleRepository defines interface for role data operations
type RoleRepository interface {
	FindOne(ctx context.Context, id uint) (*sysManagement.RoleModel, error)
	FindOneWithChildren(ctx context.Context, id uint) (*sysManagement.RoleModel, error)
	FindAll(ctx context.Context) ([]*sysManagement.RoleModel, error)
	List(ctx context.Context, page *common.PageInfo) ([]*sysManagement.RoleModel, int64, error)
	Create(ctx context.Context, role *sysManagement.RoleModel) error
	Update(ctx context.Context, role *sysManagement.RoleModel) error
	Delete(ctx context.Context, id uint) error
	AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
	GetPermissions(ctx context.Context, roleID uint) ([]uint, error)
	GetUserRoles(ctx context.Context, userID uint) ([]*sysManagement.RoleModel, error)
	GetRolesByUserIDs(ctx context.Context, userIDs []uint) (map[uint][]*sysManagement.RoleModel, error)
	AssignUserRoles(ctx context.Context, userID uint, roleIDs []uint) error
}

type roleRepository struct {
	db *sqlx.DB
}

// NewRoleRepository creates a new role repository instance
func NewRoleRepository(db *sqlx.DB) RoleRepository {
	return &roleRepository{db: db}
}

const roleColumns = `id, created_at, updated_at, deleted_at, role_name, parent_id, COALESCE(permission_hash, '') as permission_hash`

func (r *roleRepository) FindOne(ctx context.Context, id uint) (*sysManagement.RoleModel, error) {
	var role sysManagement.RoleModel
	err := r.db.GetContext(ctx, &role, "SELECT "+roleColumns+" FROM sys_management_role WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindOneWithChildren(ctx context.Context, id uint) (*sysManagement.RoleModel, error) {
	var role sysManagement.RoleModel
	err := r.db.GetContext(ctx, &role, "SELECT "+roleColumns+" FROM sys_management_role WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Load children
	var children []*sysManagement.RoleModel
	err = r.db.SelectContext(ctx, &children,
		"SELECT "+roleColumns+" FROM sys_management_role WHERE parent_id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}
	role.Children = children
	return &role, nil
}

func (r *roleRepository) FindAll(ctx context.Context) ([]*sysManagement.RoleModel, error) {
	var roles []*sysManagement.RoleModel
	err := r.db.SelectContext(ctx, &roles, "SELECT "+roleColumns+" FROM sys_management_role WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) List(ctx context.Context, page *common.PageInfo) ([]*sysManagement.RoleModel, int64, error) {
	var total int64
	if err := r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM sys_management_role WHERE deleted_at IS NULL"); err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*sysManagement.RoleModel{}, 0, nil
	}

	offset := (page.Page - 1) * page.PageSize
	var roles []*sysManagement.RoleModel
	err := r.db.SelectContext(ctx, &roles,
		"SELECT "+roleColumns+" FROM sys_management_role WHERE deleted_at IS NULL ORDER BY id ASC LIMIT $1 OFFSET $2",
		page.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

func (r *roleRepository) Create(ctx context.Context, role *sysManagement.RoleModel) error {
	query := `INSERT INTO sys_management_role (created_at, updated_at, role_name, parent_id, permission_hash)
	           VALUES (NOW(), NOW(), :role_name, :parent_id, :permission_hash)
	           RETURNING id`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.GetContext(ctx, &role.ID, role)
}

func (r *roleRepository) Update(ctx context.Context, role *sysManagement.RoleModel) error {
	query := `UPDATE sys_management_role SET
	           updated_at = NOW(),
	           role_name = :role_name,
	           parent_id = :parent_id,
	           permission_hash = :permission_hash
	         WHERE id = :id AND deleted_at IS NULL`
	_, err := r.db.NamedExecContext(ctx, query, role)
	return err
}

func (r *roleRepository) Delete(ctx context.Context, id uint) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM sys_management_role_permissions WHERE role_id = $1", id); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM sys_management_user_roles WHERE role_id = $1", id); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM sys_management_role WHERE id = $1", id); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *roleRepository) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM sys_management_role_permissions WHERE role_id = $1", roleID); err != nil {
		return err
	}

	if len(permissionIDs) > 0 {
		query, args := buildInsertPairs("sys_management_role_permissions", "role_id", "permission_id", roleID, permissionIDs)
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// buildInsertPairs constructs a single multi-row INSERT statement for
// (parentID, childID) pairs, e.g.
// INSERT INTO t (parent_col, child_col) VALUES ($1, $2), ($1, $3)...
func buildInsertPairs(table, parentCol, childCol string, parentID uint, childIDs []uint) (string, []interface{}) {
	values := make([]string, 0, len(childIDs))
	args := make([]interface{}, 0, len(childIDs)+1)
	args = append(args, parentID)
	for i, id := range childIDs {
		values = append(values, fmt.Sprintf("($1, $%d)", i+2))
		args = append(args, id)
	}
	query := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES %s", table, parentCol, childCol, strings.Join(values, ", "))
	return query, args
}

func (r *roleRepository) GetPermissions(ctx context.Context, roleID uint) ([]uint, error) {
	var permissionIDs []uint
	err := r.db.SelectContext(ctx, &permissionIDs,
		"SELECT permission_id FROM sys_management_role_permissions WHERE role_id = $1", roleID)
	if err != nil {
		return nil, err
	}
	return permissionIDs, nil
}

func (r *roleRepository) GetUserRoles(ctx context.Context, userID uint) ([]*sysManagement.RoleModel, error) {
	var roles []*sysManagement.RoleModel
	err := r.db.SelectContext(ctx, &roles,
		`SELECT r.id, r.created_at as created_at, r.updated_at as updated_at, r.deleted_at, r.role_name, r.parent_id, COALESCE(r.permission_hash, '') as permission_hash
		 FROM sys_management_role r
		 JOIN sys_management_user_roles ur ON ur.role_id = r.id
		 WHERE ur.user_id = $1 AND r.deleted_at IS NULL`, userID)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) GetRolesByUserIDs(ctx context.Context, userIDs []uint) (map[uint][]*sysManagement.RoleModel, error) {
	result := make(map[uint][]*sysManagement.RoleModel, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
	}

	placeholders := make([]string, len(userIDs))
	args := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	var rows []struct {
		UserID uint `db:"user_id"`
		sysManagement.RoleModel
	}
	err := r.db.SelectContext(ctx, &rows,
		`SELECT ur.user_id, r.id, r.created_at as created_at, r.updated_at as updated_at, r.deleted_at, r.role_name, r.parent_id, COALESCE(r.permission_hash, '') as permission_hash
		 FROM sys_management_role r
		 JOIN sys_management_user_roles ur ON ur.role_id = r.id
		 WHERE ur.user_id IN (`+strings.Join(placeholders, ", ")+`) AND r.deleted_at IS NULL`, args...)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		role := row.RoleModel
		result[row.UserID] = append(result[row.UserID], &role)
	}
	return result, nil
}

func (r *roleRepository) AssignUserRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM sys_management_user_roles WHERE user_id = $1", userID); err != nil {
		return err
	}

	if len(roleIDs) > 0 {
		query, args := buildInsertPairs("sys_management_user_roles", "user_id", "role_id", userID, roleIDs)
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return err
		}
	}

	return tx.Commit()
}
