package sysTool

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysTool"
)

// CacheRepository defines interface for cache data operations
type CacheRepository interface {
	FindOne(ctx context.Context, key string) (*sysTool.CacheModel, error)
	Create(ctx context.Context, cache *sysTool.CacheModel) error
	Update(ctx context.Context, cache *sysTool.CacheModel) error
	Delete(ctx context.Context, key string) error
	DeleteExpired(ctx context.Context) error
	List(ctx context.Context, page *common.PageInfo) ([]*sysTool.CacheModel, int64, error)
}

type cacheRepository struct {
	db *sqlx.DB
}

// NewCacheRepository creates a new cache repository instance
func NewCacheRepository(db *sqlx.DB) CacheRepository {
	return &cacheRepository{db: db}
}

const cacheColumns = `id, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at, deleted_at, username, key, value, expires_at`

func (r *cacheRepository) FindOne(ctx context.Context, key string) (*sysTool.CacheModel, error) {
	var cache sysTool.CacheModel
	err := sqlx.GetContext(ctx, r.db, &cache,
		"SELECT "+cacheColumns+" FROM sys_tool_cache WHERE key=$1 AND expires_at > $2 AND deleted_at IS NULL",
		key, time.Now())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &cache, nil
}

func (r *cacheRepository) Create(ctx context.Context, cache *sysTool.CacheModel) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO sys_tool_cache (username, key, value, expires_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		cache.Username, cache.Key, cache.Value, cache.ExpiresAt, cache.CreatedAt, cache.UpdatedAt)
	return err
}

func (r *cacheRepository) Update(ctx context.Context, cache *sysTool.CacheModel) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sys_tool_cache SET username=$1, key=$2, value=$3, expires_at=$4, updated_at=NOW() WHERE id=$5 AND deleted_at IS NULL",
		cache.Username, cache.Key, cache.Value, cache.ExpiresAt, cache.ID)
	return err
}

func (r *cacheRepository) Delete(ctx context.Context, key string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sys_tool_cache SET deleted_at=NOW() WHERE key=$1", key)
	return err
}

func (r *cacheRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sys_tool_cache SET deleted_at=NOW() WHERE expires_at <= $1", time.Now())
	return err
}

func (r *cacheRepository) List(ctx context.Context, page *common.PageInfo) ([]*sysTool.CacheModel, int64, error) {
	var total int64
	err := sqlx.GetContext(ctx, r.db, &total,
		"SELECT COUNT(*) FROM sys_tool_cache WHERE deleted_at IS NULL")
	if err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	var caches []*sysTool.CacheModel
	err = sqlx.SelectContext(ctx, r.db, &caches,
		"SELECT "+cacheColumns+" FROM sys_tool_cache WHERE deleted_at IS NULL ORDER BY expires_at ASC LIMIT $1 OFFSET $2",
		page.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return caches, total, nil
}
