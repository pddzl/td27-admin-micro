package sysManagement

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"td27/rpc/basis/internal/model/sysManagement"
	"td27/rpc/basis/internal/util"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/role_pb"
)

func (rl *RoleLogic) mapElevationToResp(e *sysManagement.RoleElevationModel) *role_pb.ElevationResp {
	if e == nil {
		return nil
	}

	resp := &role_pb.ElevationResp{
		Id:              int64(e.ID),
		UserId:          int64(e.UserID),
		RoleId:          int64(e.RoleID),
		Status:          string(e.Status),
		Reason:          e.Reason,
		DurationMinutes: int64(e.DurationMinutes),
		CreatedAt:       util.Ts(e.CreatedAt),
		UpdatedAt:       util.Ts(e.UpdatedAt),
	}
	if e.ExpiresAt != nil {
		resp.ExpiresAt = util.Ts(e.ExpiresAt)
	}
	return resp
}

func (rl *RoleLogic) CreateElevationRequest(in *role_pb.CreateElevationRequestReq) (*common_pb.SuccessResp, error) {
	_, err := rl.svcCtx.ElevService.CreateRequest(rl.ctx, uint(in.UserId), uint(in.RoleId), in.Reason, int(in.DurationMinutes))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "create elevation request failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (rl *RoleLogic) DecideElevation(in *role_pb.DecideElevationReq) (*common_pb.SuccessResp, error) {
	err := rl.svcCtx.ElevService.Decide(rl.ctx, uint(in.Id), uint(in.DecidedBy), in.Approve, nil)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "decide elevation failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (rl *RoleLogic) RevokeElevation(in *role_pb.RevokeElevationReq) (*common_pb.SuccessResp, error) {
	err := rl.svcCtx.ElevService.Revoke(rl.ctx, uint(in.Id), uint(in.RevokedBy))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "revoke elevation failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (rl *RoleLogic) ListElevations(in *common_pb.Empty) (*role_pb.ListElevationResp, error) {
	elevations, err := rl.svcCtx.ElevService.List(rl.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list elevations failed: %v", err)
	}

	resp := &role_pb.ListElevationResp{
		List:  make([]*role_pb.ElevationResp, 0, len(elevations)),
		Total: int64(len(elevations)),
	}
	for _, e := range elevations {
		resp.List = append(resp.List, rl.mapElevationToResp(e))
	}

	return resp, nil
}

func (rl *RoleLogic) GetMyElevations(in *common_pb.IdReq) (*role_pb.ListElevationResp, error) {
	elevations, err := rl.svcCtx.ElevService.ListByUser(rl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list elevations failed: %v", err)
	}

	resp := &role_pb.ListElevationResp{
		List:  make([]*role_pb.ElevationResp, 0, len(elevations)),
		Total: int64(len(elevations)),
	}
	for _, e := range elevations {
		resp.List = append(resp.List, rl.mapElevationToResp(e))
	}

	return resp, nil
}
