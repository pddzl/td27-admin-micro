package sysManagement

import (
	"context"

	"td27/rpc/basis/internal/logic/sysManagement"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/button_pb"
)

type ButtonServer struct {
	svcCtx *svc.ServiceContext
	button_pb.UnimplementedButtonServer
}

func NewButtonServer(svcCtx *svc.ServiceContext) *ButtonServer {
	return &ButtonServer{
		svcCtx: svcCtx,
	}
}

func (bs *ButtonServer) GetButton(ctx context.Context, in *common_pb.IdReq) (*button_pb.ButtonResp, error) {
	bl := sysManagement.NewButtonLogic(ctx, bs.svcCtx)
	return bl.GetButton(in)
}

func (bs *ButtonServer) ListButton(ctx context.Context, in *common_pb.PageReq) (*button_pb.ListButtonResp, error) {
	bl := sysManagement.NewButtonLogic(ctx, bs.svcCtx)
	return bl.ListButton(in)
}

func (bs *ButtonServer) GetButtonsByPagePath(ctx context.Context, in *button_pb.GetButtonsByPagePathReq) (*button_pb.ListButtonResp, error) {
	bl := sysManagement.NewButtonLogic(ctx, bs.svcCtx)
	return bl.GetButtonsByPagePath(in)
}

func (bs *ButtonServer) GetUserButtons(ctx context.Context, in *button_pb.GetUserButtonsReq) (*button_pb.ListButtonResp, error) {
	bl := sysManagement.NewButtonLogic(ctx, bs.svcCtx)
	return bl.GetUserButtons(in)
}

func (bs *ButtonServer) CreateButton(ctx context.Context, in *button_pb.CreateButtonReq) (*common_pb.SuccessResp, error) {
	bl := sysManagement.NewButtonLogic(ctx, bs.svcCtx)
	return bl.CreateButton(in)
}

func (bs *ButtonServer) UpdateButton(ctx context.Context, in *button_pb.UpdateButtonReq) (*button_pb.ButtonResp, error) {
	bl := sysManagement.NewButtonLogic(ctx, bs.svcCtx)
	return bl.UpdateButton(in)
}

func (bs *ButtonServer) DeleteButton(ctx context.Context, in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	bl := sysManagement.NewButtonLogic(ctx, bs.svcCtx)
	return bl.DeleteButton(in)
}

func (bs *ButtonServer) BatchCheckPermission(ctx context.Context, in *button_pb.BatchCheckPermissionReq) (*button_pb.BatchCheckPermissionResp, error) {
	bl := sysManagement.NewButtonLogic(ctx, bs.svcCtx)
	return bl.BatchCheckPermission(in)
}
