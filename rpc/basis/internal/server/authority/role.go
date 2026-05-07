package authority

import (
	"context"

	"td27/rpc/basis/internal/logic/authority"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/authority/role_pb"
	"td27/rpc/basis/types/common_pb"
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
	rl := authority.NewRoleLogic(ctx, rs.svcCtx)
	return rl.GetRole(in)
}

func (rs *RoleServer) GetRoleWithChildren(ctx context.Context, in *common_pb.IdReq) (*role_pb.RoleTreeResp, error) {
	rl := authority.NewRoleLogic(ctx, rs.svcCtx)
	return rl.GetRoleWithChildren(in)
}

func (rs *RoleServer) ListRole(ctx context.Context, in *common_pb.PageReq) (*role_pb.ListRoleResp, error) {
	rl := authority.NewRoleLogic(ctx, rs.svcCtx)
	return rl.ListRole(in)
}

func (rs *RoleServer) GetAllRoles(ctx context.Context, in *common_pb.Empty) (*role_pb.ListRoleResp, error) {
	rl := authority.NewRoleLogic(ctx, rs.svcCtx)
	return rl.GetAllRoles(in)
}

func (rs *RoleServer) CreateRole(ctx context.Context, in *role_pb.CreateRoleReq) (*common_pb.SuccessResp, error) {
	rl := authority.NewRoleLogic(ctx, rs.svcCtx)
	return rl.CreateRole(in)
}

func (rs *RoleServer) UpdateRole(ctx context.Context, in *role_pb.UpdateRoleReq) (*role_pb.RoleResp, error) {
	rl := authority.NewRoleLogic(ctx, rs.svcCtx)
	return rl.UpdateRole(in)
}

func (rs *RoleServer) DeleteRole(ctx context.Context, in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	rl := authority.NewRoleLogic(ctx, rs.svcCtx)
	return rl.DeleteRole(in)
}

func (rs *RoleServer) AssignPermissions(ctx context.Context, in *role_pb.AssignPermissionsReq) (*common_pb.SuccessResp, error) {
	rl := authority.NewRoleLogic(ctx, rs.svcCtx)
	return rl.AssignPermissions(in)
}

func (rs *RoleServer) GetRolePermissions(ctx context.Context, in *common_pb.IdReq) (*role_pb.GetRolePermissionsResp, error) {
	rl := authority.NewRoleLogic(ctx, rs.svcCtx)
	return rl.GetRolePermissions(in)
}
