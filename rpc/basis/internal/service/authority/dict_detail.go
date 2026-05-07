package authority

import (
	"context"
	"errors"

	"td27/rpc/basis/internal/model/authority"
	"td27/rpc/basis/internal/model/common"
	repoAuthority "td27/rpc/basis/internal/repository/authority"
)

type DictDetailService struct {
	detailRepo repoAuthority.DictDetailRepository
	dictRepo   repoAuthority.DictRepository
}

func NewDictDetailService(
	detailRepo repoAuthority.DictDetailRepository,
	dictRepo repoAuthority.DictRepository,
) *DictDetailService {
	return &DictDetailService{
		detailRepo: detailRepo,
		dictRepo:   dictRepo,
	}
}

func (s *DictDetailService) GetByID(ctx context.Context, id uint) (*authority.DictDetailModel, error) {
	return s.detailRepo.FindOne(ctx, id)
}

func (s *DictDetailService) GetByDictID(ctx context.Context, dictID uint) ([]*authority.DictDetailModel, error) {
	return s.detailRepo.FindByDictID(ctx, dictID)
}

func (s *DictDetailService) GetByDictENName(ctx context.Context, dictENName string) ([]*authority.DictDetailModel, error) {
	return s.detailRepo.FindByDictENName(ctx, dictENName)
}

func (s *DictDetailService) List(ctx context.Context, page *common.PageInfo, dictID *uint) ([]*authority.DictDetailModel, int64, error) {
	return s.detailRepo.List(ctx, page, dictID)
}

func (s *DictDetailService) Create(ctx context.Context, detail *authority.DictDetailModel) error {
	dict, err := s.dictRepo.FindOne(ctx, uint(detail.DictModelID))
	if err != nil {
		return err
	}
	if dict == nil {
		return errors.New("dictionary not found")
	}
	return s.detailRepo.Create(ctx, detail)
}

func (s *DictDetailService) Update(ctx context.Context, detail *authority.DictDetailModel) error {
	return s.detailRepo.Update(ctx, detail)
}

func (s *DictDetailService) Delete(ctx context.Context, id uint) error {
	return s.detailRepo.Delete(ctx, id)
}
