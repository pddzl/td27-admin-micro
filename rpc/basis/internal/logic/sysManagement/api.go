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
	"td27/rpc/basis/types/sysManagement/api_pb"
)

type APILogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAPILogic(ctx context.Context, svcCtx *svc.ServiceContext) *APILogic {
	return &APILogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (al *APILogic) mapAPIToResp(api *sysManagement.ApiModel) *api_pb.APIResp {
	if api == nil {
		return nil
	}

	return &api_pb.APIResp{
		Id:          uint64(api.ID),
		Path:        api.Path,
		Method:      api.Method,
		GroupEn:     api.GroupEN,
		GroupCn:     api.GroupCN,
		Description: api.Description,
		CreatedAt:   util.ToProtoTimestamp(api.CreatedAt),
		UpdatedAt:   util.ToProtoTimestamp(api.UpdatedAt),
	}
}

func (al *APILogic) GetAPI(in *common_pb.IdReq) (*api_pb.APIResp, error) {
	api, err := al.svcCtx.APIService.GetByID(al.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get api failed: %v", err)
	}
	if api == nil {
		return nil, status.Errorf(codes.NotFound, "api not found")
	}

	return al.mapAPIToResp(api), nil
}

func (al *APILogic) ListAPI(in *common_pb.PageReq) (*api_pb.ListAPIResp, error) {
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

	apis, countt, err := al.svcCtx.APIService.List(al.ctx, page, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list apis failed: %v", err)
	}

	resp := &api_pb.ListAPIResp{
		List:  make([]*api_pb.APIResp, 0, len(apis)),
		Total: countt,
	}

	for _, api := range apis {
		resp.List = append(resp.List, al.mapAPIToResp(api))
	}

	return resp, nil
}

func (al *APILogic) GetAPIsByGroup(in *api_pb.GetAPIsByGroupReq) (*api_pb.ListAPIResp, error) {
	apis, err := al.svcCtx.APIService.GetByGroup(al.ctx, in.GroupEn)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get apis by group failed: %v", err)
	}

	resp := &api_pb.ListAPIResp{
		List:  make([]*api_pb.APIResp, 0, len(apis)),
		Total: int64(len(apis)),
	}

	for _, api := range apis {
		resp.List = append(resp.List, al.mapAPIToResp(api))
	}

	return resp, nil
}

func (al *APILogic) CreateAPI(in *api_pb.CreateAPIReq) (*common_pb.SuccessResp, error) {
	existing, err := al.svcCtx.APIService.GetByPathAndMethod(al.ctx, in.Path, in.Method)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "check api exists failed: %v", err)
	}
	if existing != nil {
		return nil, status.Errorf(codes.AlreadyExists, "api with this path and method already exists")
	}

	api := &sysManagement.ApiModel{
		Path:        in.Path,
		Method:      in.Method,
		GroupEN:     in.GroupEn,
		GroupCN:     in.GroupCn,
		Description: in.Description,
	}

	err = al.svcCtx.APIService.Create(al.ctx, api)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create api failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (al *APILogic) UpdateAPI(in *api_pb.UpdateAPIReq) (*api_pb.APIResp, error) {
	api, err := al.svcCtx.APIService.GetByID(al.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get api failed: %v", err)
	}
	if api == nil {
		return nil, status.Errorf(codes.NotFound, "api not found")
	}

	if in.Path != nil {
		api.Path = *in.Path
	}
	if in.Method != nil {
		api.Method = *in.Method
	}
	if in.GroupEn != nil {
		api.GroupEN = *in.GroupEn
	}
	if in.GroupCn != nil {
		api.GroupCN = *in.GroupCn
	}
	if in.Description != nil {
		api.Description = *in.Description
	}

	err = al.svcCtx.APIService.Update(al.ctx, api)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update api failed: %v", err)
	}

	updatedAPI, err := al.svcCtx.APIService.GetByID(al.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get updated api failed: %v", err)
	}

	return al.mapAPIToResp(updatedAPI), nil
}

func (al *APILogic) DeleteAPI(in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid api id")
	}

	err := al.svcCtx.APIService.Delete(al.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete api failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (al *APILogic) DeleteByIds(in *common_pb.IdsReq) (*common_pb.SuccessResp, error) {
	if len(in.Ids) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "ids is empty")
	}

	ids := make([]uint, 0, len(in.Ids))
	for _, id := range in.Ids {
		ids = append(ids, uint(id))
	}

	err := al.svcCtx.APIService.DeleteByIds(al.ctx, ids)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "batch delete apis failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (al *APILogic) GetAPITree(in *api_pb.APITreeReq) (*api_pb.APITreeResp, error) {
	allAPIs, err := al.svcCtx.APIService.GetAll(al.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get all apis failed: %v", err)
	}

	checkedAPIMap := make(map[uint]struct{})
	if in.RoleId != 0 {
		perms, err := al.svcCtx.RoleService.GetPermissionDetails(al.ctx, uint(in.RoleId))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "get role permissions failed: %v", err)
		}
		for _, perm := range perms {
			if perm.Domain == sysManagement.PermissionDomainAPI {
				checkedAPIMap[perm.DomainID] = struct{}{}
			}
		}
	}

	type groupInfo struct {
		GroupCN string
		APIs    []*api_pb.APIResp
	}
	groupMap := make(map[string]*groupInfo)
	groupOrder := make([]string, 0)

	for _, api := range allAPIs {
		group, exists := groupMap[api.GroupEN]
		if !exists {
			group = &groupInfo{
				GroupCN: api.GroupCN,
				APIs:    make([]*api_pb.APIResp, 0),
			}
			groupMap[api.GroupEN] = group
			groupOrder = append(groupOrder, api.GroupEN)
		}
		group.APIs = append(group.APIs, al.mapAPIToResp(api))
	}

	checkedIDs := make([]uint64, 0, len(checkedAPIMap))
	list := make([]*api_pb.APITreeItem, 0, len(groupOrder))

	for _, groupEn := range groupOrder {
		gi := groupMap[groupEn]
		groupChecked := false
		for _, api := range gi.APIs {
			if _, ok := checkedAPIMap[uint(api.Id)]; ok {
				groupChecked = true
				checkedIDs = append(checkedIDs, api.Id)
			}
		}
		list = append(list, &api_pb.APITreeItem{
			GroupEn: groupEn,
			GroupCn: gi.GroupCN,
			Apis:    gi.APIs,
			Checked: groupChecked,
		})
	}

	return &api_pb.APITreeResp{
		List:       list,
		CheckedIds: checkedIDs,
	}, nil
}
