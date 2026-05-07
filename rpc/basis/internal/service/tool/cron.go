package tool

import (
	"context"
	"errors"
	"fmt"

	"github.com/robfig/cron/v3"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/tool"
	repoTool "td27/rpc/basis/internal/repository/tool"
)

type CronService struct {
	cronRepo repoTool.CronRepository
	scheduler *cron.Cron
	jobs map[uint]func()
}

func NewCronService(cronRepo repoTool.CronRepository, scheduler *cron.Cron) *CronService {
	service := &CronService{
		cronRepo:  cronRepo,
		scheduler: scheduler,
		jobs:      make(map[uint]func()),
	}
	service.loadAllEnabledJobs()
	return service
}

func (s *CronService) loadAllEnabledJobs() {
	jobs, err := s.cronRepo.FindAllEnabled(context.Background())
	if err != nil {
		return
	}

	for _, job := range jobs {
		if job.Open {
			s.scheduleJob(job)
		}
	}
}

func (s *CronService) scheduleJob(job *tool.CronModel) error {
	entryID, err := s.scheduler.AddFunc(job.Expression, func() {
		switch job.Method {
		case tool.CronMethod.ClearTable:
		case tool.CronMethod.ClearCache:
		case tool.CronMethod.Shell:
		}
	})
	if err != nil {
		return err
	}

	job.EntryId = int(entryID)
	return s.cronRepo.UpdateEntryID(context.Background(), job.ID, int(entryID))
}

func (s *CronService) GetByID(ctx context.Context, id uint) (*tool.CronModel, error) {
	return s.cronRepo.FindOne(ctx, id)
}

func (s *CronService) List(ctx context.Context, page *common.PageInfo) ([]*tool.CronModel, int64, error) {
	return s.cronRepo.List(ctx, page)
}

func (s *CronService) Create(ctx context.Context, cronJob *tool.CronModel) error {
	if err := s.cronRepo.Create(ctx, cronJob); err != nil {
		return err
	}

	if cronJob.Open {
		return s.scheduleJob(cronJob)
	}
	return nil
}

func (s *CronService) Update(ctx context.Context, cronJob *tool.CronModel) error {
	existing, err := s.cronRepo.FindOne(ctx, cronJob.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("cron job not found")
	}

	if existing.EntryId != 0 {
		s.scheduler.Remove(cron.EntryID(existing.EntryId))
	}

	if err := s.cronRepo.Update(ctx, cronJob); err != nil {
		return err
	}

	if cronJob.Open {
		return s.scheduleJob(cronJob)
	}
	return nil
}

func (s *CronService) ToggleStatus(ctx context.Context, id uint, open bool) error {
	job, err := s.cronRepo.FindOne(ctx, id)
	if err != nil {
		return err
	}
	if job == nil {
		return errors.New("cron job not found")
	}

	if job.Open == open {
		return nil
	}

	if job.EntryId != 0 {
		s.scheduler.Remove(cron.EntryID(job.EntryId))
	}

	if err := s.cronRepo.ToggleStatus(ctx, id, open); err != nil {
		return err
	}

	if open {
		job.Open = true
		return s.scheduleJob(job)
	}
	return nil
}

func (s *CronService) Delete(ctx context.Context, id uint) error {
	job, err := s.cronRepo.FindOne(ctx, id)
	if err != nil {
		return err
	}
	if job == nil {
		return errors.New("cron job not found")
	}

	if job.EntryId != 0 {
		s.scheduler.Remove(cron.EntryID(job.EntryId))
	}

	return s.cronRepo.Delete(ctx, id)
}

func (s *CronService) ExecuteNow(ctx context.Context, id uint) error {
	job, err := s.cronRepo.FindOne(ctx, id)
	if err != nil {
		return err
	}
	if job == nil {
		return errors.New("cron job not found")
	}

	switch job.Method {
	case tool.CronMethod.ClearTable:
	case tool.CronMethod.ClearCache:
	case tool.CronMethod.Shell:
	default:
		return fmt.Errorf("unknown cron method: %s", job.Method)
	}
	return nil
}
