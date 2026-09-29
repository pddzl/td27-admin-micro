package sysManagement

import (
	"context"

	"td27/rpc/basis/internal/logic/sysManagement"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/dict_pb"
)

type DictServer struct {
	svcCtx *svc.ServiceContext
	dict_pb.UnimplementedDictServer
}

func NewDictServer(svcCtx *svc.ServiceContext) *DictServer {
	return &DictServer{
		svcCtx: svcCtx,
	}
}

func (ds *DictServer) GetDict(ctx context.Context, in *common_pb.IdReq) (*dict_pb.DictResp, error) {
	dl := sysManagement.NewDictLogic(ctx, ds.svcCtx)
	return dl.GetDict(in)
}

func (ds *DictServer) GetDictByENName(ctx context.Context, in *dict_pb.GetDictByENNameReq) (*dict_pb.DictResp, error) {
	dl := sysManagement.NewDictLogic(ctx, ds.svcCtx)
	return dl.GetDictByENName(in)
}

func (ds *DictServer) ListDict(ctx context.Context, in *common_pb.PageReq) (*dict_pb.ListDictResp, error) {
	dl := sysManagement.NewDictLogic(ctx, ds.svcCtx)
	return dl.ListDict(in)
}

func (ds *DictServer) CreateDict(ctx context.Context, in *dict_pb.CreateDictReq) (*common_pb.SuccessResp, error) {
	dl := sysManagement.NewDictLogic(ctx, ds.svcCtx)
	return dl.CreateDict(in)
}

func (ds *DictServer) UpdateDict(ctx context.Context, in *dict_pb.UpdateDictReq) (*dict_pb.DictResp, error) {
	dl := sysManagement.NewDictLogic(ctx, ds.svcCtx)
	return dl.UpdateDict(in)
}

func (ds *DictServer) DeleteDict(ctx context.Context, in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	dl := sysManagement.NewDictLogic(ctx, ds.svcCtx)
	return dl.DeleteDict(in)
}
