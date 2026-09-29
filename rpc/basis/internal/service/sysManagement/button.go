package sysManagement

import (
	"context"
	"errors"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
	sysManagementRepo "td27/rpc/basis/internal/repository/sysManagement"
)

type ButtonService struct {
	buttonRepo sysManagementRepo.ButtonRepository
}

func NewButtonService(buttonRepo sysManagementRepo.ButtonRepository) *ButtonService {
	return &ButtonService{
		buttonRepo: buttonRepo,
	}
}

func (s *ButtonService) GetByID(ctx context.Context, id uint) (*sysManagement.ButtonModel, error) {
	return s.buttonRepo.FindOne(ctx, id)
}

func (s *ButtonService) GetByCode(ctx context.Context, code string) (*sysManagement.ButtonModel, error) {
	return s.buttonRepo.FindByCode(ctx, code)
}

func (s *ButtonService) GetByPagePath(ctx context.Context, pagePath string) ([]*sysManagement.ButtonModel, error) {
	return s.buttonRepo.FindByPagePath(ctx, pagePath)
}

func (s *ButtonService) GetByRoleIDs(ctx context.Context, roleIDs []uint) ([]*sysManagement.ButtonModel, error) {
	return s.buttonRepo.FindByRoleIDs(ctx, roleIDs)
}

func (s *ButtonService) List(ctx context.Context, page *common.PageInfo, pagePath *string) ([]*sysManagement.ButtonModel, int64, error) {
	return s.buttonRepo.List(ctx, page, pagePath)
}

func (s *ButtonService) Create(ctx context.Context, button *sysManagement.ButtonModel) error {
	existing, err := s.buttonRepo.FindByCode(ctx, button.ButtonCode)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("button with same code already exists")
	}
	return s.buttonRepo.Create(ctx, button)
}

func (s *ButtonService) Update(ctx context.Context, button *sysManagement.ButtonModel) error {
	return s.buttonRepo.Update(ctx, button)
}

func (s *ButtonService) Delete(ctx context.Context, id uint) error {
	return s.buttonRepo.Delete(ctx, id)
}

// BatchCheckPermission checks if the given roles have access to the specified button codes.
// Currently returns true for all codes as a simple pass-through.
// TODO: Implement proper Casbin permission check once PermissionAdapter is ready.
// The expected implementation should use the PermissionService:
//
//	permService.CheckPermission(ctx, roleIDs, buttonCode, sysManagement.ActionView)
func (s *ButtonService) BatchCheckPermission(ctx context.Context, buttonCodes []string, roleIDs []uint) (map[string]bool, error) {
	return s.buttonRepo.BatchCheckPermission(ctx, roleIDs, buttonCodes)
}
