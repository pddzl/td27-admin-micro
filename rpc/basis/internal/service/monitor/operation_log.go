package monitor

import (
	"context"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/monitor"
	repoMonitor "td27/rpc/basis/internal/repository/monitor"
)

type OperationLogService struct {
	logRepo repoMonitor.OperationLogRepository
}

func NewOperationLogService(logRepo repoMonitor.OperationLogRepository) *OperationLogService {
	return &OperationLogService{
		logRepo: logRepo,
	}
}

func (s *OperationLogService) Create(ctx context.Context, log *monitor.OperationLogModel) error {
	return s.logRepo.Create(ctx, log)
}

func (s *OperationLogService) List(ctx context.Context, page *common.PageInfo, userID *uint, status *int) ([]*monitor.OperationLogModel, int64, error) {
	return s.logRepo.List(ctx, page, userID, status)
}

func (s *OperationLogService) CleanupExpired(ctx context.Context, days int) error {
	return s.logRepo.DeleteExpired(ctx, days)
}
