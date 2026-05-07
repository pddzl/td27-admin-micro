package tool

import (
	"context"
	"td27/rpc/basis/internal/logic/tool"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/tool/service_token_pb"
	"td27/rpc/basis/types/common_pb"
)

type ServiceTokenServer struct {
	svcCtx *svc.ServiceContext
	service_token_pb.UnimplementedServiceTokenServer
}

func NewServiceTokenServer(svcCtx *svc.ServiceContext) *ServiceTokenServer {
	return &ServiceTokenServer{svcCtx: svcCtx}
}

func (s *ServiceTokenServer) CreateServiceToken(ctx context.Context, in *service_token_pb.CreateServiceTokenReq) (*service_token_pb.CreateServiceTokenResp, error) {
	return tool.NewServiceTokenLogic(ctx, s.svcCtx).CreateServiceToken(in)
}
func (s *ServiceTokenServer) GetServiceToken(ctx context.Context, in *common_pb.IdReq) (*service_token_pb.ServiceTokenResp, error) {
	return tool.NewServiceTokenLogic(ctx, s.svcCtx).GetServiceToken(in)
}
func (s *ServiceTokenServer) ListServiceToken(ctx context.Context, in *common_pb.PageReq) (*service_token_pb.ListServiceTokenResp, error) {
	return tool.NewServiceTokenLogic(ctx, s.svcCtx).ListServiceToken(in)
}
func (s *ServiceTokenServer) UpdateServiceToken(ctx context.Context, in *service_token_pb.UpdateServiceTokenReq) (*service_token_pb.ServiceTokenResp, error) {
	return tool.NewServiceTokenLogic(ctx, s.svcCtx).UpdateServiceToken(in)
}
func (s *ServiceTokenServer) ToggleTokenStatus(ctx context.Context, in *service_token_pb.ToggleTokenStatusReq) (*common_pb.SuccessResp, error) {
	return tool.NewServiceTokenLogic(ctx, s.svcCtx).ToggleTokenStatus(in)
}
func (s *ServiceTokenServer) AssignTokenPermissions(ctx context.Context, in *service_token_pb.AssignTokenPermissionsReq) (*common_pb.SuccessResp, error) {
	return tool.NewServiceTokenLogic(ctx, s.svcCtx).AssignTokenPermissions(in)
}
func (s *ServiceTokenServer) GetTokenPermissions(ctx context.Context, in *common_pb.IdReq) (*service_token_pb.GetTokenPermissionsResp, error) {
	return tool.NewServiceTokenLogic(ctx, s.svcCtx).GetTokenPermissions(in)
}
func (s *ServiceTokenServer) DeleteServiceToken(ctx context.Context, in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	return tool.NewServiceTokenLogic(ctx, s.svcCtx).DeleteServiceToken(in)
}
func (s *ServiceTokenServer) ValidateToken(ctx context.Context, in *service_token_pb.ValidateTokenReq) (*service_token_pb.ValidateTokenResp, error) {
	return tool.NewServiceTokenLogic(ctx, s.svcCtx).ValidateToken(in)
}
