package sysMonitor

import (
	"context"

	"td27/rpc/basis/internal/logic/sysMonitor"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysMonitor/dashboard_pb"
	"td27/rpc/basis/types/sysMonitor/operation_log_pb"
)

type DashboardServer struct {
	svcCtx *svc.ServiceContext
	dashboard_pb.UnimplementedDashboardServer
}

func NewDashboardServer(svcCtx *svc.ServiceContext) *DashboardServer {
	return &DashboardServer{svcCtx: svcCtx}
}

func (s *DashboardServer) GetStatistics(ctx context.Context, in *common_pb.Empty) (*dashboard_pb.DashboardStatsResp, error) {
	return sysMonitor.NewDashboardLogic(ctx, s.svcCtx).GetStatistics(in)
}

func (s *DashboardServer) GetRecentOperations(ctx context.Context, in *dashboard_pb.RecentOpsReq) (*operation_log_pb.ListOperationLogResp, error) {
	return sysMonitor.NewDashboardLogic(ctx, s.svcCtx).GetRecentOperations(in)
}

func (s *DashboardServer) GetSystemInfo(ctx context.Context, in *common_pb.Empty) (*dashboard_pb.SystemInfoResp, error) {
	return sysMonitor.NewDashboardLogic(ctx, s.svcCtx).GetSystemInfo(in)
}
