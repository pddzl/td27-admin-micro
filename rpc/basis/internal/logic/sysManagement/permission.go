package sysManagement

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/permission_pb"
"td27/rpc/basis/internal/util"
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

func (pl *PermissionLogic) mapPermissionToResp(perm *sysManagement.PermissionModel) *permission_pb.PermissionResp {
	if perm == nil {
		return nil
	}

	return &permission_pb.PermissionResp{
		Id:        int64(perm.ID),
		Name:      perm.Name,
		Domain:    permissionDomainToProto(perm.Domain),
		Resource:  perm.Resource,
		Action:    actionToProto(perm.Action),
		Effect:    effectToProto(perm.Effect),
		DomainId:  int64(perm.DomainID),
		CreatedAt: util.Ts(perm.CreatedAt),
		UpdatedAt: util.Ts(perm.UpdatedAt),
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

	perms, countt, err := pl.svcCtx.PermService.List(pl.ctx, page, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list permissions failed: %v", err)
	}

	resp := &permission_pb.ListPermissionResp{
		List:  make([]*permission_pb.PermissionResp, 0, len(perms)),
		Total: countt,
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
	perm := &sysManagement.PermissionModel{
		Name:     in.Name,
		Domain:   permissionDomainFromProto(in.Domain),
		Resource: in.Resource,
		Action:   actionFromProto(in.Action),
		Effect:   effectFromProto(in.Effect),
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
	if in.Effect != nil {
		perm.Effect = effectFromProto(*in.Effect)
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
	var roleIDs []uint
	if in.UserId > 0 {
		// Live role evaluation: roles resolved from the database override the
		// (possibly stale) role ids embedded in the caller's JWT.
		liveRoleIDs, err := pl.svcCtx.PermService.GetRolesForUser(pl.ctx, uint(in.UserId))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "check permission failed: %v", err)
		}

		// JIT elevations: approved, unexpired temporary role grants are merged
		// in, so an approved elevation authorizes on the next request and an
		// expired one stops authorizing immediately.
		elevatedRoleIDs, err := pl.svcCtx.ElevService.ActiveRoleIDs(pl.ctx, uint(in.UserId))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "check permission failed: %v", err)
		}

		if len(elevatedRoleIDs) > 0 {
			seen := make(map[uint]struct{}, len(liveRoleIDs)+len(elevatedRoleIDs))
			for _, id := range liveRoleIDs {
				seen[id] = struct{}{}
			}
			for _, id := range elevatedRoleIDs {
				if _, ok := seen[id]; !ok {
					liveRoleIDs = append(liveRoleIDs, id)
				}
			}
		}
		roleIDs = liveRoleIDs
	} else {
		// Legacy path: caller supplied role ids directly.
		roleIDs = make([]uint, 0, len(in.RoleIds))
		for _, rid := range in.RoleIds {
			roleIDs = append(roleIDs, uint(rid))
		}
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

func (pl *PermissionLogic) LintPolicies(in *common_pb.Empty) (*permission_pb.LintPoliciesResp, error) {
	issues, err := pl.svcCtx.PermService.LintPolicies(pl.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lint policies failed: %v", err)
	}

	resp := &permission_pb.LintPoliciesResp{
		Issues: make([]*permission_pb.PolicyIssue, 0, len(issues)),
	}
	for _, issue := range issues {
		resp.Issues = append(resp.Issues, &permission_pb.PolicyIssue{
			Severity:     issue.Severity,
			Message:      issue.Message,
			RoleId:       int64(issue.RoleID),
			PermissionId: int64(issue.PermissionID),
		})
	}
	resp.Total = int64(len(issues))

	return resp, nil
}

func permissionDomainToProto(domain sysManagement.PermissionDomain) permission_pb.PermissionDomain {
	switch domain {
	case sysManagement.PermissionDomainMenu:
		return permission_pb.PermissionDomain_DOMAIN_MENU
	case sysManagement.PermissionDomainAPI:
		return permission_pb.PermissionDomain_DOMAIN_API
	case sysManagement.PermissionDomainButton:
		return permission_pb.PermissionDomain_DOMAIN_BUTTON
	case sysManagement.PermissionDomainData:
		return permission_pb.PermissionDomain_DOMAIN_DATA
	default:
		return permission_pb.PermissionDomain_DOMAIN_MENU
	}
}

func permissionDomainFromProto(domain permission_pb.PermissionDomain) sysManagement.PermissionDomain {
	switch domain {
	case permission_pb.PermissionDomain_DOMAIN_MENU:
		return sysManagement.PermissionDomainMenu
	case permission_pb.PermissionDomain_DOMAIN_API:
		return sysManagement.PermissionDomainAPI
	case permission_pb.PermissionDomain_DOMAIN_BUTTON:
		return sysManagement.PermissionDomainButton
	case permission_pb.PermissionDomain_DOMAIN_DATA:
		return sysManagement.PermissionDomainData
	default:
		return sysManagement.PermissionDomainMenu
	}
}

func actionToProto(action sysManagement.Action) permission_pb.Action {
	switch action {
	case sysManagement.ActionAll:
		return permission_pb.Action_ACTION_ALL
	case sysManagement.ActionView:
		return permission_pb.Action_ACTION_VIEW
	case sysManagement.ActionRead:
		return permission_pb.Action_ACTION_READ
	case sysManagement.ActionCreate:
		return permission_pb.Action_ACTION_CREATE
	case sysManagement.ActionUpdate:
		return permission_pb.Action_ACTION_UPDATE
	case sysManagement.ActionDelete:
		return permission_pb.Action_ACTION_DELETE
	case sysManagement.ActionExecute:
		return permission_pb.Action_ACTION_EXECUTE
	default:
		return permission_pb.Action_ACTION_ALL
	}
}

func actionFromProto(action permission_pb.Action) sysManagement.Action {
	switch action {
	case permission_pb.Action_ACTION_ALL:
		return sysManagement.ActionAll
	case permission_pb.Action_ACTION_VIEW:
		return sysManagement.ActionView
	case permission_pb.Action_ACTION_READ:
		return sysManagement.ActionRead
	case permission_pb.Action_ACTION_CREATE:
		return sysManagement.ActionCreate
	case permission_pb.Action_ACTION_UPDATE:
		return sysManagement.ActionUpdate
	case permission_pb.Action_ACTION_DELETE:
		return sysManagement.ActionDelete
	case permission_pb.Action_ACTION_EXECUTE:
		return sysManagement.ActionExecute
	default:
		return sysManagement.ActionAll
	}
}

func effectToProto(effect sysManagement.Effect) permission_pb.PermissionEffect {
	if effect == sysManagement.EffectDeny {
		return permission_pb.PermissionEffect_EFFECT_DENY
	}
	return permission_pb.PermissionEffect_EFFECT_ALLOW
}

func effectFromProto(effect permission_pb.PermissionEffect) sysManagement.Effect {
	if effect == permission_pb.PermissionEffect_EFFECT_DENY {
		return sysManagement.EffectDeny
	}
	return sysManagement.EffectAllow
}
