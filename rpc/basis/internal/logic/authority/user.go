package authority

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"basis/internal/model/authority"
	"basis/internal/pkg"
	"basis/internal/svc"
	"basis/types/authority/user_pb"
	"basis/types/common_pb"
)

type UserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (ul *UserLogic) DeleteUser(in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 {
		return nil, errors.New("invalid user id")
	}

	err := ul.svcCtx.AuthorityUserRepo.Delete(ul.ctx, in.Id)
	if err != nil {
		ul.Errorf("delete user failed, id=%d, err=%v", in.Id, err)
		return nil, err
	}

	ul.Infof("user deleted, id=%d", in.Id)

	return &common_pb.SuccessResp{Success: true}, nil
}

func (ul *UserLogic) CreateUser(in *user_pb.CreateUserReq) (*common_pb.SuccessResp, error) {
	// check role exists
	exists, err := ul.svcCtx.AuthorityUserRepo.ExistsById(ul.ctx, uint(in.RoleId))
	if err != nil {
		return nil, err
	}
	if !exists {
		ul.Logger.Errorf("role not found")
		return nil, errors.New("role not found")
	}

	// build entity
	user := &authority.AuthorityUserEntity{
		Username:    in.Username,
		Password:    pkg.MD5V([]byte(in.Password)),
		Phone:       in.Phone,
		Email:       in.Email,
		Active:      in.Active,
		RoleModelID: uint(in.RoleId),
	}

	err = ul.svcCtx.AuthorityUserRepo.Insert(ul.ctx, user)
	if err != nil {
		ul.Logger.Errorf("insert user entity failed, err=%v", err)
		return nil, err
	}

	return &common_pb.SuccessResp{Success: true}, nil
}
