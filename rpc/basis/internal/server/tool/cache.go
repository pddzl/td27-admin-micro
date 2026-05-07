package tool

import (
	"context"
	"td27/rpc/basis/internal/logic/tool"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/tool/cache_pb"
	"td27/rpc/basis/types/common_pb"
)

type CacheServer struct {
	svcCtx *svc.ServiceContext
	cache_pb.UnimplementedCacheServer
}

func NewCacheServer(svcCtx *svc.ServiceContext) *CacheServer {
	return &CacheServer{svcCtx: svcCtx}
}

func (s *CacheServer) GetCache(ctx context.Context, in *cache_pb.GetCacheReq) (*cache_pb.GetCacheResp, error) {
	return tool.NewCacheLogic(ctx, s.svcCtx).GetCache(in)
}
func (s *CacheServer) SetCache(ctx context.Context, in *cache_pb.SetCacheReq) (*common_pb.SuccessResp, error) {
	return tool.NewCacheLogic(ctx, s.svcCtx).SetCache(in)
}
func (s *CacheServer) DeleteCache(ctx context.Context, in *cache_pb.DeleteCacheReq) (*common_pb.SuccessResp, error) {
	return tool.NewCacheLogic(ctx, s.svcCtx).DeleteCache(in)
}
func (s *CacheServer) ListCache(ctx context.Context, in *common_pb.PageReq) (*cache_pb.ListCacheResp, error) {
	return tool.NewCacheLogic(ctx, s.svcCtx).ListCache(in)
}
func (s *CacheServer) CleanupExpired(ctx context.Context, in *common_pb.Empty) (*cache_pb.CleanupExpiredResp, error) {
	return tool.NewCacheLogic(ctx, s.svcCtx).CleanupExpired(in)
}
