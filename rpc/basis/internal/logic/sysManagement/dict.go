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
	"td27/rpc/basis/types/sysManagement/dict_pb"
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

func (dl *DictLogic) mapDictToResp(dict *sysManagement.DictModel) *dict_pb.DictResp {
	if dict == nil {
		return nil
	}

	details := make([]*dict_pb.DictDetailResp, 0, len(dict.DictDetails))
	for _, detail := range dict.DictDetails {
		details = append(details, &dict_pb.DictDetailResp{
			Id:          int64(detail.ID),
			Label:       detail.Label,
			Value:       detail.Value,
			Sort:        int32(detail.Sort),
			Description: detail.Description,
			CreatedAt:   util.ToProtoTimestamp(detail.CreatedAt),
			UpdatedAt:   util.ToProtoTimestamp(detail.UpdatedAt),
		})
	}

	return &dict_pb.DictResp{
		Id:        int64(dict.ID),
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

	dicts, countt, err := dl.svcCtx.DictService.List(dl.ctx, page)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list dictionaries failed: %v", err)
	}

	resp := &dict_pb.ListDictResp{
		List:  make([]*dict_pb.DictResp, 0, len(dicts)),
		Total: countt,
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

	dict := &sysManagement.DictModel{
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




