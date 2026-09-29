package sysTool

import (
	"context"
	"errors"
	"time"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysTool"
	sysToolRepo "td27/rpc/basis/internal/repository/sysTool"
)

type CacheService struct {
	cacheRepo sysToolRepo.CacheRepository
}

func NewCacheService(cacheRepo sysToolRepo.CacheRepository) *CacheService {
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
	cache := &sysTool.CacheModel{
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

func (s *CacheService) List(ctx context.Context, page *common.PageInfo) ([]*sysTool.CacheModel, int64, error) {
	return s.cacheRepo.List(ctx, page)
}
