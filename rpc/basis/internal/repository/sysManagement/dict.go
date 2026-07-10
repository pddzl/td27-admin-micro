package sysManagement

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
)

// DictRepository defines interface for dictionary data operations
type DictRepository interface {
	FindOne(ctx context.Context, id uint) (*sysManagement.DictModel, error)
	FindByENName(ctx context.Context, enName string) (*sysManagement.DictModel, error)
	FindAll(ctx context.Context) ([]*sysManagement.DictModel, error)
	List(ctx context.Context, page *common.PageInfo) ([]*sysManagement.DictModel, int64, error)
	Create(ctx context.Context, dict *sysManagement.DictModel) error
	Update(ctx context.Context, dict *sysManagement.DictModel) error
	Delete(ctx context.Context, id uint) error
}

type dictRepository struct {
	db *sqlx.DB
}

// NewDictRepository creates a new dictionary repository instance
func NewDictRepository(db *sqlx.DB) DictRepository {
	return &dictRepository{db: db}
}

const dictTable = "sys_management_dict"
const dictColumns = `id, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at, deleted_at, cn_name, en_name`

func (r *dictRepository) FindOne(ctx context.Context, id uint) (*sysManagement.DictModel, error) {
	var dict sysManagement.DictModel
	err := r.db.GetContext(ctx, &dict, "SELECT "+dictColumns+" FROM "+dictTable+" WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &dict, nil
}

func (r *dictRepository) FindByENName(ctx context.Context, enName string) (*sysManagement.DictModel, error) {
	var dict sysManagement.DictModel
	err := r.db.GetContext(ctx, &dict, "SELECT "+dictColumns+" FROM "+dictTable+" WHERE en_name=$1 AND deleted_at IS NULL", enName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &dict, nil
}

func (r *dictRepository) FindAll(ctx context.Context) ([]*sysManagement.DictModel, error) {
	var dicts []*sysManagement.DictModel
	err := r.db.SelectContext(ctx, &dicts, "SELECT "+dictColumns+" FROM "+dictTable+" WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	return dicts, nil
}

func (r *dictRepository) List(ctx context.Context, page *common.PageInfo) ([]*sysManagement.DictModel, int64, error) {
	var total int64
	err := r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM "+dictTable+" WHERE deleted_at IS NULL")
	if err != nil {
		return nil, 0, err
	}

	var dicts []*sysManagement.DictModel
	offset := (page.Page - 1) * page.PageSize
	err = r.db.SelectContext(ctx, &dicts,
		"SELECT "+dictColumns+" FROM "+dictTable+" WHERE deleted_at IS NULL ORDER BY id LIMIT $1 OFFSET $2",
		page.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return dicts, total, nil
}

func (r *dictRepository) Create(ctx context.Context, dict *sysManagement.DictModel) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO sys_management_dict (cn_name, en_name)
		 VALUES ($1, $2)`, dict.CNName, dict.ENName)
	return err
}

func (r *dictRepository) Update(ctx context.Context, dict *sysManagement.DictModel) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sys_management_dict SET cn_name=$1, en_name=$2, updated_at=NOW() WHERE id=$3 AND deleted_at IS NULL",
		dict.CNName, dict.ENName, dict.ID)
	return err
}

func (r *dictRepository) Delete(ctx context.Context, id uint) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		"UPDATE sys_management_dict_detail SET deleted_at=NOW() WHERE dict_id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE sys_management_dict SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
