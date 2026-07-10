package sysManagement

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/internal/util"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/role_pb"
	"td27/rpc/basis/types/sysManagement/user_pb"
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

func (ul *UserLogic) mapUserToResp(user *sysManagement.UserModel) *user_pb.UserResp {
	if user == nil {
		return nil
	}

	roles := make([]*role_pb.RoleResp, 0, len(user.Roles))
	for _, role := range user.Roles {
		var roleParentID uint64
		if role.ParentID != nil {
			roleParentID = uint64(*role.ParentID)
		}
		roles = append(roles, &role_pb.RoleResp{
			Id:        uint64(role.ID),
			RoleName:  role.RoleName,
			ParentId:  &roleParentID,
			CreatedAt: util.ToProtoTimestamp(role.CreatedAt),
			UpdatedAt: util.ToProtoTimestamp(role.UpdatedAt),
		})
	}

	return &user_pb.UserResp{
		Id:        uint64(user.ID),
		Username:  user.Username,
		Phone:     user.Phone,
		Email:     user.Email,
		Active:    user.Active,
		DeptId:    uint64(user.DeptID),
		Roles:     roles,
		CreatedAt: util.ToProtoTimestamp(user.CreatedAt),
		UpdatedAt: util.ToProtoTimestamp(user.UpdatedAt),
	}
}

func (ul *UserLogic) Login(in *user_pb.LoginReq) (*user_pb.LoginResp, error) {
	user, err := ul.svcCtx.UserService.GetByUsername(ul.ctx, in.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login failed: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Internal, "username doesn't exist")
	}

	if !ul.svcCtx.UserService.VerifyPassword(user.Password, in.Password) {
		return nil, status.Errorf(codes.Internal, "invalid password")
	}

	// Load user roles for JWT claims
	userWithRoles, err := ul.svcCtx.UserService.GetByIDWithRoles(ul.ctx, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login failed: %v", err)
	}

	roleIds := make([]uint64, 0, len(userWithRoles.Roles))
	for _, role := range userWithRoles.Roles {
		roleIds = append(roleIds, uint64(role.ID))
	}

	token, expireAt, err := ul.svcCtx.JWT.CreateToken(uint64(user.ID), user.Username, roleIds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &user_pb.LoginResp{
		Token:     token,
		User:      ul.mapUserToResp(userWithRoles),
		ExpiresAt: util.ToProtoTimestamp(expireAt),
	}, nil
}

func (ul *UserLogic) GetUserInfo(in *common_pb.IdReq) (*user_pb.UserResp, error) {
	user, err := ul.svcCtx.UserService.GetByIDWithRoles(ul.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetUserInfo failed: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return ul.mapUserToResp(user), nil
}

func (ul *UserLogic) ListUser(in *common_pb.PageReq) (*user_pb.ListUserResp, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}
	page := &common.PageInfo{
		Page:     int(in.Page),
		PageSize: int(in.PageSize),
	}

	users, count, err := ul.svcCtx.UserService.List(ul.ctx, page, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ListUser failed: %v", err)
	}

	resp := &user_pb.ListUserResp{
		List:  make([]*user_pb.UserResp, 0, len(users)),
		Total: count,
	}

	for _, user := range users {
		resp.List = append(resp.List, ul.mapUserToResp(user))
	}

	return resp, nil
}

func (ul *UserLogic) DeleteUser(in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id")
	}

	err := ul.svcCtx.UserService.Delete(ul.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete user failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (ul *UserLogic) CreateUser(in *user_pb.CreateUserReq) (*common_pb.SuccessResp, error) {
	user := &sysManagement.UserModel{
		Username: in.Username,
		Phone:    *in.Phone,
		Email:    *in.Email,
		Active:   *in.Active,
		DeptID:   uint(*in.DeptId),
	}

	err := ul.svcCtx.UserService.Create(ul.ctx, user, in.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create user failed: %v", err)
	}

	if len(in.RoleIds) > 0 {
		roleIDs := make([]uint, 0, len(in.RoleIds))
		for _, rid := range in.RoleIds {
			roleIDs = append(roleIDs, uint(rid))
		}
		err = ul.svcCtx.UserService.AssignRoles(ul.ctx, uint(user.ID), roleIDs)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "assign roles failed: %v", err)
		}
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (ul *UserLogic) UpdateUser(in *user_pb.UpdateUserReq) (*user_pb.UserResp, error) {
	user, err := ul.svcCtx.UserService.GetByID(ul.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get user failed: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	if in.Username != nil {
		user.Username = *in.Username
	}
	if in.Phone != nil {
		user.Phone = *in.Phone
	}
	if in.Email != nil {
		user.Email = *in.Email
	}
	if in.Active != nil {
		user.Active = *in.Active
	}
	if in.DeptId != nil {
		user.DeptID = uint(*in.DeptId)
	}

	err = ul.svcCtx.UserService.Update(ul.ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update user failed: %v", err)
	}

	if in.RoleIds != nil && len(in.RoleIds) > 0 {
		roleIDs := make([]uint, 0, len(in.RoleIds))
		for _, rid := range in.RoleIds {
			roleIDs = append(roleIDs, uint(rid))
		}
		err = ul.svcCtx.UserService.AssignRoles(ul.ctx, uint(user.ID), roleIDs)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "update roles failed: %v", err)
		}
	}

	updatedUser, err := ul.svcCtx.UserService.GetByIDWithRoles(ul.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get updated user failed: %v", err)
	}

	return ul.mapUserToResp(updatedUser), nil
}

func (ul *UserLogic) ModifyPassword(in *user_pb.ModifyPasswdReq) (*common_pb.SuccessResp, error) {
	err := ul.svcCtx.UserService.ChangePassword(ul.ctx, uint(in.Id), in.OldPassword, in.NewPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "modify password failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (ul *UserLogic) SwitchUserActive(in *user_pb.SwitchActiveReq) (*common_pb.SuccessResp, error) {
	err := ul.svcCtx.UserService.ToggleActive(ul.ctx, uint(in.Id), in.Active)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "switch active failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (ul *UserLogic) AssignRoles(in *user_pb.AssignRolesReq) (*common_pb.SuccessResp, error) {
	roleIDs := make([]uint, 0, len(in.RoleIds))
	for _, rid := range in.RoleIds {
		roleIDs = append(roleIDs, uint(rid))
	}

	err := ul.svcCtx.UserService.AssignRoles(ul.ctx, uint(in.UserId), roleIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "assign roles failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}
