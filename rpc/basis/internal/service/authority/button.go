package authority

import (
	"context"
	"errors"

	"td27/rpc/basis/internal/model/authority"
	"td27/rpc/basis/internal/model/common"
	repoAuthority "td27/rpc/basis/internal/repository/authority"
)

type ButtonService struct {
	buttonRepo repoAuthority.ButtonRepository
}

func NewButtonService(buttonRepo repoAuthority.ButtonRepository) *ButtonService {
	return &ButtonService{
		buttonRepo: buttonRepo,
	}
}

func (s *ButtonService) GetByID(ctx context.Context, id uint) (*authority.ButtonModel, error) {
	return s.buttonRepo.FindOne(ctx, id)
}

func (s *ButtonService) GetByCode(ctx context.Context, code string) (*authority.ButtonModel, error) {
	return s.buttonRepo.FindByCode(ctx, code)
}

func (s *ButtonService) GetByPagePath(ctx context.Context, pagePath string) ([]*authority.ButtonModel, error) {
	return s.buttonRepo.FindByPagePath(ctx, pagePath)
}

func (s *ButtonService) GetByRoleIDs(ctx context.Context, roleIDs []uint) ([]*authority.ButtonModel, error) {
	return s.buttonRepo.FindByRoleIDs(ctx, roleIDs)
}

func (s *ButtonService) List(ctx context.Context, page *common.PageInfo, pagePath *string) ([]*authority.ButtonModel, int64, error) {
	return s.buttonRepo.List(ctx, page, pagePath)
}

func (s *ButtonService) Create(ctx context.Context, button *authority.ButtonModel) error {
	existing, err := s.buttonRepo.FindByCode(ctx, button.ButtonCode)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("button with same code already exists")
	}
	return s.buttonRepo.Create(ctx, button)
}

func (s *ButtonService) Update(ctx context.Context, button *authority.ButtonModel) error {
	return s.buttonRepo.Update(ctx, button)
}

func (s *ButtonService) Delete(ctx context.Context, id uint) error {
	return s.buttonRepo.Delete(ctx, id)
}
