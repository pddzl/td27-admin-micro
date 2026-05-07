package monitor

import (
	"context"
	"time"

	"gorm.io/gorm"

	"td27/rpc/basis/internal/model/monitor"
	"td27/rpc/basis/internal/model/common"
)

// OperationLogRepository defines interface for operation log data operations
type OperationLogRepository interface {
	Create(ctx context.Context, log *monitor.OperationLogModel) error
	List(ctx context.Context, page *common.PageInfo, userID *uint, status *int) ([]*monitor.OperationLogModel, int64, error)
	DeleteExpired(ctx context.Context, days int) error
}

type operationLogRepository struct {
	db *gorm.DB
}

// NewOperationLogRepository creates a new operation log repository instance
func NewOperationLogRepository(db *gorm.DB) OperationLogRepository {
	return &operationLogRepository{db: db}
}

func (r *operationLogRepository) Create(ctx context.Context, log *monitor.OperationLogModel) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *operationLogRepository) List(ctx context.Context, page *common.PageInfo, userID *uint, status *int) ([]*monitor.OperationLogModel, int64, error) {
	var logs []*monitor.OperationLogModel
	var total int64

	query := r.db.WithContext(ctx).Model(&monitor.OperationLogModel{})
	
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	if err := query.Order("created_at desc").Offset(offset).Limit(page.PageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *operationLogRepository) DeleteExpired(ctx context.Context, days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	return r.db.WithContext(ctx).Where("created_at < ?", cutoff).Delete(&monitor.OperationLogModel{}).Error
}
