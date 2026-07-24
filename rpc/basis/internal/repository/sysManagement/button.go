package sysManagement

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
)

// ButtonRepository defines interface for button permission data operations
type ButtonRepository interface {
	FindOne(ctx context.Context, id uint) (*sysManagement.ButtonModel, error)
	FindByCode(ctx context.Context, code string) (*sysManagement.ButtonModel, error)
	FindByPagePath(ctx context.Context, pagePath string) ([]*sysManagement.ButtonModel, error)
	FindByRoleIDs(ctx context.Context, roleIDs []uint) ([]*sysManagement.ButtonModel, error)
	List(ctx context.Context, page *common.PageInfo, pagePath *string) ([]*sysManagement.ButtonModel, int64, error)
	Create(ctx context.Context, button *sysManagement.ButtonModel) error
	Update(ctx context.Context, button *sysManagement.ButtonModel) error
	Delete(ctx context.Context, id uint) error
	BatchCheckPermission(ctx context.Context, roleIDs []uint, buttonCodes []string) (map[string]bool, error)
}

type buttonRepository struct {
	db *sqlx.DB
}

// NewButtonRepository creates a new button repository instance
func NewButtonRepository(db *sqlx.DB) ButtonRepository {
	return &buttonRepository{db: db}
}

const buttonTable = "sys_management_button"
const buttonColumns = `id, created_at, updated_at, deleted_at, button_code, button_name, description, page_path`

func (r *buttonRepository) FindOne(ctx context.Context, id uint) (*sysManagement.ButtonModel, error) {
	var button sysManagement.ButtonModel
	if err := r.db.GetContext(ctx, &button, "SELECT "+buttonColumns+" FROM "+buttonTable+" WHERE id=$1 AND deleted_at IS NULL", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &button, nil
}

func (r *buttonRepository) FindByCode(ctx context.Context, code string) (*sysManagement.ButtonModel, error) {
	var button sysManagement.ButtonModel
	if err := r.db.GetContext(ctx, &button, "SELECT "+buttonColumns+" FROM "+buttonTable+" WHERE button_code=$1 AND deleted_at IS NULL", code); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &button, nil
}

func (r *buttonRepository) FindByPagePath(ctx context.Context, pagePath string) ([]*sysManagement.ButtonModel, error) {
	var buttons []*sysManagement.ButtonModel
	if err := r.db.SelectContext(ctx, &buttons, "SELECT "+buttonColumns+" FROM "+buttonTable+" WHERE page_path=$1 AND deleted_at IS NULL", pagePath); err != nil {
		return nil, err
	}
	return buttons, nil
}

func (r *buttonRepository) FindByRoleIDs(ctx context.Context, roleIDs []uint) ([]*sysManagement.ButtonModel, error) {
	query, args, err := sqlx.In(`
		SELECT DISTINCT b.id, b.created_at as created_at, b.updated_at as updated_at, b.deleted_at, b.button_code, b.button_name, b.description, b.page_path FROM `+buttonTable+` b
		JOIN sys_management_permission p ON p.domain_id = b.id AND p.domain = 'button'
		JOIN sys_management_role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id IN (?) AND b.deleted_at IS NULL`, roleIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	var buttons []*sysManagement.ButtonModel
	if err = r.db.SelectContext(ctx, &buttons, query, args...); err != nil {
		return nil, err
	}
	return buttons, nil
}

func (r *buttonRepository) List(ctx context.Context, page *common.PageInfo, pagePath *string) ([]*sysManagement.ButtonModel, int64, error) {
	var total int64
	var buttons []*sysManagement.ButtonModel

	where := "WHERE deleted_at IS NULL"
	var args []interface{}

	if pagePath != nil {
		where += " AND page_path = ?"
		args = append(args, *pagePath)
	}

	countQuery := "SELECT COUNT(*) FROM " + buttonTable + " " + where
	countQuery = r.db.Rebind(countQuery)
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	dataQuery := "SELECT " + buttonColumns + " FROM " + buttonTable + " " + where + " ORDER BY id LIMIT ? OFFSET ?"
	dataQuery = r.db.Rebind(dataQuery)
	dataArgs := append(args, page.PageSize, offset)
	if err := r.db.SelectContext(ctx, &buttons, dataQuery, dataArgs...); err != nil {
		return nil, 0, err
	}

	return buttons, total, nil
}

func (r *buttonRepository) Create(ctx context.Context, button *sysManagement.ButtonModel) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO `+buttonTable+` (created_at, updated_at, button_code, button_name, description, page_path)
		 VALUES (NOW(), NOW(), $1, $2, $3, $4) RETURNING id`,
		button.ButtonCode, button.ButtonName, button.Description, button.PagePath).Scan(&button.ID)
	return err
}

func (r *buttonRepository) Update(ctx context.Context, button *sysManagement.ButtonModel) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE `+buttonTable+` SET button_code=$1, button_name=$2, description=$3, page_path=$4, updated_at=NOW() WHERE id=$5`,
		button.ButtonCode, button.ButtonName, button.Description, button.PagePath, button.ID)
	return err
}

func (r *buttonRepository) BatchCheckPermission(ctx context.Context, roleIDs []uint, buttonCodes []string) (map[string]bool, error) {
	result := make(map[string]bool, len(buttonCodes))
	for _, code := range buttonCodes {
		result[code] = false
	}

	if len(roleIDs) == 0 || len(buttonCodes) == 0 {
		return result, nil
	}

	query, args, err := sqlx.In(`
		SELECT DISTINCT p.resource FROM sys_management_permission p
		JOIN sys_management_role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id IN (?) AND p.domain = 'button' AND p.resource IN (?)`, roleIDs, buttonCodes)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	var permittedCodes []string
	if err = r.db.SelectContext(ctx, &permittedCodes, query, args...); err != nil {
		return nil, err
	}

	for _, code := range permittedCodes {
		result[code] = true
	}

	return result, nil
}

func (r *buttonRepository) Delete(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM "+buttonTable+" WHERE id=$1", id)
	return err
}
