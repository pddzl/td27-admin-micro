package authority

import (
	"context"

	"td27/rpc/basis/internal/logic/authority"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/authority/user_pb"
	"td27/rpc/basis/types/common_pb"
)

type UserServer struct {
	svcCtx *svc.ServiceContext
	user_pb.UnimplementedUserServer
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

func (us *UserServer) DeleteUser(ctx context.Context, in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	ul := authority.NewUserLogic(ctx, us.svcCtx)
	return ul.DeleteUser(in)
}

func (us *UserServer) GetUserInfo(ctx context.Context, in *common_pb.IdReq) (*user_pb.UserResp, error) {
	ul := authority.NewUserLogic(ctx, us.svcCtx)
	return ul.GetUserInfo(in)
}

func (us *UserServer) ListUser(ctx context.Context, in *common_pb.PageReq) (*user_pb.ListUserResp, error) {
	ul := authority.NewUserLogic(ctx, us.svcCtx)
	return ul.ListUser(in)
}

func (us *UserServer) CreateUser(ctx context.Context, in *user_pb.CreateUserReq) (*common_pb.SuccessResp, error) {
	ul := authority.NewUserLogic(ctx, us.svcCtx)
	return ul.CreateUser(in)
}

func (us *UserServer) UpdateUser(ctx context.Context, in *user_pb.UpdateUserReq) (*user_pb.UserResp, error) {
	ul := authority.NewUserLogic(ctx, us.svcCtx)
	return ul.UpdateUser(in)
}

func (us *UserServer) ModifyPassword(ctx context.Context, in *user_pb.ModifyPasswdReq) (*common_pb.SuccessResp, error) {
	ul := authority.NewUserLogic(ctx, us.svcCtx)
	return ul.ModifyPassword(in)
}

func (us *UserServer) SwitchUserActive(ctx context.Context, in *user_pb.SwitchActiveReq) (*common_pb.SuccessResp, error) {
	ul := authority.NewUserLogic(ctx, us.svcCtx)
	return ul.SwitchUserActive(in)
}

func (us *UserServer) AssignRoles(ctx context.Context, in *user_pb.AssignRolesReq) (*common_pb.SuccessResp, error) {
	ul := authority.NewUserLogic(ctx, us.svcCtx)
	return ul.AssignRoles(in)
}

func (us *UserServer) Login(ctx context.Context, in *user_pb.LoginReq) (*user_pb.LoginResp, error) {
	ul := authority.NewUserLogic(ctx, us.svcCtx)
	return ul.Login(in)
}
