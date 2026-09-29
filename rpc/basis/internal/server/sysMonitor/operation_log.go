package sysMonitor

import (
	"context"

	"td27/rpc/basis/internal/logic/sysMonitor"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysMonitor/operation_log_pb"
)

type OperationLogServer struct {
	svcCtx *svc.ServiceContext
	operation_log_pb.UnimplementedOperationLogServer
}

func NewOperationLogServer(svcCtx *svc.ServiceContext) *OperationLogServer {
	return &OperationLogServer{svcCtx: svcCtx}
}

func (s *OperationLogServer) CreateOperationLog(ctx context.Context, in *operation_log_pb.CreateOperationLogReq) (*common_pb.SuccessResp, error) {
	return sysMonitor.NewOperationLogLogic(ctx, s.svcCtx).CreateOperationLog(in)
}

func (s *OperationLogServer) ListOperationLog(ctx context.Context, in *operation_log_pb.ListOperationLogReq) (*operation_log_pb.ListOperationLogResp, error) {
	return sysMonitor.NewOperationLogLogic(ctx, s.svcCtx).ListOperationLog(in)
}

func (s *OperationLogServer) CleanupExpiredLogs(ctx context.Context, in *operation_log_pb.CleanupExpiredLogsReq) (*common_pb.SuccessResp, error) {
	return sysMonitor.NewOperationLogLogic(ctx, s.svcCtx).CleanupExpiredLogs(in)
}

func (s *OperationLogServer) Delete(ctx context.Context, in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	return sysMonitor.NewOperationLogLogic(ctx, s.svcCtx).DeleteOperationLog(in)
}

func (s *OperationLogServer) DeleteByIds(ctx context.Context, in *common_pb.IdsReq) (*common_pb.SuccessResp, error) {
	return sysMonitor.NewOperationLogLogic(ctx, s.svcCtx).DeleteOperationLogByIds(in)
}
