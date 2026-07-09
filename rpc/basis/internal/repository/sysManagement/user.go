package sysManagement

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
)

// UserRepository defines interface for user data operations
type UserRepository interface {
	FindOne(ctx context.Context, id uint) (*sysManagement.UserModel, error)
	FindOneByUsername(ctx context.Context, username string) (*sysManagement.UserModel, error)
	FindOneWithRoles(ctx context.Context, id uint) (*sysManagement.UserModel, error)
	List(ctx context.Context, page *common.PageInfo, deptID *uint) ([]*sysManagement.UserModel, int64, error)
	Create(ctx context.Context, user *sysManagement.UserModel) error
	Update(ctx context.Context, user *sysManagement.UserModel) error
	Delete(ctx context.Context, id uint) error
	UpdatePassword(ctx context.Context, id uint, newPasswordHash string) error
	ToggleActive(ctx context.Context, id uint, active bool) error
	CountByRoleID(ctx context.Context, roleID uint) (int64, error)
}

type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

const userColumns = `id, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at, deleted_at, username, password, phone, email, active, dept_id`

func (r *userRepository) FindOne(ctx context.Context, id uint) (*sysManagement.UserModel, error) {
	var user sysManagement.UserModel
	err := r.db.GetContext(ctx, &user, "SELECT "+userColumns+" FROM sys_management_user WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindOneByUsername(ctx context.Context, username string) (*sysManagement.UserModel, error) {
	var user sysManagement.UserModel
	err := r.db.GetContext(ctx, &user, "SELECT "+userColumns+" FROM sys_management_user WHERE username=$1 AND deleted_at IS NULL", username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindOneWithRoles(ctx context.Context, id uint) (*sysManagement.UserModel, error) {
	var user sysManagement.UserModel
	err := r.db.GetContext(ctx, &user, "SELECT "+userColumns+" FROM sys_management_user WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Load associated roles
	var roles []*sysManagement.RoleModel
	err = r.db.SelectContext(ctx, &roles,
		`SELECT r.id, COALESCE(r.created_at, NOW()) as created_at, COALESCE(r.updated_at, NOW()) as updated_at, r.deleted_at, r.role_name, r.parent_id, COALESCE(r.permission_hash, '') as permission_hash
		 FROM sys_management_role r
		 JOIN sys_management_user_roles ur ON ur.role_id = r.id
		 WHERE ur.user_id = $1 AND r.deleted_at IS NULL`, id)
	if err != nil {
		return nil, err
	}
	user.Roles = roles
	return &user, nil
}

func (r *userRepository) List(ctx context.Context, page *common.PageInfo, deptID *uint) ([]*sysManagement.UserModel, int64, error) {
	var total int64
	baseWhere := "WHERE deleted_at IS NULL"
	args := []interface{}{}

	if deptID != nil {
		baseWhere += " AND dept_id = $1"
		args = append(args, *deptID)
	}

	countQuery := "SELECT COUNT(*) FROM sys_management_user " + baseWhere
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*sysManagement.UserModel{}, 0, nil
	}

	offset := (page.Page - 1) * page.PageSize
	selectQuery := "SELECT " + userColumns + " FROM sys_management_user " + baseWhere
	if deptID != nil {
		selectQuery += " ORDER BY id ASC LIMIT $2 OFFSET $3"
		args = append(args, page.PageSize, offset)
	} else {
		selectQuery += " ORDER BY id ASC LIMIT $1 OFFSET $2"
		args = append(args, page.PageSize, offset)
	}

	var users []*sysManagement.UserModel
	if err := r.db.SelectContext(ctx, &users, selectQuery, args...); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepository) Create(ctx context.Context, user *sysManagement.UserModel) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `INSERT INTO sys_management_user (created_at, updated_at, username, password, phone, email, active, dept_id)
	           VALUES (:created_at, :updated_at, :username, :password, :phone, :email, :active, :dept_id)
	           RETURNING id`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.GetContext(ctx, &user.ID, user)
}

func (r *userRepository) Update(ctx context.Context, user *sysManagement.UserModel) error {
	query := `UPDATE sys_management_user SET
	           updated_at = NOW(),
	           username = :username,
	           phone = :phone,
	           email = :email,
	           active = :active,
	           dept_id = :dept_id
	         WHERE id = :id AND deleted_at IS NULL`
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx, "UPDATE sys_management_user SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id)
	return err
}

func (r *userRepository) UpdatePassword(ctx context.Context, id uint, newPasswordHash string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sys_management_user SET password = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL",
		newPasswordHash, id)
	return err
}

func (r *userRepository) ToggleActive(ctx context.Context, id uint, active bool) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sys_management_user SET active = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL",
		active, id)
	return err
}

func (r *userRepository) CountByRoleID(ctx context.Context, roleID uint) (int64, error) {
	var count int64
	err := r.db.GetContext(ctx, &count,
		"SELECT COUNT(*) FROM sys_management_user_roles WHERE role_id = $1", roleID)
	return count, err
}
