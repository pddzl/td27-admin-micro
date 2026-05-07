package tool

import (
	"context"
	"errors"
	"time"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/tool"
	repoTool "td27/rpc/basis/internal/repository/tool"
)

type CacheService struct {
	cacheRepo repoTool.CacheRepository
}

func NewCacheService(cacheRepo repoTool.CacheRepository) *CacheService {
	return &CacheService{
		cacheRepo: cacheRepo,
	}
}

func (s *CacheService) Get(ctx context.Context, key string) (string, error) {
	cache, err := s.cacheRepo.FindOne(ctx, key)
	if err != nil {
		return "", err
	}
	if cache == nil {
		return "", errors.New("cache not found")
	}
	return cache.Value, nil
}

func (s *CacheService) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	expiresAt := time.Now().Add(ttl)
	cache := &tool.CacheModel{
		Key:       key,
		Value:     value,
		ExpiresAt: expiresAt,
	}
	return s.cacheRepo.Create(ctx, cache)
}

func (s *CacheService) Delete(ctx context.Context, key string) error {
	return s.cacheRepo.Delete(ctx, key)
}

func (s *CacheService) CleanupExpired(ctx context.Context) error {
	return s.cacheRepo.DeleteExpired(ctx)
}

func (s *CacheService) List(ctx context.Context, page *common.PageInfo) ([]*tool.CacheModel, int64, error) {
	return s.cacheRepo.List(ctx, page)
}
