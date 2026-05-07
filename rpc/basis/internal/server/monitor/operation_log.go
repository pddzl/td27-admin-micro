package monitor

import (
	"context"
	"td27/rpc/basis/internal/logic/monitor"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/monitor/operation_log_pb"
	"td27/rpc/basis/types/common_pb"
)

type OperationLogServer struct {
	svcCtx *svc.ServiceContext
	operation_log_pb.UnimplementedOperationLogServer
}

func NewOperationLogServer(svcCtx *svc.ServiceContext) *OperationLogServer {
	return &OperationLogServer{svcCtx: svcCtx}
}

func (s *OperationLogServer) ListOperationLog(ctx context.Context, in *operation_log_pb.ListOperationLogReq) (*operation_log_pb.ListOperationLogResp, error) {
	return monitor.NewOperationLogLogic(ctx, s.svcCtx).ListOperationLog(in)
}

func (s *OperationLogServer) CleanupExpiredLogs(ctx context.Context, in *operation_log_pb.CleanupExpiredLogsReq) (*common_pb.SuccessResp, error) {
	return monitor.NewOperationLogLogic(ctx, s.svcCtx).CleanupExpiredLogs(in)
}
