package sysManagement

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"td27/rpc/basis/internal/model/sysManagement"
)

// MenuRepository defines interface for menu data operations
type MenuRepository interface {
	FindOne(ctx context.Context, id uint) (*sysManagement.MenuModel, error)
	FindAll(ctx context.Context) ([]*sysManagement.MenuModel, error)
	FindByParentID(ctx context.Context, parentID uint) ([]*sysManagement.MenuModel, error)
	FindByRoleIDs(ctx context.Context, roleIDs []uint) ([]*sysManagement.MenuModel, error)
	Create(ctx context.Context, menu *sysManagement.MenuModel) error
	Update(ctx context.Context, menu *sysManagement.MenuModel) error
	Delete(ctx context.Context, id uint) error
}

const menuColumns = `id, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at, deleted_at,
	COALESCE(menu_name, '') as menu_name, COALESCE(icon, '') as icon, COALESCE(path, '') as path,
	COALESCE(component, '') as component, COALESCE(redirect, '') as redirect, parent_id, sort, hidden, keep_alive, affix, always_show, COALESCE(title, '') as title`

type menuRepository struct {
	db *sqlx.DB
}

// NewMenuRepository creates a new menu repository instance
func NewMenuRepository(db *sqlx.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) FindOne(ctx context.Context, id uint) (*sysManagement.MenuModel, error) {
	var menu sysManagement.MenuModel
	err := r.db.GetContext(ctx, &menu, "SELECT "+menuColumns+" FROM sys_management_menu WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) FindAll(ctx context.Context) ([]*sysManagement.MenuModel, error) {
	var menus []*sysManagement.MenuModel
	err := 	r.db.SelectContext(ctx, &menus,
		"SELECT "+menuColumns+" FROM sys_management_menu WHERE deleted_at IS NULL ORDER BY sort ASC")
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *menuRepository) FindByParentID(ctx context.Context, parentID uint) ([]*sysManagement.MenuModel, error) {
	var menus []*sysManagement.MenuModel
	err := r.db.SelectContext(ctx, &menus,
		"SELECT "+menuColumns+" FROM sys_management_menu WHERE parent_id=$1 AND deleted_at IS NULL ORDER BY sort ASC", parentID)
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *menuRepository) FindByRoleIDs(ctx context.Context, roleIDs []uint) ([]*sysManagement.MenuModel, error) {
	if len(roleIDs) == 0 {
		return []*sysManagement.MenuModel{}, nil
	}

	query := `SELECT DISTINCT m.* FROM sys_management_menu m
		JOIN sys_management_permission p ON p.domain_id = m.id AND p.domain = 'menu'
		JOIN sys_management_role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id IN (?) AND m.deleted_at IS NULL
		ORDER BY m.sort ASC`

	query, args, err := sqlx.In(query, roleIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	var menus []*sysManagement.MenuModel
	err = r.db.SelectContext(ctx, &menus, query, args...)
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *menuRepository) Create(ctx context.Context, menu *sysManagement.MenuModel) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO sys_management_menu (menu_name, icon, path, component, redirect, parent_id, sort, hidden, keep_alive, affix, always_show, title)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		menu.MenuName, menu.Icon, menu.Path, menu.Component, menu.Redirect,
		menu.ParentID, menu.Sort, menu.Hidden, menu.KeepAlive, menu.Affix, menu.AlwaysShow, menu.Title)
	return err
}

func (r *menuRepository) Update(ctx context.Context, menu *sysManagement.MenuModel) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE sys_management_menu
		 SET menu_name=$1, icon=$2, path=$3, component=$4, redirect=$5, parent_id=$6, sort=$7,
		     hidden=$8, keep_alive=$9, affix=$10, always_show=$11, title=$12, updated_at=NOW()
		 WHERE id=$13 AND deleted_at IS NULL`,
		menu.MenuName, menu.Icon, menu.Path, menu.Component, menu.Redirect,
		menu.ParentID, menu.Sort, menu.Hidden, menu.KeepAlive, menu.Affix, menu.AlwaysShow, menu.Title,
		menu.ID)
	return err
}

func (r *menuRepository) Delete(ctx context.Context, id uint) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		"UPDATE sys_management_menu SET deleted_at=NOW() WHERE parent_id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE sys_management_menu SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
