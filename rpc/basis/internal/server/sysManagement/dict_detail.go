package sysManagement

import (
	"context"

	"td27/rpc/basis/internal/logic/sysManagement"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/dict_detail_pb"
	"td27/rpc/basis/types/sysManagement/dict_pb"
)

type DictDetailServer struct {
	svcCtx *svc.ServiceContext
	dict_detail_pb.UnimplementedDictDetailServer
}

func NewDictDetailServer(svcCtx *svc.ServiceContext) *DictDetailServer {
	return &DictDetailServer{
		svcCtx: svcCtx,
	}
}

func (ds *DictDetailServer) CreateDictDetail(ctx context.Context, in *dict_detail_pb.CreateDictDetailReq) (*dict_pb.DictDetailResp, error) {
	dl := sysManagement.NewDictDetailLogic(ctx, ds.svcCtx)
	return dl.CreateDictDetail(in)
}

func (ds *DictDetailServer) UpdateDictDetail(ctx context.Context, in *dict_detail_pb.UpdateDictDetailReq) (*dict_pb.DictDetailResp, error) {
	dl := sysManagement.NewDictDetailLogic(ctx, ds.svcCtx)
	return dl.UpdateDictDetail(in)
}

func (ds *DictDetailServer) DeleteDictDetail(ctx context.Context, in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	dl := sysManagement.NewDictDetailLogic(ctx, ds.svcCtx)
	return dl.DeleteDictDetail(in)
}

func (ds *DictDetailServer) FlatDictDetails(ctx context.Context, in *dict_detail_pb.FlatDictDetailsReq) (*dict_detail_pb.FlatDictDetailsResp, error) {
	dl := sysManagement.NewDictDetailLogic(ctx, ds.svcCtx)
	return dl.FlatDictDetails(in)
}

func (ds *DictDetailServer) ListDictDetail(ctx context.Context, in *dict_detail_pb.ListDictDetailReq) (*dict_detail_pb.ListDictDetailResp, error) {
	dl := sysManagement.NewDictDetailLogic(ctx, ds.svcCtx)
	return dl.ListDictDetail(in)
}
