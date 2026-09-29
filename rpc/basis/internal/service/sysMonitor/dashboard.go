package sysMonitor

import (
	"context"
	"runtime"
	"sync"
	"time"

	sysMonitorModel "td27/rpc/basis/internal/model/sysMonitor"
	sysMonitorRepo "td27/rpc/basis/internal/repository/sysMonitor"
)

var (
	startTime = time.Now()
	startOnce sync.Once
)

func init() {
	startOnce.Do(func() {
		startTime = time.Now()
	})
}

// DashboardService handles dashboard business logic
type DashboardService struct {
	dashRepo sysMonitorRepo.DashboardRepository
}

// NewDashboardService creates a new dashboard service instance
func NewDashboardService(dashRepo sysMonitorRepo.DashboardRepository) *DashboardService {
	return &DashboardService{
		dashRepo: dashRepo,
	}
}

// GetUserCount returns the total number of users
func (s *DashboardService) GetUserCount(ctx context.Context) (int64, error) {
	return s.dashRepo.GetUserCount(ctx)
}

// GetRoleCount returns the total number of roles
func (s *DashboardService) GetRoleCount(ctx context.Context) (int64, error) {
	return s.dashRepo.GetRoleCount(ctx)
}

// GetAPICount returns the total number of APIs
func (s *DashboardService) GetAPICount(ctx context.Context) (int64, error) {
	return s.dashRepo.GetAPICount(ctx)
}

// GetDeptCount returns the total number of departments
func (s *DashboardService) GetDeptCount(ctx context.Context) (int64, error) {
	return s.dashRepo.GetDeptCount(ctx)
}

// GetMenuCount returns the total number of menus
func (s *DashboardService) GetMenuCount(ctx context.Context) (int64, error) {
	return s.dashRepo.GetMenuCount(ctx)
}

// GetDictCount returns the total number of dictionaries
func (s *DashboardService) GetDictCount(ctx context.Context) (int64, error) {
	return s.dashRepo.GetDictCount(ctx)
}

// GetLogCount returns the total number of operation logs
func (s *DashboardService) GetLogCount(ctx context.Context) (int64, error) {
	return s.dashRepo.GetLogCount(ctx)
}

// GetFileCount returns the total number of uploaded files
func (s *DashboardService) GetFileCount(ctx context.Context) (int64, error) {
	return s.dashRepo.GetFileCount(ctx)
}

// GetRecentOperations returns the most recent operation logs
func (s *DashboardService) GetRecentOperations(ctx context.Context, limit int) ([]*sysMonitorModel.OperationLogModel, error) {
	return s.dashRepo.GetRecentOperations(ctx, limit)
}

// GoVersion returns the Go runtime version
func (s *DashboardService) GoVersion() string {
	return runtime.Version()
}

// OS returns the operating system
func (s *DashboardService) OS() string {
	return runtime.GOOS
}

// Arch returns the architecture
func (s *DashboardService) Arch() string {
	return runtime.GOARCH
}

// CPUCores returns the number of CPU cores
func (s *DashboardService) CPUCores() int32 {
	return int32(runtime.NumCPU())
}

// Uptime returns the server uptime in seconds
func (s *DashboardService) Uptime() int64 {
	return int64(time.Since(startTime).Seconds())
}

// GoroutineCount returns the number of goroutines
func (s *DashboardService) GoroutineCount() int64 {
	return int64(runtime.NumGoroutine())
}

// MemoryStats returns memory allocation stats
func (s *DashboardService) MemoryStats() (allocatedMb, totalAllocatedMb float64, gcCount int64) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	allocatedMb = float64(m.Alloc) / 1024 / 1024
	totalAllocatedMb = float64(m.TotalAlloc) / 1024 / 1024
	gcCount = int64(m.NumGC)
	return
}
