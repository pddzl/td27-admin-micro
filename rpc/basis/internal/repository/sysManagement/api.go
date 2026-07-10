package sysManagement

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
)

// APIRepository defines interface for API data operations
type APIRepository interface {
	FindOne(ctx context.Context, id uint) (*sysManagement.ApiModel, error)
	FindAll(ctx context.Context) ([]*sysManagement.ApiModel, error)
	FindByPathAndMethod(ctx context.Context, path string, method string) (*sysManagement.ApiModel, error)
	FindByGroup(ctx context.Context, groupEN string) ([]*sysManagement.ApiModel, error)
	List(ctx context.Context, page *common.PageInfo, groupEN *string) ([]*sysManagement.ApiModel, int64, error)
	Create(ctx context.Context, api *sysManagement.ApiModel) error
	Update(ctx context.Context, api *sysManagement.ApiModel) error
	Delete(ctx context.Context, id uint) error
	DeleteByIds(ctx context.Context, ids []uint) error
}

const apiColumns = `id, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at, deleted_at, COALESCE(path, '') as path, COALESCE(method, '') as method, COALESCE(group_en, '') as group_en, COALESCE(group_cn, '') as group_cn, COALESCE(description, '') as description`

type apiRepository struct {
	db *sqlx.DB
}

// NewAPIRepository creates a new API repository instance
func NewAPIRepository(db *sqlx.DB) APIRepository {
	return &apiRepository{db: db}
}

const apiTable = "sys_management_api"

func (r *apiRepository) FindOne(ctx context.Context, id uint) (*sysManagement.ApiModel, error) {
	var api sysManagement.ApiModel
	if err := r.db.GetContext(ctx, &api, "SELECT " + apiColumns + " FROM " + apiTable + " WHERE id=$1 AND deleted_at IS NULL", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &api, nil
}

func (r *apiRepository) FindAll(ctx context.Context) ([]*sysManagement.ApiModel, error) {
	var apis []*sysManagement.ApiModel
	if err := r.db.SelectContext(ctx, &apis, "SELECT " + apiColumns + " FROM " + apiTable + " WHERE deleted_at IS NULL"); err != nil {
		return nil, err
	}
	return apis, nil
}

func (r *apiRepository) FindByPathAndMethod(ctx context.Context, path string, method string) (*sysManagement.ApiModel, error) {
	var api sysManagement.ApiModel
	if err := r.db.GetContext(ctx, &api, "SELECT " + apiColumns + " FROM " + apiTable + " WHERE path=$1 AND method=$2 AND deleted_at IS NULL", path, method); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &api, nil
}

func (r *apiRepository) FindByGroup(ctx context.Context, groupEN string) ([]*sysManagement.ApiModel, error) {
	var apis []*sysManagement.ApiModel
	if err := r.db.SelectContext(ctx, &apis, "SELECT " + apiColumns + " FROM " + apiTable + " WHERE group_en=$1 AND deleted_at IS NULL", groupEN); err != nil {
		return nil, err
	}
	return apis, nil
}

func (r *apiRepository) List(ctx context.Context, page *common.PageInfo, groupEN *string) ([]*sysManagement.ApiModel, int64, error) {
	var total int64
	var apis []*sysManagement.ApiModel

	where := "WHERE deleted_at IS NULL"
	var args []interface{}

	if groupEN != nil {
		where += " AND group_en = ?"
		args = append(args, *groupEN)
	}

	countQuery := "SELECT COUNT(*) FROM " + apiTable + " " + where
	countQuery = r.db.Rebind(countQuery)
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	dataQuery := "SELECT " + apiColumns + " FROM " + apiTable + " " + where + " ORDER BY id LIMIT ? OFFSET ?"
	dataQuery = r.db.Rebind(dataQuery)
	dataArgs := append(args, page.PageSize, offset)
	if err := r.db.SelectContext(ctx, &apis, dataQuery, dataArgs...); err != nil {
		return nil, 0, err
	}

	return apis, total, nil
}

func (r *apiRepository) Create(ctx context.Context, api *sysManagement.ApiModel) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO `+apiTable+` (path, method, group_en, group_cn, description, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		api.Path, api.Method, api.GroupEN, api.GroupCN, api.Description, api.CreatedAt, api.UpdatedAt).Scan(&api.ID)
	return err
}

func (r *apiRepository) Update(ctx context.Context, api *sysManagement.ApiModel) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE `+apiTable+` SET path=$1, method=$2, group_en=$3, group_cn=$4, description=$5, updated_at=$6 WHERE id=$7`,
		api.Path, api.Method, api.GroupEN, api.GroupCN, api.Description, api.UpdatedAt, api.ID)
	return err
}

func (r *apiRepository) Delete(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx, "UPDATE "+apiTable+" SET deleted_at=NOW() WHERE id=$1", id)
	return err
}

func (r *apiRepository) DeleteByIds(ctx context.Context, ids []uint) error {
	query, args, err := sqlx.In("UPDATE "+apiTable+" SET deleted_at=NOW() WHERE id IN (?) AND deleted_at IS NULL", ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
