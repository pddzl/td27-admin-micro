package authority

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"td27/pkg/tool"
	"td27/rpc/basis/internal/model/authority"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/authority/user_pb"
	"td27/rpc/basis/types/common_pb"
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

func (ul *UserLogic) mapUserEntityToUserResp(entity *authority.AuthorityUserEntity) *user_pb.UserResp {
	if entity == nil {
		return nil
	}

	return &user_pb.UserResp{
		Id:        entity.ID,
		Username:  entity.Username,
		Phone:     entity.Phone,
		Email:     entity.Email,
		Active:    entity.Active,
		RoleId:    entity.RoleModelID,
		CreatedAt: tool.ToProtoTimestamp(entity.CreatedAt),
		UpdatedAt: tool.ToProtoTimestamp(entity.UpdatedAt),
	}
}

func (ul *UserLogic) mapUserPlusRoleNameDTOToUserRoleResp(dto *authority.UserPlusRoleNameDTO) *user_pb.UserRoleResp {
	if dto == nil {
		return nil
	}

	userResp := ul.mapUserEntityToUserResp(&dto.AuthorityUserEntity)

	return &user_pb.UserRoleResp{
		User:     userResp,
		RoleName: dto.RoleName,
	}
}

func (ul *UserLogic) GetUserInfo(in *common_pb.IdReq) (*user_pb.UserRoleResp, error) {
	userPlusRoleNameDTO, err := ul.svcCtx.AuthorityUserRepo.FindOne(ul.ctx, in.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetUserInfo failed: %v", err)
	}

	return ul.mapUserPlusRoleNameDTOToUserRoleResp(userPlusRoleNameDTO), nil
}

func (ul *UserLogic) ListUser(in *common_pb.PageReq) (*user_pb.ListUserResp, error) {
	list, count, err := ul.svcCtx.AuthorityUserRepo.List(ul.ctx, int(in.Page), int(in.PageSize))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ListUser failed: %v", err)
	}

	resp := &user_pb.ListUserResp{
		List:  make([]*user_pb.UserRoleResp, 0, len(list)),
		Total: count,
	}

	for _, user := range list {
		resp.List = append(
			resp.List,
			ul.mapUserPlusRoleNameDTOToUserRoleResp(&user),
		)
	}

	return resp, nil
}

func (ul *UserLogic) DeleteUser(in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 {
		return nil, status.Errorf(codes.Internal, "invalid user id")
	}

	err := ul.svcCtx.AuthorityUserRepo.Delete(ul.ctx, in.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete user failed, id=%d, err=%v", in.Id, err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (ul *UserLogic) CreateUser(in *user_pb.CreateUserReq) (*common_pb.SuccessResp, error) {
	// check role exists
	exists, err := ul.svcCtx.AuthorityRoleRepo.ExistsById(ul.ctx, in.RoleId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !exists {
		return nil, status.Errorf(codes.Internal, "role not found")
	}

	// build entity
	user := &authority.AuthorityUserEntity{
		Username:    in.Username,
		Password:    tool.MD5V([]byte(in.Password)),
		Phone:       *in.Phone,
		Email:       *in.Email,
		Active:      *in.Active,
		RoleModelID: in.RoleId,
	}

	// create entity
	err = ul.svcCtx.AuthorityUserRepo.Insert(ul.ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "insert user entity failed, err=%v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (ul *UserLogic) UpdateUser(in *user_pb.UpdateUserReq) (*user_pb.UserResp, error) {
	// check user exists
	userExists, err := ul.svcCtx.AuthorityUserRepo.ExistsById(ul.ctx, in.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if !userExists {
		return nil, status.Errorf(codes.Internal, "user not found")
	}

	// check role exists
	roleExists, err := ul.svcCtx.AuthorityRoleRepo.ExistsById(ul.ctx, in.RoleId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if !roleExists {
		return nil, status.Errorf(codes.Internal, "role not found")
	}

	// update user
	userEntity, err := ul.svcCtx.AuthorityUserRepo.Update(
		ul.ctx,
		&authority.UpdateUserDTO{
			ID:          in.Id,
			Username:    in.Username,
			Password:    in.Password,
			Phone:       in.Phone,
			Email:       in.Email,
			Active:      in.Active,
			RoleModelId: in.RoleId,
		},
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "update user failed: %v", err)
	}

	// entity → proto
	return &user_pb.UserResp{
		Id:        userEntity.ID,
		Username:  userEntity.Username,
		Phone:     userEntity.Phone,
		Email:     userEntity.Email,
		Active:    userEntity.Active,
		RoleId:    userEntity.RoleModelID,
		CreatedAt: tool.ToProtoTimestamp(userEntity.CreatedAt),
		UpdatedAt: tool.ToProtoTimestamp(userEntity.UpdatedAt),
	}, nil
}

func (ul *UserLogic) ModifyPassword(in *user_pb.ModifyPasswdReq) (*common_pb.SuccessResp, error) {
	err := ul.svcCtx.AuthorityUserRepo.ModifyPassword(ul.ctx, in.Id, in.OldPassword, in.NewPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "modify password failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (ul *UserLogic) SwitchUserActive(in *user_pb.SwitchActiveReq) (*common_pb.SuccessResp, error) {
	err := ul.svcCtx.AuthorityUserRepo.SwitchUserActive(ul.ctx, in.Id, in.Active)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "switch active failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}
