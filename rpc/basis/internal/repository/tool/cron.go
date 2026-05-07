package tool

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"td27/rpc/basis/internal/model/tool"
	"td27/rpc/basis/internal/model/common"
)

// CronRepository defines interface for scheduled task data operations
type CronRepository interface {
	FindOne(ctx context.Context, id uint) (*tool.CronModel, error)
	FindAllEnabled(ctx context.Context) ([]*tool.CronModel, error)
	List(ctx context.Context, page *common.PageInfo) ([]*tool.CronModel, int64, error)
	Create(ctx context.Context, cron *tool.CronModel) error
	Update(ctx context.Context, cron *tool.CronModel) error
	UpdateEntryID(ctx context.Context, id uint, entryID int) error
	ToggleStatus(ctx context.Context, id uint, open bool) error
	Delete(ctx context.Context, id uint) error
}

type cronRepository struct {
	db *gorm.DB
}

// NewCronRepository creates a new cron repository instance
func NewCronRepository(db *gorm.DB) CronRepository {
	return &cronRepository{db: db}
}

func (r *cronRepository) FindOne(ctx context.Context, id uint) (*tool.CronModel, error) {
	var cron tool.CronModel
	if err := r.db.WithContext(ctx).First(&cron, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cron, nil
}

func (r *cronRepository) FindAllEnabled(ctx context.Context) ([]*tool.CronModel, error) {
	var crons []*tool.CronModel
	if err := r.db.WithContext(ctx).Where("open = ?", true).Find(&crons).Error; err != nil {
		return nil, err
	}
	return crons, nil
}

func (r *cronRepository) List(ctx context.Context, page *common.PageInfo) ([]*tool.CronModel, int64, error) {
	var crons []*tool.CronModel
	var total int64

	query := r.db.WithContext(ctx).Model(&tool.CronModel{})
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	if err := query.Order("created_at desc").Offset(offset).Limit(page.PageSize).Find(&crons).Error; err != nil {
		return nil, 0, err
	}

	return crons, total, nil
}

func (r *cronRepository) Create(ctx context.Context, cron *tool.CronModel) error {
	return r.db.WithContext(ctx).Create(cron).Error
}

func (r *cronRepository) Update(ctx context.Context, cron *tool.CronModel) error {
	return r.db.WithContext(ctx).Model(cron).Updates(cron).Error
}

func (r *cronRepository) UpdateEntryID(ctx context.Context, id uint, entryID int) error {
	return r.db.WithContext(ctx).Model(&tool.CronModel{}).Where("id = ?", id).Update("entry_id", entryID).Error
}

func (r *cronRepository) ToggleStatus(ctx context.Context, id uint, open bool) error {
	return r.db.WithContext(ctx).Model(&tool.CronModel{}).Where("id = ?", id).Update("open", open).Error
}

func (r *cronRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&tool.CronModel{}, id).Error
}
