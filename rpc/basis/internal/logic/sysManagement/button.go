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
	"td27/rpc/basis/types/sysManagement/button_pb"
)

type ButtonLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewButtonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ButtonLogic {
	return &ButtonLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (bl *ButtonLogic) mapButtonToResp(button *sysManagement.ButtonModel) *button_pb.ButtonResp {
	if button == nil {
		return nil
	}

	return &button_pb.ButtonResp{
		Id:          int64(button.ID),
		ButtonCode:  button.ButtonCode,
		ButtonName:  button.ButtonName,
		Description: button.Description,
		PagePath:    button.PagePath,
		CreatedAt:   util.ToProtoTimestamp(button.CreatedAt),
		UpdatedAt:   util.ToProtoTimestamp(button.UpdatedAt),
	}
}

func (bl *ButtonLogic) GetButton(in *common_pb.IdReq) (*button_pb.ButtonResp, error) {
	button, err := bl.svcCtx.ButtonService.GetByID(bl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get button failed: %v", err)
	}
	if button == nil {
		return nil, status.Errorf(codes.NotFound, "button not found")
	}

	return bl.mapButtonToResp(button), nil
}

func (bl *ButtonLogic) ListButton(in *common_pb.PageReq) (*button_pb.ListButtonResp, error) {
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

	buttons, countt, err := bl.svcCtx.ButtonService.List(bl.ctx, page, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list buttons failed: %v", err)
	}

	resp := &button_pb.ListButtonResp{
		List:  make([]*button_pb.ButtonResp, 0, len(buttons)),
		Total: countt,
	}

	for _, button := range buttons {
		resp.List = append(resp.List, bl.mapButtonToResp(button))
	}

	return resp, nil
}

func (bl *ButtonLogic) GetButtonsByPagePath(in *button_pb.GetButtonsByPagePathReq) (*button_pb.ListButtonResp, error) {
	buttons, err := bl.svcCtx.ButtonService.GetByPagePath(bl.ctx, in.PagePath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get buttons by page path failed: %v", err)
	}

	resp := &button_pb.ListButtonResp{
		List:  make([]*button_pb.ButtonResp, 0, len(buttons)),
		Total: int64(len(buttons)),
	}

	for _, button := range buttons {
		resp.List = append(resp.List, bl.mapButtonToResp(button))
	}

	return resp, nil
}

func (bl *ButtonLogic) GetUserButtons(in *button_pb.GetUserButtonsReq) (*button_pb.ListButtonResp, error) {
	roleIDs := make([]uint, 0, len(in.RoleIds))
	for _, rid := range in.RoleIds {
		roleIDs = append(roleIDs, uint(rid))
	}

	buttons, err := bl.svcCtx.ButtonService.GetByRoleIDs(bl.ctx, roleIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get user buttons failed: %v", err)
	}

	resp := &button_pb.ListButtonResp{
		List:  make([]*button_pb.ButtonResp, 0, len(buttons)),
		Total: int64(len(buttons)),
	}

	for _, button := range buttons {
		resp.List = append(resp.List, bl.mapButtonToResp(button))
	}

	return resp, nil
}

func (bl *ButtonLogic) CreateButton(in *button_pb.CreateButtonReq) (*common_pb.SuccessResp, error) {
	existing, err := bl.svcCtx.ButtonService.GetByCode(bl.ctx, in.ButtonCode)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "check button exists failed: %v", err)
	}
	if existing != nil {
		return nil, status.Errorf(codes.AlreadyExists, "button with this code already exists")
	}

	button := &sysManagement.ButtonModel{
		ButtonCode:  in.ButtonCode,
		ButtonName:  in.ButtonName,
		Description: *in.Description,
		PagePath:    in.PagePath,
	}

	err = bl.svcCtx.ButtonService.Create(bl.ctx, button)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create button failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (bl *ButtonLogic) UpdateButton(in *button_pb.UpdateButtonReq) (*button_pb.ButtonResp, error) {
	button, err := bl.svcCtx.ButtonService.GetByID(bl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get button failed: %v", err)
	}
	if button == nil {
		return nil, status.Errorf(codes.NotFound, "button not found")
	}

	if in.ButtonCode != nil {
		button.ButtonCode = *in.ButtonCode
	}
	if in.ButtonName != nil {
		button.ButtonName = *in.ButtonName
	}
	if in.Description != nil {
		button.Description = *in.Description
	}
	if in.PagePath != nil {
		button.PagePath = *in.PagePath
	}

	err = bl.svcCtx.ButtonService.Update(bl.ctx, button)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update button failed: %v", err)
	}

	updatedButton, err := bl.svcCtx.ButtonService.GetByID(bl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get updated button failed: %v", err)
	}

	return bl.mapButtonToResp(updatedButton), nil
}

func (bl *ButtonLogic) BatchCheckPermission(in *button_pb.BatchCheckPermissionReq) (*button_pb.BatchCheckPermissionResp, error) {
	roleIDs := make([]uint, 0, len(in.RoleIds))
	for _, rid := range in.RoleIds {
		roleIDs = append(roleIDs, uint(rid))
	}

	results, err := bl.svcCtx.ButtonService.BatchCheckPermission(bl.ctx, in.ButtonCodes, roleIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "batch check permission failed: %v", err)
	}

	return &button_pb.BatchCheckPermissionResp{
		Results: results,
	}, nil
}

func (bl *ButtonLogic) DeleteButton(in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid button id")
	}

	err := bl.svcCtx.ButtonService.Delete(bl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete button failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}
