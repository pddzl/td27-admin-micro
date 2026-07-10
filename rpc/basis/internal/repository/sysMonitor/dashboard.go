package sysMonitor

import (
	"context"

	"github.com/jmoiron/sqlx"

	sysMonitorModel "td27/rpc/basis/internal/model/sysMonitor"
)

// DashboardRepository defines interface for dashboard data aggregation
type DashboardRepository interface {
	GetUserCount(ctx context.Context) (int64, error)
	GetRoleCount(ctx context.Context) (int64, error)
	GetAPICount(ctx context.Context) (int64, error)
	GetDeptCount(ctx context.Context) (int64, error)
	GetMenuCount(ctx context.Context) (int64, error)
	GetDictCount(ctx context.Context) (int64, error)
	GetLogCount(ctx context.Context) (int64, error)
	GetFileCount(ctx context.Context) (int64, error)
	GetRecentOperations(ctx context.Context, limit int) ([]*sysMonitorModel.OperationLogModel, error)
}

type dashboardRepository struct {
	db *sqlx.DB
}

// NewDashboardRepository creates a new dashboard repository instance
func NewDashboardRepository(db *sqlx.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetUserCount(ctx context.Context) (int64, error) {
	var count int64
	err := sqlx.GetContext(ctx, r.db, &count,
		"SELECT COUNT(*) FROM sys_management_user WHERE deleted_at IS NULL")
	return count, err
}

func (r *dashboardRepository) GetRoleCount(ctx context.Context) (int64, error) {
	var count int64
	err := sqlx.GetContext(ctx, r.db, &count,
		"SELECT COUNT(*) FROM sys_management_role WHERE deleted_at IS NULL")
	return count, err
}

func (r *dashboardRepository) GetAPICount(ctx context.Context) (int64, error) {
	var count int64
	err := sqlx.GetContext(ctx, r.db, &count,
		"SELECT COUNT(*) FROM sys_management_api WHERE deleted_at IS NULL")
	return count, err
}

func (r *dashboardRepository) GetDeptCount(ctx context.Context) (int64, error) {
	var count int64
	err := sqlx.GetContext(ctx, r.db, &count,
		"SELECT COUNT(*) FROM sys_management_dept WHERE deleted_at IS NULL")
	return count, err
}

func (r *dashboardRepository) GetMenuCount(ctx context.Context) (int64, error) {
	var count int64
	err := sqlx.GetContext(ctx, r.db, &count,
		"SELECT COUNT(*) FROM sys_management_menu WHERE deleted_at IS NULL")
	return count, err
}

func (r *dashboardRepository) GetDictCount(ctx context.Context) (int64, error) {
	var count int64
	err := sqlx.GetContext(ctx, r.db, &count,
		"SELECT COUNT(*) FROM sys_management_dict WHERE deleted_at IS NULL")
	return count, err
}

func (r *dashboardRepository) GetLogCount(ctx context.Context) (int64, error) {
	var count int64
	err := sqlx.GetContext(ctx, r.db, &count,
		"SELECT COUNT(*) FROM sys_monitor_operation_log WHERE deleted_at IS NULL")
	return count, err
}

func (r *dashboardRepository) GetFileCount(ctx context.Context) (int64, error) {
	var count int64
	err := sqlx.GetContext(ctx, r.db, &count,
		"SELECT COUNT(*) FROM sys_tool_file WHERE deleted_at IS NULL")
	return count, err
}

func (r *dashboardRepository) GetRecentOperations(ctx context.Context, limit int) ([]*sysMonitorModel.OperationLogModel, error) {
	var logs []*sysMonitorModel.OperationLogModel
	err := sqlx.SelectContext(ctx, r.db, &logs,
		"SELECT "+operationLogColumns+" FROM sys_monitor_operation_log WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1", limit)
	return logs, err
}
