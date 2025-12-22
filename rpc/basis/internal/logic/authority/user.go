package authority

import (
	"basis/types/authority/user_pb"
	"context"
	"database/sql"
	"errors"

	modelAuthority "basis/internal/model/authority"
	"basis/internal/pkg"
	"basis/internal/svc"
	"basis/types/common_pb"

	"github.com/zeromicro/go-zero/core/logx"
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

	return &common_pb.SuccessResp{
		Success: true,
	}, nil
}

func (ul *UserLogic) CreateUser(in *user_pb.CreateUserReq) (*common_pb.SuccessResp, error) {
	// check role exists
	//exists, err := ul.svcCtx.AuthorityUserRepo.RoleExists(l.ctx, uint(in.RoleId))
	//if err != nil {
	//	return nil, err
	//}
	//if !exists {
	//	return nil, errors.New("角色不存在")
	//}

	// create user model
	_, err := ul.svcCtx.AuthorityUserRepo.Insert(ul.ctx, &modelAuthority.AuthorityUser{
		Username: sql.NullString{
			String: in.Username,
			Valid:  in.Username != "",
		},
		Password: pkg.MD5V([]byte(in.Password)),
		Phone: sql.NullString{
			String: in.Phone,
			Valid:  in.Phone != "",
		},
		Email: sql.NullString{
			String: in.Email,
			Valid:  in.Email != "",
		},
		Active: sql.NullInt64{
			Int64: func() int64 {
				if in.Active {
					return 1
				}
				return 0
			}(),
			Valid: true,
		},
		RoleModelId: in.RoleId,
	})
	if err != nil {
		return nil, err
	}

	return &common_pb.SuccessResp{Success: true}, nil
}
