package monitor

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	monitorModel "td27/rpc/basis/internal/model/monitor"
	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/internal/util"
	"td27/rpc/basis/types/monitor/operation_log_pb"
	"td27/rpc/basis/types/common_pb"
)

type OperationLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOperationLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationLogLogic {
	return &OperationLogLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *OperationLogLogic) mapToResp(log *monitorModel.OperationLogModel) *operation_log_pb.OperationLogResp {
	if log == nil { return nil }
	return &operation_log_pb.OperationLogResp{
		Id: uint64(log.ID), Ip: log.Ip, Method: log.Method, Path: log.Path,
		Status: int32(log.Status), UserAgent: log.UserAgent, ReqParam: log.ReqParam,
		RespData: log.RespData, RespTime: log.RespTime, UserId: uint64(log.UserID),
		UserName: log.UserName, CreatedAt: util.ToProtoTimestamp(log.CreatedAt),
	}
}

func (l *OperationLogLogic) ListOperationLog(in *operation_log_pb.ListOperationLogReq) (*operation_log_pb.ListOperationLogResp, error) {
	page := &common.PageInfo{Page: int(in.Page.Page), PageSize: int(in.Page.PageSize)}
	var userID *uint
	if in.UserId != nil { v := uint(*in.UserId); userID = &v }
	var statusVal *int
	if in.Status != nil { v := int(*in.Status); statusVal = &v }
	logs, count, err := l.svcCtx.LogService.List(l.ctx, page, userID, statusVal)
	if err != nil { return nil, status.Errorf(codes.Internal, "list logs failed: %v", err) }
	resp := &operation_log_pb.ListOperationLogResp{List: make([]*operation_log_pb.OperationLogResp, 0, len(logs)), Total: count}
	for _, log := range logs { resp.List = append(resp.List, l.mapToResp(log)) }
	return resp, nil
}

func (l *OperationLogLogic) CleanupExpiredLogs(in *operation_log_pb.CleanupExpiredLogsReq) (*common_pb.SuccessResp, error) {
	if in.Days <= 0 { return nil, status.Errorf(codes.InvalidArgument, "days must be positive") }
	if err := l.svcCtx.LogService.CleanupExpired(l.ctx, int(in.Days)); err != nil {
		return nil, status.Errorf(codes.Internal, "cleanup logs failed: %v", err)
	}
	return &common_pb.SuccessResp{Success: true}, nil
}
