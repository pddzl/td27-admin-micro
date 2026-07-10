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

func (rl *RoleLogic) mapRoleToResp(role *sysManagement.RoleModel) *role_pb.RoleResp {
	if role == nil {
		return nil
	}

	var parentID uint64
	if role.ParentID != nil {
		parentID = uint64(*role.ParentID)
	}

	return &role_pb.RoleResp{
		Id:        uint64(role.ID),
		RoleName:  role.RoleName,
		ParentId:  &parentID,
		CreatedAt: util.ToProtoTimestamp(role.CreatedAt),
		UpdatedAt: util.ToProtoTimestamp(role.UpdatedAt),
	}
}

func (rl *RoleLogic) mapRoleToTree(role *sysManagement.RoleModel, allRoles []*sysManagement.RoleModel) *role_pb.RoleTreeResp {
	resp := &role_pb.RoleTreeResp{
		Role:     rl.mapRoleToResp(role),
		Children: make([]*role_pb.RoleTreeResp, 0),
	}

	for _, child := range allRoles {
		if child.ParentID != nil && *child.ParentID == role.ID {
			resp.Children = append(resp.Children, rl.mapRoleToTree(child, allRoles))
		}
	}

	return resp
}

func (rl *RoleLogic) GetRole(in *common_pb.IdReq) (*role_pb.RoleResp, error) {
	role, err := rl.svcCtx.RoleService.GetByID(rl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get role failed: %v", err)
	}
	if role == nil {
		return nil, status.Errorf(codes.NotFound, "role not found")
	}

	return rl.mapRoleToResp(role), nil
}

func (rl *RoleLogic) GetRoleWithChildren(in *common_pb.IdReq) (*role_pb.RoleTreeResp, error) {
	role, err := rl.svcCtx.RoleService.GetByID(rl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get role failed: %v", err)
	}
	if role == nil {
		return nil, status.Errorf(codes.NotFound, "role not found")
	}

	allRoles, err := rl.svcCtx.RoleService.GetAll(rl.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get all roles failed: %v", err)
	}

	return rl.mapRoleToTree(role, allRoles), nil
}

func (rl *RoleLogic) ListRole(in *common_pb.PageReq) (*role_pb.ListRoleResp, error) {
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

	roles, count, err := rl.svcCtx.RoleService.List(rl.ctx, page)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list roles failed: %v", err)
	}

	resp := &role_pb.ListRoleResp{
		List:  make([]*role_pb.RoleResp, 0, len(roles)),
		Total: count,
	}

	for _, role := range roles {
		resp.List = append(resp.List, rl.mapRoleToResp(role))
	}

	return resp, nil
}

func (rl *RoleLogic) GetAllRoles(in *common_pb.Empty) (*role_pb.ListRoleResp, error) {
	roles, err := rl.svcCtx.RoleService.GetAll(rl.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get all roles failed: %v", err)
	}

	resp := &role_pb.ListRoleResp{
		List:  make([]*role_pb.RoleResp, 0, len(roles)),
		Total: int64(len(roles)),
	}

	for _, role := range roles {
		resp.List = append(resp.List, rl.mapRoleToResp(role))
	}

	return resp, nil
}

func (rl *RoleLogic) CreateRole(in *role_pb.CreateRoleReq) (*common_pb.SuccessResp, error) {
	parentID := uint(*in.ParentId)
	role := &sysManagement.RoleModel{
		RoleName: in.RoleName,
		ParentID: &parentID,
	}

	err := rl.svcCtx.RoleService.Create(rl.ctx, role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create role failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (rl *RoleLogic) UpdateRole(in *role_pb.UpdateRoleReq) (*role_pb.RoleResp, error) {
	role, err := rl.svcCtx.RoleService.GetByID(rl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get role failed: %v", err)
	}
	if role == nil {
		return nil, status.Errorf(codes.NotFound, "role not found")
	}

	if in.RoleName != nil {
		role.RoleName = *in.RoleName
	}
	if in.ParentId != nil {
		parentID := uint(*in.ParentId)
		role.ParentID = &parentID
	}

	err = rl.svcCtx.RoleService.Update(rl.ctx, role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update role failed: %v", err)
	}

	updatedRole, err := rl.svcCtx.RoleService.GetByID(rl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get updated role failed: %v", err)
	}

	return rl.mapRoleToResp(updatedRole), nil
}

func (rl *RoleLogic) DeleteRole(in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid role id")
	}

	err := rl.svcCtx.RoleService.Delete(rl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete role failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (rl *RoleLogic) AssignPermissions(in *role_pb.AssignPermissionsReq) (*common_pb.SuccessResp, error) {
	permIDs := make([]uint, 0, len(in.PermissionIds))
	for _, pid := range in.PermissionIds {
		permIDs = append(permIDs, uint(pid))
	}

	err := rl.svcCtx.RoleService.AssignPermissions(rl.ctx, uint(in.RoleId), permIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "assign permissions failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (rl *RoleLogic) GetRolePermissions(in *common_pb.IdReq) (*role_pb.GetRolePermissionsResp, error) {
	permIDs, err := rl.svcCtx.RoleService.GetPermissions(rl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get role permissions failed: %v", err)
	}

	resp := &role_pb.GetRolePermissionsResp{
		PermissionIds: make([]uint64, 0, len(permIDs)),
	}

	for _, pid := range permIDs {
		resp.PermissionIds = append(resp.PermissionIds, uint64(pid))
	}

	return resp, nil
}
