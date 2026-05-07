package authority

import (
	"context"
	"errors"

	"td27/rpc/basis/internal/model/authority"
	"td27/rpc/basis/internal/model/common"
	repoAuthority "td27/rpc/basis/internal/repository/authority"
)

type DictService struct {
	dictRepo repoAuthority.DictRepository
}

func NewDictService(dictRepo repoAuthority.DictRepository) *DictService {
	return &DictService{
		dictRepo: dictRepo,
	}
}

func (s *DictService) GetByID(ctx context.Context, id uint) (*authority.DictModel, error) {
	return s.dictRepo.FindOne(ctx, id)
}

func (s *DictService) GetByENName(ctx context.Context, enName string) (*authority.DictModel, error) {
	return s.dictRepo.FindByENName(ctx, enName)
}

func (s *DictService) GetAll(ctx context.Context) ([]*authority.DictModel, error) {
	return s.dictRepo.FindAll(ctx)
}

func (s *DictService) List(ctx context.Context, page *common.PageInfo) ([]*authority.DictModel, int64, error) {
	return s.dictRepo.List(ctx, page)
}

func (s *DictService) Create(ctx context.Context, dict *authority.DictModel) error {
	existing, err := s.dictRepo.FindByENName(ctx, dict.ENName)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("dictionary with same english name already exists")
	}
	return s.dictRepo.Create(ctx, dict)
}

func (s *DictService) Update(ctx context.Context, dict *authority.DictModel) error {
	return s.dictRepo.Update(ctx, dict)
}

func (s *DictService) Delete(ctx context.Context, id uint) error {
	return s.dictRepo.Delete(ctx, id)
}
