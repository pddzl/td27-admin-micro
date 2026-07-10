package sysTool

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/internal/util"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysTool/cache_pb"
)

type CacheLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCacheLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CacheLogic {
	return &CacheLogic{
		ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx),
	}
}

func (cl *CacheLogic) GetCache(in *cache_pb.GetCacheReq) (*cache_pb.GetCacheResp, error) {
	value, err := cl.svcCtx.CacheService.Get(cl.ctx, in.Key)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get cache failed: %v", err)
	}
	return &cache_pb.GetCacheResp{Key: in.Key, Value: value}, nil
}

func (cl *CacheLogic) SetCache(in *cache_pb.SetCacheReq) (*common_pb.SuccessResp, error) {
	if err := cl.svcCtx.CacheService.Set(cl.ctx, in.Key, in.Value, time.Duration(in.TtlSeconds)*time.Second); err != nil {
		return nil, status.Errorf(codes.Internal, "set cache failed: %v", err)
	}
	return &common_pb.SuccessResp{Success: true}, nil
}

func (cl *CacheLogic) DeleteCache(in *cache_pb.DeleteCacheReq) (*common_pb.SuccessResp, error) {
	if err := cl.svcCtx.CacheService.Delete(cl.ctx, in.Key); err != nil {
		return nil, status.Errorf(codes.Internal, "delete cache failed: %v", err)
	}
	return &common_pb.SuccessResp{Success: true}, nil
}

func (cl *CacheLogic) ListCache(in *common_pb.PageReq) (*cache_pb.ListCacheResp, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}

	page := &common.PageInfo{Page: int(in.Page), PageSize: int(in.PageSize)}
	caches, count, err := cl.svcCtx.CacheService.List(cl.ctx, page)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list cache failed: %v", err)
	}
	resp := &cache_pb.ListCacheResp{List: make([]*cache_pb.CacheResp, 0, len(caches)), Total: count	}
	for _, c := range caches {
		resp.List = append(resp.List, &cache_pb.CacheResp{
			Id: int64(c.ID), Username: c.Username, Key: c.Key, Value: c.Value,
			ExpiresAt: util.ToProtoTimestamp(c.ExpiresAt),
			CreatedAt: util.ToProtoTimestamp(c.CreatedAt), UpdatedAt: util.ToProtoTimestamp(c.UpdatedAt),
		})
	}
	return resp, nil
}

func (cl *CacheLogic) CleanupExpired(in *common_pb.Empty) (*cache_pb.CleanupExpiredResp, error) {
	if err := cl.svcCtx.CacheService.CleanupExpired(cl.ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "cleanup cache failed: %v", err)
	}
	return &cache_pb.CleanupExpiredResp{DeletedCount: 0}, nil
}
