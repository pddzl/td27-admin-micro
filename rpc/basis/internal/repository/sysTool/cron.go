package sysTool

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysTool"
)

// CronRepository defines interface for scheduled task data operations
type CronRepository interface {
	FindOne(ctx context.Context, id uint) (*sysTool.CronModel, error)
	FindAllEnabled(ctx context.Context) ([]*sysTool.CronModel, error)
	List(ctx context.Context, page *common.PageInfo) ([]*sysTool.CronModel, int64, error)
	Create(ctx context.Context, cron *sysTool.CronModel) error
	Update(ctx context.Context, cron *sysTool.CronModel) error
	UpdateEntryID(ctx context.Context, id uint, entryID int) error
	ToggleStatus(ctx context.Context, id uint, open bool) error
	Delete(ctx context.Context, id uint) error
	DeleteByIds(ctx context.Context, ids []uint) error
}

type cronRepository struct {
	db *sqlx.DB
}

// NewCronRepository creates a new cron repository instance
func NewCronRepository(db *sqlx.DB) CronRepository {
	return &cronRepository{db: db}
}

const cronColumns = `id, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at, deleted_at, name, method, expression, strategy, open, extraParams, entryId, comment`

func (r *cronRepository) FindOne(ctx context.Context, id uint) (*sysTool.CronModel, error) {
	var cron sysTool.CronModel
	err := sqlx.GetContext(ctx, r.db, &cron,
		"SELECT "+cronColumns+" FROM sys_tool_cron WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &cron, nil
}

func (r *cronRepository) FindAllEnabled(ctx context.Context) ([]*sysTool.CronModel, error) {
	var crons []*sysTool.CronModel
	err := sqlx.SelectContext(ctx, r.db, &crons,
		"SELECT "+cronColumns+" FROM sys_tool_cron WHERE open=$1 AND deleted_at IS NULL", true)
	if err != nil {
		return nil, err
	}
	return crons, nil
}

func (r *cronRepository) List(ctx context.Context, page *common.PageInfo) ([]*sysTool.CronModel, int64, error) {
	var total int64
	err := sqlx.GetContext(ctx, r.db, &total,
		"SELECT COUNT(*) FROM sys_tool_cron WHERE deleted_at IS NULL")
	if err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	var crons []*sysTool.CronModel
	err = sqlx.SelectContext(ctx, r.db, &crons,
		"SELECT "+cronColumns+" FROM sys_tool_cron WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		page.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return crons, total, nil
}

func (r *cronRepository) Create(ctx context.Context, cron *sysTool.CronModel) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO sys_tool_cron (name, method, expression, strategy, open, extraParams, entryId, comment, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		cron.Name, cron.Method, cron.Expression, cron.Strategy, cron.Open, cron.ExtraParams, cron.EntryId, cron.Comment, cron.CreatedAt, cron.UpdatedAt)
	return err
}

func (r *cronRepository) Update(ctx context.Context, cron *sysTool.CronModel) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sys_tool_cron SET name=$1, method=$2, expression=$3, strategy=$4, open=$5, extraParams=$6, entryId=$7, comment=$8, updated_at=NOW() WHERE id=$9 AND deleted_at IS NULL",
		cron.Name, cron.Method, cron.Expression, cron.Strategy, cron.Open, cron.ExtraParams, cron.EntryId, cron.Comment, cron.ID)
	return err
}

func (r *cronRepository) UpdateEntryID(ctx context.Context, id uint, entryID int) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sys_tool_cron SET entryId=$1, updated_at=NOW() WHERE id=$2 AND deleted_at IS NULL",
		entryID, id)
	return err
}

func (r *cronRepository) ToggleStatus(ctx context.Context, id uint, open bool) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sys_tool_cron SET open=$1, updated_at=NOW() WHERE id=$2 AND deleted_at IS NULL",
		open, id)
	return err
}

func (r *cronRepository) Delete(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sys_tool_cron SET deleted_at=NOW() WHERE id=$1", id)
	return err
}

func (r *cronRepository) DeleteByIds(ctx context.Context, ids []uint) error {
	query, args, err := sqlx.In("UPDATE sys_tool_cron SET deleted_at=NOW() WHERE id IN (?)", ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
