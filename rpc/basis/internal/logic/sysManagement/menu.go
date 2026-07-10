package sysManagement

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"td27/rpc/basis/internal/model/sysManagement"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/internal/util"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/menu_pb"
)

type MenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuLogic {
	return &MenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (ml *MenuLogic) mapMenuToResp(menu *sysManagement.MenuModel) *menu_pb.MenuResp {
	if menu == nil {
		return nil
	}

	return &menu_pb.MenuResp{
		Id:         int64(menu.ID),
		MenuName:   menu.MenuName,
		Icon:       menu.Icon,
		Path:       menu.Path,
		Component:  menu.Component,
		Redirect:   menu.Redirect,
		ParentId:   int64(menu.ParentID),
		Sort:       uint32(menu.Sort),
		Hidden:     menu.Hidden,
		KeepAlive:  menu.KeepAlive,
		Affix:      menu.Affix,
		AlwaysShow: menu.AlwaysShow,
		Title:      menu.Title,
		CreatedAt:  util.ToProtoTimestamp(menu.CreatedAt),
		UpdatedAt:  util.ToProtoTimestamp(menu.UpdatedAt),
	}
}

func (ml *MenuLogic) mapMenuToTree(menu *sysManagement.MenuModel, allMenus []*sysManagement.MenuModel) *menu_pb.MenuTreeResp {
	resp := &menu_pb.MenuTreeResp{
		Menu:     ml.mapMenuToResp(menu),
		Children: make([]*menu_pb.MenuTreeResp, 0),
	}

	for _, child := range allMenus {
		if child.ParentID == menu.ID {
			resp.Children = append(resp.Children, ml.mapMenuToTree(child, allMenus))
		}
	}

	return resp
}

func (ml *MenuLogic) GetMenu(in *common_pb.IdReq) (*menu_pb.MenuResp, error) {
	menu, err := ml.svcCtx.MenuService.GetByID(ml.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get menu failed: %v", err)
	}
	if menu == nil {
		return nil, status.Errorf(codes.NotFound, "menu not found")
	}

	return ml.mapMenuToResp(menu), nil
}

func (ml *MenuLogic) GetMenuTree(in *common_pb.Empty) (*menu_pb.GetMenuTreeResp, error) {
	allMenus, err := ml.svcCtx.MenuService.GetAll(ml.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get all menus failed: %v", err)
	}

	tree := make([]*menu_pb.MenuTreeResp, 0)
	for _, menu := range allMenus {
		if menu.ParentID == 0 {
			tree = append(tree, ml.mapMenuToTree(menu, allMenus))
		}
	}

	return &menu_pb.GetMenuTreeResp{
		Tree: tree,
	}, nil
}

func (ml *MenuLogic) GetUserMenus(in *menu_pb.GetUserMenusReq) (*menu_pb.GetMenuTreeResp, error) {
	roleIDs := make([]uint, 0, len(in.RoleIds))
	for _, rid := range in.RoleIds {
		roleIDs = append(roleIDs, uint(rid))
	}

	menus, err := ml.svcCtx.MenuService.GetByRoleIDs(ml.ctx, roleIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get user menus failed: %v", err)
	}

	allMenus, err := ml.svcCtx.MenuService.GetAll(ml.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get all menus failed: %v", err)
	}

	tree := make([]*menu_pb.MenuTreeResp, 0)
	menuMap := make(map[uint]*sysManagement.MenuModel)
	for _, menu := range menus {
		menuMap[menu.ID] = menu
	}

	for _, menu := range menus {
		if menu.ParentID == 0 || menuMap[menu.ParentID] == nil {
			tree = append(tree, ml.mapMenuToTree(menu, allMenus))
		}
	}

	return &menu_pb.GetMenuTreeResp{
		Tree: tree,
	}, nil
}

func (ml *MenuLogic) CreateMenu(in *menu_pb.CreateMenuReq) (*common_pb.SuccessResp, error) {
	menu := &sysManagement.MenuModel{
		MenuName:   in.MenuName,
		Icon:       *in.Icon,
		Path:       in.Path,
		Component:  *in.Component,
		Redirect:   *in.Redirect,
		ParentID:   uint(*in.ParentId),
		Sort:       uint(*in.Sort),
		Hidden:     *in.Hidden,
		KeepAlive:  *in.KeepAlive,
		Affix:      *in.Affix,
		AlwaysShow: *in.AlwaysShow,
		Title:      in.Title,
	}

	err := ml.svcCtx.MenuService.Create(ml.ctx, menu)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create menu failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (ml *MenuLogic) UpdateMenu(in *menu_pb.UpdateMenuReq) (*menu_pb.MenuResp, error) {
	menu, err := ml.svcCtx.MenuService.GetByID(ml.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get menu failed: %v", err)
	}
	if menu == nil {
		return nil, status.Errorf(codes.NotFound, "menu not found")
	}

	if in.MenuName != nil {
		menu.MenuName = *in.MenuName
	}
	if in.Icon != nil {
		menu.Icon = *in.Icon
	}
	if in.Path != nil {
		menu.Path = *in.Path
	}
	if in.Component != nil {
		menu.Component = *in.Component
	}
	if in.Redirect != nil {
		menu.Redirect = *in.Redirect
	}
	if in.ParentId != nil {
		menu.ParentID = uint(*in.ParentId)
	}
	if in.Sort != nil {
		menu.Sort = uint(*in.Sort)
	}
	if in.Hidden != nil {
		menu.Hidden = *in.Hidden
	}
	if in.KeepAlive != nil {
		menu.KeepAlive = *in.KeepAlive
	}
	if in.Affix != nil {
		menu.Affix = *in.Affix
	}
	if in.AlwaysShow != nil {
		menu.AlwaysShow = *in.AlwaysShow
	}
	if in.Title != nil {
		menu.Title = *in.Title
	}

	err = ml.svcCtx.MenuService.Update(ml.ctx, menu)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update menu failed: %v", err)
	}

	updatedMenu, err := ml.svcCtx.MenuService.GetByID(ml.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get updated menu failed: %v", err)
	}

	return ml.mapMenuToResp(updatedMenu), nil
}

func (ml *MenuLogic) DeleteMenu(in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid menu id")
	}

	err := ml.svcCtx.MenuService.Delete(ml.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete menu failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (ml *MenuLogic) ListMenu(in *common_pb.PageReq) (*menu_pb.ListMenuResp, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}

	menus, err := ml.svcCtx.MenuService.GetAll(ml.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list menus failed: %v", err)
	}
	list := make([]*menu_pb.MenuResp, 0, len(menus))
	for _, m := range menus {
		list = append(list, ml.mapMenuToResp(m))
	}
	return &menu_pb.ListMenuResp{List: list, Total: int64(len(list))}, nil
}
