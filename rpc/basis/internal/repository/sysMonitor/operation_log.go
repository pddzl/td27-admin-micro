package sysMonitor

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysMonitor"
)

// OperationLogRepository defines interface for operation log data operations
type OperationLogRepository interface {
	Create(ctx context.Context, log *sysMonitor.OperationLogModel) error
	List(ctx context.Context, page *common.PageInfo, userID *uint, status *int) ([]*sysMonitor.OperationLogModel, int64, error)
	DeleteExpired(ctx context.Context, days int) error
	Delete(ctx context.Context, id uint) error
	DeleteByIds(ctx context.Context, ids []uint) error
}

type operationLogRepository struct {
	db *sqlx.DB
}

// NewOperationLogRepository creates a new operation log repository instance
func NewOperationLogRepository(db *sqlx.DB) OperationLogRepository {
	return &operationLogRepository{db: db}
}

func (r *operationLogRepository) Create(ctx context.Context, log *sysMonitor.OperationLogModel) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO sys_monitor_operation_log (ip, method, path, status, user_agent, req_param, resp_data, resp_time, user_id, user_name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		log.Ip, log.Method, log.Path, log.Status, log.UserAgent, log.ReqParam, log.RespData, log.RespTime, log.UserID, log.UserName, log.CreatedAt, log.UpdatedAt)
	return err
}

const operationLogColumns = `id, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at, deleted_at, ip, method, path, status, user_agent, req_param, resp_data, resp_time, user_id, user_name`

func (r *operationLogRepository) List(ctx context.Context, page *common.PageInfo, userID *uint, status *int) ([]*sysMonitor.OperationLogModel, int64, error) {
	baseWhere := "WHERE deleted_at IS NULL"
	args := []interface{}{}

	if userID != nil {
		baseWhere += fmt.Sprintf(" AND user_id=$%d", len(args)+1)
		args = append(args, *userID)
	}
	if status != nil {
		baseWhere += fmt.Sprintf(" AND status=$%d", len(args)+1)
		args = append(args, *status)
	}

	var total int64
	countQuery := "SELECT COUNT(*) FROM sys_monitor_operation_log " + baseWhere
	err := sqlx.GetContext(ctx, r.db, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	dataQuery := "SELECT " + operationLogColumns + " FROM sys_monitor_operation_log " + baseWhere + " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	dataArgs := append(args, page.PageSize, offset)

	var logs []*sysMonitor.OperationLogModel
	err = sqlx.SelectContext(ctx, r.db, &logs, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *operationLogRepository) DeleteExpired(ctx context.Context, days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	_, err := r.db.ExecContext(ctx,
		"UPDATE sys_monitor_operation_log SET deleted_at=NOW() WHERE created_at < $1", cutoff)
	return err
}

func (r *operationLogRepository) Delete(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sys_monitor_operation_log SET deleted_at=NOW() WHERE id=$1", id)
	return err
}

func (r *operationLogRepository) DeleteByIds(ctx context.Context, ids []uint) error {
	query, args, err := sqlx.In("UPDATE sys_monitor_operation_log SET deleted_at=NOW() WHERE id IN (?)", ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
