package sysMonitor

import (
	"context"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysMonitor"
	sysMonitorRepo "td27/rpc/basis/internal/repository/sysMonitor"
)

type OperationLogService struct {
	logRepo sysMonitorRepo.OperationLogRepository
}

func NewOperationLogService(logRepo sysMonitorRepo.OperationLogRepository) *OperationLogService {
	return &OperationLogService{
		logRepo: logRepo,
	}
}

func (s *OperationLogService) Create(ctx context.Context, log *sysMonitor.OperationLogModel) error {
	return s.logRepo.Create(ctx, log)
}

func (s *OperationLogService) List(ctx context.Context, page *common.PageInfo, userID *uint, status *int, path, method *string) ([]*sysMonitor.OperationLogModel, int64, error) {
	return s.logRepo.List(ctx, page, userID, status, path, method)
}

func (s *OperationLogService) CleanupExpired(ctx context.Context, days int) error {
	return s.logRepo.DeleteExpired(ctx, days)
}

func (s *OperationLogService) Delete(ctx context.Context, id uint) error {
	return s.logRepo.Delete(ctx, id)
}

func (s *OperationLogService) DeleteByIds(ctx context.Context, ids []uint) error {
	return s.logRepo.DeleteByIds(ctx, ids)
}
