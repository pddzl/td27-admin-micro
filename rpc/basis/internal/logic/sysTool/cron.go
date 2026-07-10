package sysTool

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"td27/rpc/basis/internal/model/common"
	sysToolModel "td27/rpc/basis/internal/model/sysTool"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/internal/util"
	
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysTool/cron_pb"
)

type CronLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCronLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CronLogic {
	return &CronLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (cl *CronLogic) mapCronToResp(cron *sysToolModel.CronModel) *cron_pb.CronResp {
	if cron == nil {
		return nil
	}
	extraParams := &cron_pb.CronExtraParams{Command: cron.ExtraParams.Command}
	for _, t := range cron.ExtraParams.TableInfo {
		extraParams.TableInfo = append(extraParams.TableInfo, &cron_pb.ClearTableParam{
			TableName: t.TableName, CompareField: t.CompareField, Interval: t.Interval,
		})
	}
	return &cron_pb.CronResp{
		Id: uint64(cron.ID), Name: cron.Name, Method: methodToProto(cron.Method),
		Expression: cron.Expression, Strategy: cron.Strategy, Open: cron.Open,
		ExtraParams: extraParams, EntryId: int32(cron.EntryId), Comment: cron.Comment,
		CreatedAt: util.ToProtoTimestamp(cron.CreatedAt), UpdatedAt: util.ToProtoTimestamp(cron.UpdatedAt),
	}
}

func (cl *CronLogic) GetCron(in *common_pb.IdReq) (*cron_pb.CronResp, error) {
	cron, err := cl.svcCtx.CronService.GetByID(cl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get cron failed: %v", err)
	}
	if cron == nil {
		return nil, status.Errorf(codes.NotFound, "cron not found")
	}
	return cl.mapCronToResp(cron), nil
}

func (cl *CronLogic) ListCron(in *common_pb.PageReq) (*cron_pb.ListCronResp, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}

	page := &common.PageInfo{Page: int(in.Page), PageSize: int(in.PageSize)}
	crons, count, err := cl.svcCtx.CronService.List(cl.ctx, page)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list cron failed: %v", err)
	}
	resp := &cron_pb.ListCronResp{List: make([]*cron_pb.CronResp, 0, len(crons)), Total: count	}
	for _, c := range crons {
		resp.List = append(resp.List, cl.mapCronToResp(c))
	}
	return resp, nil
}

func (cl *CronLogic) CreateCron(in *cron_pb.CreateCronReq) (*common_pb.SuccessResp, error) {
	cron := &sysToolModel.CronModel{Name: in.Name, Expression: in.Expression}
	cron.Method = in.Method.String()
	if in.Strategy != nil {
		cron.Strategy = *in.Strategy
	}
	if in.Open != nil {
		cron.Open = *in.Open
	}
	if in.Comment != nil {
		cron.Comment = *in.Comment
	}
	if in.ExtraParams != nil {
		ep := sysToolModel.ExtraParams{Command: in.ExtraParams.Command}
		for _, t := range in.ExtraParams.TableInfo {
			ep.TableInfo = append(ep.TableInfo, sysToolModel.ClearTable{
				TableName: t.TableName, CompareField: t.CompareField, Interval: t.Interval,
			})
		}
		cron.ExtraParams = ep
	}
	if err := cl.svcCtx.CronService.Create(cl.ctx, cron); err != nil {
		return nil, status.Errorf(codes.Internal, "create cron failed: %v", err)
	}
	return &common_pb.SuccessResp{Success: true}, nil
}

func (cl *CronLogic) UpdateCron(in *cron_pb.UpdateCronReq) (*cron_pb.CronResp, error) {
	cron, err := cl.svcCtx.CronService.GetByID(cl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get cron failed: %v", err)
	}
	if cron == nil {
		return nil, status.Errorf(codes.NotFound, "cron not found")
	}
	if in.Name != nil {
		cron.Name = *in.Name
	}
	if in.Method != nil {
		cron.Method = in.Method.String()
	}
	if in.Expression != nil {
		cron.Expression = *in.Expression
	}
	if in.Strategy != nil {
		cron.Strategy = *in.Strategy
	}
	if in.Open != nil {
		cron.Open = *in.Open
	}
	if in.Comment != nil {
		cron.Comment = *in.Comment
	}
	if err := cl.svcCtx.CronService.Update(cl.ctx, cron); err != nil {
		return nil, status.Errorf(codes.Internal, "update cron failed: %v", err)
	}
	updated, _ := cl.svcCtx.CronService.GetByID(cl.ctx, uint(in.Id))
	return cl.mapCronToResp(updated), nil
}

func (cl *CronLogic) ToggleCronStatus(in *cron_pb.ToggleCronStatusReq) (*common_pb.SuccessResp, error) {
	if err := cl.svcCtx.CronService.ToggleStatus(cl.ctx, uint(in.Id), in.Open); err != nil {
		return nil, status.Errorf(codes.Internal, "toggle cron status failed: %v", err)
	}
	return &common_pb.SuccessResp{Success: true}, nil
}

func (cl *CronLogic) ExecuteCronNow(in *cron_pb.ExecuteCronNowReq) (*common_pb.SuccessResp, error) {
	if err := cl.svcCtx.CronService.ExecuteNow(cl.ctx, uint(in.Id)); err != nil {
		return nil, status.Errorf(codes.Internal, "execute cron failed: %v", err)
	}
	return &common_pb.SuccessResp{Success: true}, nil
}

func (cl *CronLogic) DeleteCron(in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid cron id")
	}
	if err := cl.svcCtx.CronService.Delete(cl.ctx, uint(in.Id)); err != nil {
		return nil, status.Errorf(codes.Internal, "delete cron failed: %v", err)
	}
	return &common_pb.SuccessResp{Success: true}, nil
}

func (cl *CronLogic) DeleteByIds(in *common_pb.IdsReq) (*common_pb.SuccessResp, error) {
	if len(in.Ids) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "ids is empty")
	}

	ids := make([]uint, 0, len(in.Ids))
	for _, id := range in.Ids {
		ids = append(ids, uint(id))
	}

	if err := cl.svcCtx.CronService.DeleteByIds(cl.ctx, ids); err != nil {
		return nil, status.Errorf(codes.Internal, "batch delete crons failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

// methodToProto converts model method string to proto CronMethod enum
func methodToProto(method string) cron_pb.CronMethod {
	switch method {
	case "clearTable":
		return cron_pb.CronMethod_METHOD_CLEAR_TABLE
	case "clearCache":
		return cron_pb.CronMethod_METHOD_CLEAR_CACHE
	case "shell":
		return cron_pb.CronMethod_METHOD_SHELL
	default:
		return cron_pb.CronMethod_METHOD_CLEAR_TABLE
	}
}
