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
	"td27/rpc/basis/types/authority/dict_pb"
	"td27/rpc/basis/types/common_pb"
)

type DictLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictLogic {
	return &DictLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (dl *DictLogic) mapDictToResp(dict *authority.DictModel) *dict_pb.DictResp {
	if dict == nil {
		return nil
	}

	details := make([]*dict_pb.DictDetailResp, 0, len(dict.DictDetails))
	for _, detail := range dict.DictDetails {
		details = append(details, &dict_pb.DictDetailResp{
			Id:          uint64(detail.ID),
			Label:       detail.Label,
			Value:       detail.Value,
			Sort:        int32(detail.Sort),
			Description: detail.Description,
			CreatedAt:   util.ToProtoTimestamp(detail.CreatedAt),
			UpdatedAt:   util.ToProtoTimestamp(detail.UpdatedAt),
		})
	}

	return &dict_pb.DictResp{
		Id:        uint64(dict.ID),
		CnName:    dict.CNName,
		EnName:    dict.ENName,
		Details:   details,
		CreatedAt: util.ToProtoTimestamp(dict.CreatedAt),
		UpdatedAt: util.ToProtoTimestamp(dict.UpdatedAt),
	}
}

func (dl *DictLogic) GetDict(in *common_pb.IdReq) (*dict_pb.DictResp, error) {
	dict, err := dl.svcCtx.DictService.GetByID(dl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get dictionary failed: %v", err)
	}
	if dict == nil {
		return nil, status.Errorf(codes.NotFound, "dictionary not found")
	}

	return dl.mapDictToResp(dict), nil
}

func (dl *DictLogic) GetDictByENName(in *dict_pb.GetDictByENNameReq) (*dict_pb.DictResp, error) {
	dict, err := dl.svcCtx.DictService.GetByENName(dl.ctx, in.EnName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get dictionary by name failed: %v", err)
	}
	if dict == nil {
		return nil, status.Errorf(codes.NotFound, "dictionary not found")
	}

	return dl.mapDictToResp(dict), nil
}

func (dl *DictLogic) ListDict(in *common_pb.PageReq) (*dict_pb.ListDictResp, error) {
	page := &common.PageInfo{
		Page:     int(in.Page),
		PageSize: int(in.PageSize),
	}

	dicts, count, err := dl.svcCtx.DictService.List(dl.ctx, page)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list dictionaries failed: %v", err)
	}

	resp := &dict_pb.ListDictResp{
		List:  make([]*dict_pb.DictResp, 0, len(dicts)),
		Total: count,
	}

	for _, dict := range dicts {
		resp.List = append(resp.List, dl.mapDictToResp(dict))
	}

	return resp, nil
}

func (dl *DictLogic) CreateDict(in *dict_pb.CreateDictReq) (*common_pb.SuccessResp, error) {
	existing, err := dl.svcCtx.DictService.GetByENName(dl.ctx, in.EnName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "check dictionary exists failed: %v", err)
	}
	if existing != nil {
		return nil, status.Errorf(codes.AlreadyExists, "dictionary with this english name already exists")
	}

	dict := &authority.DictModel{
		CNName: in.CnName,
		ENName: in.EnName,
	}

	err = dl.svcCtx.DictService.Create(dl.ctx, dict)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create dictionary failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (dl *DictLogic) UpdateDict(in *dict_pb.UpdateDictReq) (*dict_pb.DictResp, error) {
	dict, err := dl.svcCtx.DictService.GetByID(dl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get dictionary failed: %v", err)
	}
	if dict == nil {
		return nil, status.Errorf(codes.NotFound, "dictionary not found")
	}

	if in.CnName != nil {
		dict.CNName = *in.CnName
	}
	if in.EnName != nil {
		dict.ENName = *in.EnName
	}

	err = dl.svcCtx.DictService.Update(dl.ctx, dict)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update dictionary failed: %v", err)
	}

	updatedDict, err := dl.svcCtx.DictService.GetByID(dl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get updated dictionary failed: %v", err)
	}

	return dl.mapDictToResp(updatedDict), nil
}

func (dl *DictLogic) DeleteDict(in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid dictionary id")
	}

	err := dl.svcCtx.DictService.Delete(dl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete dictionary failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (dl *DictLogic) CreateDictDetail(in *dict_pb.CreateDictDetailReq) (*common_pb.SuccessResp, error) {
	dict, err := dl.svcCtx.DictService.GetByID(dl.ctx, uint(in.DictId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get dictionary failed: %v", err)
	}
	if dict == nil {
		return nil, status.Errorf(codes.NotFound, "dictionary not found")
	}

	parentID := int(*in.ParentId)
	detail := &authority.DictDetailModel{
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

func (dl *DictLogic) UpdateDictDetail(in *dict_pb.UpdateDictDetailReq) (*dict_pb.DictDetailResp, error) {
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
		Id:          uint64(updatedDetail.ID),
		Label:       updatedDetail.Label,
		Value:       updatedDetail.Value,
		Sort:        int32(updatedDetail.Sort),
		Description: updatedDetail.Description,
		CreatedAt:   util.ToProtoTimestamp(updatedDetail.CreatedAt),
		UpdatedAt:   util.ToProtoTimestamp(updatedDetail.UpdatedAt),
	}, nil
}

func (dl *DictLogic) DeleteDictDetail(in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid detail id")
	}

	err := dl.svcCtx.DictDetRepo.Delete(dl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete detail failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}
