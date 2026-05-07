package authority

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"td27/rpc/basis/internal/model/authority"
	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/internal/util"
	"td27/rpc/basis/types/authority/permission_pb"
	"td27/rpc/basis/types/common_pb"
)

type PermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionLogic {
	return &PermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (pl *PermissionLogic) mapPermissionToResp(perm *authority.PermissionModel) *permission_pb.PermissionResp {
	if perm == nil {
		return nil
	}

	return &permission_pb.PermissionResp{
		Id:         uint64(perm.ID),
		Name:       perm.Name,
		Domain:     permissionDomainToProto(perm.Domain),
		Resource:   perm.Resource,
		Action:     actionToProto(perm.Action),
		DomainId:   uint64(perm.DomainID),
		CreatedAt:  util.ToProtoTimestamp(perm.CreatedAt),
		UpdatedAt:  util.ToProtoTimestamp(perm.UpdatedAt),
	}
}

func (pl *PermissionLogic) GetPermission(in *common_pb.IdReq) (*permission_pb.PermissionResp, error) {
	perm, err := pl.svcCtx.PermService.GetByID(pl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get permission failed: %v", err)
	}
	if perm == nil {
		return nil, status.Errorf(codes.NotFound, "permission not found")
	}

	return pl.mapPermissionToResp(perm), nil
}

func (pl *PermissionLogic) ListPermission(in *common_pb.PageReq) (*permission_pb.ListPermissionResp, error) {
	page := &common.PageInfo{
		Page:     int(in.Page),
		PageSize: int(in.PageSize),
	}

	perms, count, err := pl.svcCtx.PermService.List(pl.ctx, page, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list permissions failed: %v", err)
	}

	resp := &permission_pb.ListPermissionResp{
		List:  make([]*permission_pb.PermissionResp, 0, len(perms)),
		Total: count,
	}

	for _, perm := range perms {
		resp.List = append(resp.List, pl.mapPermissionToResp(perm))
	}

	return resp, nil
}

func (pl *PermissionLogic) GetAllPermissions(in *common_pb.Empty) (*permission_pb.ListPermissionResp, error) {
	perms, err := pl.svcCtx.PermService.GetAll(pl.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get all permissions failed: %v", err)
	}

	resp := &permission_pb.ListPermissionResp{
		List:  make([]*permission_pb.PermissionResp, 0, len(perms)),
		Total: int64(len(perms)),
	}

	for _, perm := range perms {
		resp.List = append(resp.List, pl.mapPermissionToResp(perm))
	}

	return resp, nil
}

func (pl *PermissionLogic) GetPermissionsByRoleId(in *common_pb.IdReq) (*permission_pb.ListPermissionResp, error) {
	perms, err := pl.svcCtx.PermService.GetByRoleID(pl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get permissions by role id failed: %v", err)
	}

	resp := &permission_pb.ListPermissionResp{
		List:  make([]*permission_pb.PermissionResp, 0, len(perms)),
		Total: int64(len(perms)),
	}

	for _, perm := range perms {
		resp.List = append(resp.List, pl.mapPermissionToResp(perm))
	}

	return resp, nil
}

func (pl *PermissionLogic) CreatePermission(in *permission_pb.CreatePermissionReq) (*common_pb.SuccessResp, error) {
	perm := &authority.PermissionModel{
		Name:     in.Name,
		Domain:   permissionDomainFromProto(in.Domain),
		Resource: in.Resource,
		Action:   actionFromProto(in.Action),
		DomainID: uint(in.DomainId),
	}

	err := pl.svcCtx.PermService.Create(pl.ctx, perm)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create permission failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (pl *PermissionLogic) UpdatePermission(in *permission_pb.UpdatePermissionReq) (*permission_pb.PermissionResp, error) {
	perm, err := pl.svcCtx.PermService.GetByID(pl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get permission failed: %v", err)
	}
	if perm == nil {
		return nil, status.Errorf(codes.NotFound, "permission not found")
	}

	if in.Name != nil {
		perm.Name = *in.Name
	}
	if in.Domain != nil {
		perm.Domain = permissionDomainFromProto(*in.Domain)
	}
	if in.Resource != nil {
		perm.Resource = *in.Resource
	}
	if in.Action != nil {
		perm.Action = actionFromProto(*in.Action)
	}
	if in.DomainId != nil {
		perm.DomainID = uint(*in.DomainId)
	}

	err = pl.svcCtx.PermService.Update(pl.ctx, perm)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update permission failed: %v", err)
	}

	updatedPerm, err := pl.svcCtx.PermService.GetByID(pl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get updated permission failed: %v", err)
	}

	return pl.mapPermissionToResp(updatedPerm), nil
}

func (pl *PermissionLogic) DeletePermission(in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid permission id")
	}

	err := pl.svcCtx.PermService.Delete(pl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete permission failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (pl *PermissionLogic) CheckPermission(in *permission_pb.CheckPermissionReq) (*permission_pb.CheckPermissionResp, error) {
	roleIDs := make([]uint, 0, len(in.RoleIds))
	for _, rid := range in.RoleIds {
		roleIDs = append(roleIDs, uint(rid))
	}

	allowed, err := pl.svcCtx.PermService.CheckPermission(pl.ctx, roleIDs, in.Resource, actionFromProto(in.Action))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "check permission failed: %v", err)
	}

	return &permission_pb.CheckPermissionResp{
		Allowed: allowed,
	}, nil
}

