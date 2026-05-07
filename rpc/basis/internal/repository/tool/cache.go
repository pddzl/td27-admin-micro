package tool

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"td27/rpc/basis/internal/model/tool"
	"td27/rpc/basis/internal/model/common"
)

// CacheRepository defines interface for cache data operations
type CacheRepository interface {
	FindOne(ctx context.Context, key string) (*tool.CacheModel, error)
	Create(ctx context.Context, cache *tool.CacheModel) error
	Update(ctx context.Context, cache *tool.CacheModel) error
	Delete(ctx context.Context, key string) error
	DeleteExpired(ctx context.Context) error
	List(ctx context.Context, page *common.PageInfo) ([]*tool.CacheModel, int64, error)
}

type cacheRepository struct {
	db *gorm.DB
}

// NewCacheRepository creates a new cache repository instance
func NewCacheRepository(db *gorm.DB) CacheRepository {
	return &cacheRepository{db: db}
}

func (r *cacheRepository) FindOne(ctx context.Context, key string) (*tool.CacheModel, error) {
	var cache tool.CacheModel
	if err := r.db.WithContext(ctx).Where("key = ? AND expires_at > ?", key, time.Now()).First(&cache).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cache, nil
}

func (r *cacheRepository) Create(ctx context.Context, cache *tool.CacheModel) error {
	return r.db.WithContext(ctx).Create(cache).Error
}

func (r *cacheRepository) Update(ctx context.Context, cache *tool.CacheModel) error {
	return r.db.WithContext(ctx).Model(cache).Updates(cache).Error
}

func (r *cacheRepository) Delete(ctx context.Context, key string) error {
	return r.db.WithContext(ctx).Where("key = ?", key).Delete(&tool.CacheModel{}).Error
}

func (r *cacheRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at <= ?", time.Now()).Delete(&tool.CacheModel{}).Error
}

func (r *cacheRepository) List(ctx context.Context, page *common.PageInfo) ([]*tool.CacheModel, int64, error) {
	var caches []*tool.CacheModel
	var total int64

	query := r.db.WithContext(ctx).Model(&tool.CacheModel{})
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	if err := query.Order("expires_at asc").Offset(offset).Limit(page.PageSize).Find(&caches).Error; err != nil {
		return nil, 0, err
	}

	return caches, total, nil
}
