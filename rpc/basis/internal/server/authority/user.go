package authority

import (
	"context"

	"basis/internal/logic/authority"
	"basis/internal/svc"
	"basis/types/authority/user_pb"
	"basis/types/common_pb"
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
	nl := authority.NewUserLogic(ctx, us.svcCtx)
	return nl.DeleteUser(in)
}
