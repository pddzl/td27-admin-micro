package sysManagement

import (
	"context"

	"td27/rpc/basis/internal/logic/sysManagement"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/role_pb"
)

type RoleServer struct {
	svcCtx *svc.ServiceContext
	role_pb.UnimplementedRoleServer
}

func NewRoleServer(svcCtx *svc.ServiceContext) *RoleServer {
	return &RoleServer{
		svcCtx: svcCtx,
	}
}

func (rs *RoleServer) GetRole(ctx context.Context, in *common_pb.IdReq) (*role_pb.RoleResp, error) {
	rl := sysManagement.NewRoleLogic(ctx, rs.svcCtx)
	return rl.GetRole(in)
}

func (rs *RoleServer) GetRoleWithChildren(ctx context.Context, in *common_pb.IdReq) (*role_pb.RoleTreeResp, error) {
	rl := sysManagement.NewRoleLogic(ctx, rs.svcCtx)
	return rl.GetRoleWithChildren(in)
}

func (rs *RoleServer) ListRole(ctx context.Context, in *common_pb.PageReq) (*role_pb.ListRoleResp, error) {
	rl := sysManagement.NewRoleLogic(ctx, rs.svcCtx)
	return rl.ListRole(in)
}

func (rs *RoleServer) GetAllRoles(ctx context.Context, in *common_pb.Empty) (*role_pb.ListRoleResp, error) {
	rl := sysManagement.NewRoleLogic(ctx, rs.svcCtx)
	return rl.GetAllRoles(in)
}

func (rs *RoleServer) CreateRole(ctx context.Context, in *role_pb.CreateRoleReq) (*common_pb.SuccessResp, error) {
	rl := sysManagement.NewRoleLogic(ctx, rs.svcCtx)
	return rl.CreateRole(in)
}

func (rs *RoleServer) UpdateRole(ctx context.Context, in *role_pb.UpdateRoleReq) (*role_pb.RoleResp, error) {
	rl := sysManagement.NewRoleLogic(ctx, rs.svcCtx)
	return rl.UpdateRole(in)
}

func (rs *RoleServer) DeleteRole(ctx context.Context, in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	rl := sysManagement.NewRoleLogic(ctx, rs.svcCtx)
	return rl.DeleteRole(in)
}

func (rs *RoleServer) AssignPermissions(ctx context.Context, in *role_pb.AssignPermissionsReq) (*common_pb.SuccessResp, error) {
	rl := sysManagement.NewRoleLogic(ctx, rs.svcCtx)
	return rl.AssignPermissions(in)
}

func (rs *RoleServer) GetRolePermissions(ctx context.Context, in *common_pb.IdReq) (*role_pb.GetRolePermissionsResp, error) {
	rl := sysManagement.NewRoleLogic(ctx, rs.svcCtx)
	return rl.GetRolePermissions(in)
}

func (rs *RoleServer) CreateElevationRequest(ctx context.Context, in *role_pb.CreateElevationRequestReq) (*common_pb.SuccessResp, error) {
	rl := sysManagement.NewRoleLogic(ctx, rs.svcCtx)
	return rl.CreateElevationRequest(in)
}

func (rs *RoleServer) DecideElevation(ctx context.Context, in *role_pb.DecideElevationReq) (*common_pb.SuccessResp, error) {
	rl := sysManagement.NewRoleLogic(ctx, rs.svcCtx)
	return rl.DecideElevation(in)
}

func (rs *RoleServer) RevokeElevation(ctx context.Context, in *role_pb.RevokeElevationReq) (*common_pb.SuccessResp, error) {
	rl := sysManagement.NewRoleLogic(ctx, rs.svcCtx)
	return rl.RevokeElevation(in)
}

func (rs *RoleServer) ListElevations(ctx context.Context, in *common_pb.Empty) (*role_pb.ListElevationResp, error) {
	rl := sysManagement.NewRoleLogic(ctx, rs.svcCtx)
	return rl.ListElevations(in)
}

func (rs *RoleServer) GetMyElevations(ctx context.Context, in *common_pb.IdReq) (*role_pb.ListElevationResp, error) {
	rl := sysManagement.NewRoleLogic(ctx, rs.svcCtx)
	return rl.GetMyElevations(in)
}
