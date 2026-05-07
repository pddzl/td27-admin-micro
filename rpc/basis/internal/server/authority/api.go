package authority

import (
	"context"

	"td27/rpc/basis/internal/logic/authority"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/authority/api_pb"
	"td27/rpc/basis/types/common_pb"
)

type APIServer struct {
	svcCtx *svc.ServiceContext
	api_pb.UnimplementedAPIServer
}

func NewAPIServer(svcCtx *svc.ServiceContext) *APIServer {
	return &APIServer{
		svcCtx: svcCtx,
	}
}

func (as *APIServer) GetAPI(ctx context.Context, in *common_pb.IdReq) (*api_pb.APIResp, error) {
	al := authority.NewAPILogic(ctx, as.svcCtx)
	return al.GetAPI(in)
}

func (as *APIServer) ListAPI(ctx context.Context, in *common_pb.PageReq) (*api_pb.ListAPIResp, error) {
	al := authority.NewAPILogic(ctx, as.svcCtx)
	return al.ListAPI(in)
}

func (as *APIServer) GetAPIsByGroup(ctx context.Context, in *api_pb.GetAPIsByGroupReq) (*api_pb.ListAPIResp, error) {
	al := authority.NewAPILogic(ctx, as.svcCtx)
	return al.GetAPIsByGroup(in)
}

func (as *APIServer) CreateAPI(ctx context.Context, in *api_pb.CreateAPIReq) (*common_pb.SuccessResp, error) {
	al := authority.NewAPILogic(ctx, as.svcCtx)
	return al.CreateAPI(in)
}

func (as *APIServer) UpdateAPI(ctx context.Context, in *api_pb.UpdateAPIReq) (*api_pb.APIResp, error) {
	al := authority.NewAPILogic(ctx, as.svcCtx)
	return al.UpdateAPI(in)
}

func (as *APIServer) DeleteAPI(ctx context.Context, in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	al := authority.NewAPILogic(ctx, as.svcCtx)
	return al.DeleteAPI(in)
}
