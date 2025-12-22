package authority

import (
	"context"
	"errors"

	"basis/internal/svc"
	"basis/types/common_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleLogic {
	return &RoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (rl *RoleLogic) FindRoleById(in *common_pb.IdReq) {}

func (rl *RoleLogic) DeleteRole(in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 {
		return nil, errors.New("invalid Role id")
	}

	err := rl.svcCtx.AuthorityRoleRepo.Delete(rl.ctx, in.Id)
	if err != nil {
		rl.Errorf("delete Role failed, id=%d, err=%v", in.Id, err)
		return nil, err
	}

	rl.Infof("Role deleted, id=%d", in.Id)

	return &common_pb.SuccessResp{
		Success: true,
	}, nil
}

func (rl *RoleLogic) CreateRole() (*common_pb.SuccessResp, error) {
	return &common_pb.SuccessResp{Success: true}, nil
}