func (pl *PermissionLogic) ReloadPolicy(in *common_pb.Empty) (*common_pb.SuccessResp, error) {
	err := pl.svcCtx.PermService.ReloadPolicy(pl.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "reload policy failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func permissionDomainToProto(domain authority.PermissionDomain) permission_pb.PermissionDomain {
	switch domain {
	case authority.PermissionDomainMenu:
		return permission_pb.PermissionDomain_DOMAIN_MENU
	case authority.PermissionDomainAPI:
		return permission_pb.PermissionDomain_DOMAIN_API
	case authority.PermissionDomainButton:
		return permission_pb.PermissionDomain_DOMAIN_BUTTON
	case authority.PermissionDomainData:
		return permission_pb.PermissionDomain_DOMAIN_DATA
	default:
		return permission_pb.PermissionDomain_DOMAIN_MENU
	}
}

func permissionDomainFromProto(domain permission_pb.PermissionDomain) authority.PermissionDomain {
	switch domain {
	case permission_pb.PermissionDomain_DOMAIN_MENU:
		return authority.PermissionDomainMenu
	case permission_pb.PermissionDomain_DOMAIN_API:
		return authority.PermissionDomainAPI
	case permission_pb.PermissionDomain_DOMAIN_BUTTON:
		return authority.PermissionDomainButton
	case permission_pb.PermissionDomain_DOMAIN_DATA:
		return authority.PermissionDomainData
	default:
		return authority.PermissionDomainMenu
	}
}

func actionToProto(action authority.Action) permission_pb.Action {
	switch action {
	case authority.ActionAll:
		return permission_pb.Action_ACTION_ALL
	case authority.ActionView:
		return permission_pb.Action_ACTION_VIEW
	case authority.ActionRead:
		return permission_pb.Action_ACTION_READ
	case authority.ActionCreate:
		return permission_pb.Action_ACTION_CREATE
	case authority.ActionUpdate:
		return permission_pb.Action_ACTION_UPDATE
	case authority.ActionDelete:
		return permission_pb.Action_ACTION_DELETE
	case authority.ActionExecute:
		return permission_pb.Action_ACTION_EXECUTE
	default:
		return permission_pb.Action_ACTION_ALL
	}
}

func actionFromProto(action permission_pb.Action) authority.Action {
	switch action {
	case permission_pb.Action_ACTION_ALL:
		return authority.ActionAll
	case permission_pb.Action_ACTION_VIEW:
		return authority.ActionView
	case permission_pb.Action_ACTION_READ:
		return authority.ActionRead
	case permission_pb.Action_ACTION_CREATE:
		return authority.ActionCreate
	case permission_pb.Action_ACTION_UPDATE:
		return authority.ActionUpdate
	case permission_pb.Action_ACTION_DELETE:
		return authority.ActionDelete
	case permission_pb.Action_ACTION_EXECUTE:
		return authority.ActionExecute
	default:
		return authority.ActionAll
	}
}
