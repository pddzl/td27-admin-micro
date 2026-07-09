package sysManagement

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"td27/rpc/basis/internal/model/sysManagement"
)

// PermissionRepository defines interface for permission data operations
type PermissionRepository interface {
	FindOne(ctx context.Context, id uint) (*sysManagement.PermissionModel, error)
	FindAll(ctx context.Context) ([]*sysManagement.PermissionModel, error)
	FindByDomain(ctx context.Context, domain sysManagement.PermissionDomain) ([]*sysManagement.PermissionModel, error)
	FindByRoleID(ctx context.Context, roleID uint) ([]*sysManagement.PermissionModel, error)
	FindByResourceAndAction(ctx context.Context, resource string, action sysManagement.Action) (*sysManagement.PermissionModel, error)
	Create(ctx context.Context, permission *sysManagement.PermissionModel) error
	Update(ctx context.Context, permission *sysManagement.PermissionModel) error
	Delete(ctx context.Context, id uint) error
}

type permissionRepository struct {
	db *sqlx.DB
}

// NewPermissionRepository creates a new permission repository instance
func NewPermissionRepository(db *sqlx.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

const permColumns = `id, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at, deleted_at, name, domain, resource, action, domain_id`

func (r *permissionRepository) FindOne(ctx context.Context, id uint) (*sysManagement.PermissionModel, error) {
	var perm sysManagement.PermissionModel
	err := r.db.GetContext(ctx, &perm, "SELECT "+permColumns+" FROM sys_management_permission WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &perm, nil
}

func (r *permissionRepository) FindAll(ctx context.Context) ([]*sysManagement.PermissionModel, error) {
	var perms []*sysManagement.PermissionModel
	err := r.db.SelectContext(ctx, &perms, "SELECT "+permColumns+" FROM sys_management_permission WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	return perms, nil
}

func (r *permissionRepository) FindByDomain(ctx context.Context, domain sysManagement.PermissionDomain) ([]*sysManagement.PermissionModel, error) {
	var perms []*sysManagement.PermissionModel
	err := r.db.SelectContext(ctx, &perms,
		"SELECT "+permColumns+" FROM sys_management_permission WHERE domain=$1 AND deleted_at IS NULL", domain)
	if err != nil {
		return nil, err
	}
	return perms, nil
}

func (r *permissionRepository) FindByRoleID(ctx context.Context, roleID uint) ([]*sysManagement.PermissionModel, error) {
	var perms []*sysManagement.PermissionModel
	err := r.db.SelectContext(ctx, &perms,
		`SELECT p.id, COALESCE(p.created_at, NOW()) as created_at, COALESCE(p.updated_at, NOW()) as updated_at, p.deleted_at, p.name, p.domain, p.resource, p.action, p.domain_id
		 FROM sys_management_permission p
		 JOIN sys_management_role_permissions rp ON rp.permission_id = p.id
		 WHERE rp.role_id = $1 AND p.deleted_at IS NULL`, roleID)
	if err != nil {
		return nil, err
	}
	return perms, nil
}

func (r *permissionRepository) FindByResourceAndAction(ctx context.Context, resource string, action sysManagement.Action) (*sysManagement.PermissionModel, error) {
	var perm sysManagement.PermissionModel
	err := r.db.GetContext(ctx, &perm,
		"SELECT "+permColumns+" FROM sys_management_permission WHERE resource=$1 AND action=$2 AND deleted_at IS NULL",
		resource, action)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &perm, nil
}

func (r *permissionRepository) Create(ctx context.Context, permission *sysManagement.PermissionModel) error {
	query := `INSERT INTO sys_management_permission (created_at, updated_at, name, domain, resource, action, domain_id)
	           VALUES (NOW(), NOW(), :name, :domain, :resource, :action, :domain_id)
	           RETURNING id`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.GetContext(ctx, &permission.ID, permission)
}

func (r *permissionRepository) Update(ctx context.Context, permission *sysManagement.PermissionModel) error {
	query := `UPDATE sys_management_permission SET
	           updated_at = NOW(),
	           name = :name,
	           domain = :domain,
	           resource = :resource,
	           action = :action,
	           domain_id = :domain_id
	         WHERE id = :id AND deleted_at IS NULL`
	_, err := r.db.NamedExecContext(ctx, query, permission)
	return err
}

func (r *permissionRepository) Delete(ctx context.Context, id uint) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM sys_management_role_permissions WHERE permission_id = $1", id); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "UPDATE sys_management_permission SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id); err != nil {
		return err
	}
	return tx.Commit()
}
