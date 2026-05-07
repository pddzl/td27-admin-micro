package authority

import (
	"context"

	"td27/rpc/basis/internal/logic/authority"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/authority/menu_pb"
	"td27/rpc/basis/types/common_pb"
)

type MenuServer struct {
	svcCtx *svc.ServiceContext
	menu_pb.UnimplementedMenuServer
}

func NewMenuServer(svcCtx *svc.ServiceContext) *MenuServer {
	return &MenuServer{
		svcCtx: svcCtx,
	}
}

func (ms *MenuServer) GetMenu(ctx context.Context, in *common_pb.IdReq) (*menu_pb.MenuResp, error) {
	ml := authority.NewMenuLogic(ctx, ms.svcCtx)
	return ml.GetMenu(in)
}

func (ms *MenuServer) GetMenuTree(ctx context.Context, in *common_pb.Empty) (*menu_pb.GetMenuTreeResp, error) {
	ml := authority.NewMenuLogic(ctx, ms.svcCtx)
	return ml.GetMenuTree(in)
}

func (ms *MenuServer) GetUserMenus(ctx context.Context, in *menu_pb.GetUserMenusReq) (*menu_pb.GetMenuTreeResp, error) {
	ml := authority.NewMenuLogic(ctx, ms.svcCtx)
	return ml.GetUserMenus(in)
}

func (ms *MenuServer) CreateMenu(ctx context.Context, in *menu_pb.CreateMenuReq) (*common_pb.SuccessResp, error) {
	ml := authority.NewMenuLogic(ctx, ms.svcCtx)
	return ml.CreateMenu(in)
}

func (ms *MenuServer) UpdateMenu(ctx context.Context, in *menu_pb.UpdateMenuReq) (*menu_pb.MenuResp, error) {
	ml := authority.NewMenuLogic(ctx, ms.svcCtx)
	return ml.UpdateMenu(in)
}

func (ms *MenuServer) DeleteMenu(ctx context.Context, in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	ml := authority.NewMenuLogic(ctx, ms.svcCtx)
	return ml.DeleteMenu(in)
}
