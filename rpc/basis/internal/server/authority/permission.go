package authority

import (
	"context"

	"td27/rpc/basis/internal/logic/authority"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/authority/permission_pb"
	"td27/rpc/basis/types/common_pb"
)

type PermissionServer struct {
	svcCtx *svc.ServiceContext
	permission_pb.UnimplementedPermissionServer
}

func NewPermissionServer(svcCtx *svc.ServiceContext) *PermissionServer {
	return &PermissionServer{
		svcCtx: svcCtx,
	}
}

func (ps *PermissionServer) GetPermission(ctx context.Context, in *common_pb.IdReq) (*permission_pb.PermissionResp, error) {
	pl := authority.NewPermissionLogic(ctx, ps.svcCtx)
	return pl.GetPermission(in)
}

func (ps *PermissionServer) ListPermission(ctx context.Context, in *common_pb.PageReq) (*permission_pb.ListPermissionResp, error) {
	pl := authority.NewPermissionLogic(ctx, ps.svcCtx)
	return pl.ListPermission(in)
}

func (ps *PermissionServer) GetAllPermissions(ctx context.Context, in *common_pb.Empty) (*permission_pb.ListPermissionResp, error) {
	pl := authority.NewPermissionLogic(ctx, ps.svcCtx)
	return pl.GetAllPermissions(in)
}

func (ps *PermissionServer) GetPermissionsByRoleId(ctx context.Context, in *common_pb.IdReq) (*permission_pb.ListPermissionResp, error) {
	pl := authority.NewPermissionLogic(ctx, ps.svcCtx)
	return pl.GetPermissionsByRoleId(in)
}

func (ps *PermissionServer) CreatePermission(ctx context.Context, in *permission_pb.CreatePermissionReq) (*common_pb.SuccessResp, error) {
	pl := authority.NewPermissionLogic(ctx, ps.svcCtx)
	return pl.CreatePermission(in)
}

func (ps *PermissionServer) UpdatePermission(ctx context.Context, in *permission_pb.UpdatePermissionReq) (*permission_pb.PermissionResp, error) {
	pl := authority.NewPermissionLogic(ctx, ps.svcCtx)
	return pl.UpdatePermission(in)
}

func (ps *PermissionServer) DeletePermission(ctx context.Context, in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	pl := authority.NewPermissionLogic(ctx, ps.svcCtx)
	return pl.DeletePermission(in)
}

func (ps *PermissionServer) CheckPermission(ctx context.Context, in *permission_pb.CheckPermissionReq) (*permission_pb.CheckPermissionResp, error) {
	pl := authority.NewPermissionLogic(ctx, ps.svcCtx)
	return pl.CheckPermission(in)
}

func (ps *PermissionServer) ReloadPolicy(ctx context.Context, in *common_pb.Empty) (*common_pb.SuccessResp, error) {
	pl := authority.NewPermissionLogic(ctx, ps.svcCtx)
	return pl.ReloadPolicy(in)
}
