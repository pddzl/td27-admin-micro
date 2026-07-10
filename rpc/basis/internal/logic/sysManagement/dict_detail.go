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
	"td27/rpc/basis/types/sysManagement/dict_detail_pb"
	"td27/rpc/basis/types/sysManagement/dict_pb"
)

type DictDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDetailLogic {
	return &DictDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (dl *DictDetailLogic) CreateDictDetail(in *dict_detail_pb.CreateDictDetailReq) (*common_pb.SuccessResp, error) {
	dict, err := dl.svcCtx.DictService.GetByID(dl.ctx, uint(in.DictId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get dictionary failed: %v", err)
	}
	if dict == nil {
		return nil, status.Errorf(codes.NotFound, "dictionary not found")
	}

	parentID := int(*in.ParentId)
	detail := &sysManagement.DictDetailModel{
		DictModelID: int(in.DictId),
		Label:       in.Label,
		Value:       in.Value,
		Sort:        int(*in.Sort),
		Description: *in.Description,
		ParentID:    &parentID,
	}

	err = dl.svcCtx.DictDetRepo.Create(dl.ctx, detail)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create dictionary detail failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (dl *DictDetailLogic) UpdateDictDetail(in *dict_detail_pb.UpdateDictDetailReq) (*dict_pb.DictDetailResp, error) {
	detail, err := dl.svcCtx.DictDetRepo.FindOne(dl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get dictionary detail failed: %v", err)
	}
	if detail == nil {
		return nil, status.Errorf(codes.NotFound, "dictionary detail not found")
	}

	if in.Label != nil {
		detail.Label = *in.Label
	}
	if in.Value != nil {
		detail.Value = *in.Value
	}
	if in.Sort != nil {
		detail.Sort = int(*in.Sort)
	}
	if in.Description != nil {
		detail.Description = *in.Description
	}
	if in.ParentId != nil {
		parentID := int(*in.ParentId)
		detail.ParentID = &parentID
	}

	err = dl.svcCtx.DictDetRepo.Update(dl.ctx, detail)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update dictionary detail failed: %v", err)
	}

	updatedDetail, err := dl.svcCtx.DictDetRepo.FindOne(dl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get updated detail failed: %v", err)
	}

	return &dict_pb.DictDetailResp{
		Id:          int64(updatedDetail.ID),
		Label:       updatedDetail.Label,
		Value:       updatedDetail.Value,
		Sort:        int32(updatedDetail.Sort),
		Description: updatedDetail.Description,
		CreatedAt:   util.ToProtoTimestamp(updatedDetail.CreatedAt),
		UpdatedAt:   util.ToProtoTimestamp(updatedDetail.UpdatedAt),
	}, nil
}

func (dl *DictDetailLogic) DeleteDictDetail(in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid detail id")
	}

	err := dl.svcCtx.DictDetRepo.Delete(dl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete detail failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (dl *DictDetailLogic) FlatDictDetails(in *dict_detail_pb.FlatDictDetailsReq) (*dict_detail_pb.FlatDictDetailsResp, error) {
	if in.DictId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid dict id")
	}

	dict, err := dl.svcCtx.DictService.GetByID(dl.ctx, uint(in.DictId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get dictionary failed: %v", err)
	}
	if dict == nil {
		return nil, status.Errorf(codes.NotFound, "dictionary not found")
	}

	details, err := dl.svcCtx.DictDetRepo.FindByDictID(dl.ctx, uint(in.DictId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get dict details failed: %v", err)
	}

	list := make([]*dict_pb.DictDetailResp, 0, len(details))
	for _, detail := range details {
		list = append(list, &dict_pb.DictDetailResp{
			Id:          int64(detail.ID),
			Label:       detail.Label,
			Value:       detail.Value,
			Sort:        int32(detail.Sort),
			Description: detail.Description,
			CreatedAt:   util.ToProtoTimestamp(detail.CreatedAt),
			UpdatedAt:   util.ToProtoTimestamp(detail.UpdatedAt),
		})
	}

	return &dict_detail_pb.FlatDictDetailsResp{List: list}, nil
}

func (dl *DictDetailLogic) ListDictDetail(in *dict_detail_pb.ListDictDetailReq) (*dict_detail_pb.ListDictDetailResp, error) {
	page := &common.PageInfo{
		Page:     int(in.Page),
		PageSize: int(in.PageSize),
	}

	var dictID *uint
	if in.DictId != nil {
		id := uint(*in.DictId)
		dictID = &id
	}

	details, count, err := dl.svcCtx.DictDetRepo.List(dl.ctx, page, dictID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list dict details failed: %v", err)
	}

	list := make([]*dict_pb.DictDetailResp, 0, len(details))
	for _, detail := range details {
		list = append(list, &dict_pb.DictDetailResp{
			Id:          int64(detail.ID),
			Label:       detail.Label,
			Value:       detail.Value,
			Sort:        int32(detail.Sort),
			Description: detail.Description,
			CreatedAt:   util.ToProtoTimestamp(detail.CreatedAt),
			UpdatedAt:   util.ToProtoTimestamp(detail.UpdatedAt),
		})
	}

	return &dict_detail_pb.ListDictDetailResp{
		List:  list,
		Total: count,
	}, nil
}
