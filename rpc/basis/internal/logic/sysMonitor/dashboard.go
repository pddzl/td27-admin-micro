package sysMonitor

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/internal/util"
	sysMonitorModel "td27/rpc/basis/internal/model/sysMonitor"

	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysMonitor/dashboard_pb"
	"td27/rpc/basis/types/sysMonitor/operation_log_pb"
)

type DashboardLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DashboardLogic {
	return &DashboardLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DashboardLogic) mapOpsLogToResp(log *sysMonitorModel.OperationLogModel) *operation_log_pb.OperationLogResp {
	if log == nil {
		return nil
	}
	return &operation_log_pb.OperationLogResp{
		Id:        uint64(log.ID),
		Ip:        log.Ip,
		Method:    log.Method,
		Path:      log.Path,
		Status:    int32(log.Status),
		UserAgent: log.UserAgent,
		ReqParam:  log.ReqParam,
		RespData:  log.RespData,
		RespTime:  log.RespTime,
		UserId:    uint64(log.UserID),
		UserName:  log.UserName,
		CreatedAt: util.ToProtoTimestamp(log.CreatedAt),
	}
}

func (l *DashboardLogic) GetStatistics(_ *common_pb.Empty) (*dashboard_pb.DashboardStatsResp, error) {
	userCount, err := l.svcCtx.DashboardService.GetUserCount(l.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get user count failed: %v", err)
	}
	roleCount, err := l.svcCtx.DashboardService.GetRoleCount(l.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get role count failed: %v", err)
	}
	apiCount, err := l.svcCtx.DashboardService.GetAPICount(l.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get api count failed: %v", err)
	}
	deptCount, err := l.svcCtx.DashboardService.GetDeptCount(l.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get dept count failed: %v", err)
	}
	menuCount, err := l.svcCtx.DashboardService.GetMenuCount(l.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get menu count failed: %v", err)
	}
	dictCount, err := l.svcCtx.DashboardService.GetDictCount(l.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get dict count failed: %v", err)
	}
	logCount, err := l.svcCtx.DashboardService.GetLogCount(l.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get log count failed: %v", err)
	}
	fileCount, err := l.svcCtx.DashboardService.GetFileCount(l.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get file count failed: %v", err)
	}

	return &dashboard_pb.DashboardStatsResp{
		UserCount: userCount,
		RoleCount: roleCount,
		ApiCount:  apiCount,
		DeptCount: deptCount,
		MenuCount: menuCount,
		DictCount: dictCount,
		LogCount:  logCount,
		FileCount: fileCount,
	}, nil
}

func (l *DashboardLogic) GetRecentOperations(in *dashboard_pb.RecentOpsReq) (*operation_log_pb.ListOperationLogResp, error) {
	limit := int(in.Limit)
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	logs, err := l.svcCtx.DashboardService.GetRecentOperations(l.ctx, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get recent operations failed: %v", err)
	}
	resp := &operation_log_pb.ListOperationLogResp{
		List: make([]*operation_log_pb.OperationLogResp, 0, len(logs)),
	}
	for _, log := range logs {
		resp.List = append(resp.List, l.mapOpsLogToResp(log))
	}
	return resp, nil
}

func (l *DashboardLogic) GetSystemInfo(_ *common_pb.Empty) (*dashboard_pb.SystemInfoResp, error) {
	allocatedMb, totalAllocatedMb, gcCount := l.svcCtx.DashboardService.MemoryStats()
	return &dashboard_pb.SystemInfoResp{
		GoVersion:         l.svcCtx.DashboardService.GoVersion(),
		Os:                l.svcCtx.DashboardService.OS(),
		Arch:              l.svcCtx.DashboardService.Arch(),
		CpuCores:          l.svcCtx.DashboardService.CPUCores(),
		UptimeSeconds:     l.svcCtx.DashboardService.Uptime(),
		GoroutineCount:    l.svcCtx.DashboardService.GoroutineCount(),
		MemoryAllocatedMb: allocatedMb,
		TotalAllocatedMb:  totalAllocatedMb,
		GcCount:           gcCount,
	}, nil
}
