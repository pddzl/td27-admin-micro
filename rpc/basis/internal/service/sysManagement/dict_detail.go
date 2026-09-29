package sysManagement

import (
	"context"
	"errors"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
	sysManagementRepo "td27/rpc/basis/internal/repository/sysManagement"
)

type DictDetailService struct {
	detailRepo sysManagementRepo.DictDetailRepository
	dictRepo   sysManagementRepo.DictRepository
}

func NewDictDetailService(
	detailRepo sysManagementRepo.DictDetailRepository,
	dictRepo sysManagementRepo.DictRepository,
) *DictDetailService {
	return &DictDetailService{
		detailRepo: detailRepo,
		dictRepo:   dictRepo,
	}
}

func (s *DictDetailService) GetByID(ctx context.Context, id uint) (*sysManagement.DictDetailModel, error) {
	return s.detailRepo.FindOne(ctx, id)
}

func (s *DictDetailService) GetByDictID(ctx context.Context, dictID uint) ([]*sysManagement.DictDetailModel, error) {
	return s.detailRepo.FindByDictID(ctx, dictID)
}

func (s *DictDetailService) GetByDictENName(ctx context.Context, dictENName string) ([]*sysManagement.DictDetailModel, error) {
	return s.detailRepo.FindByDictENName(ctx, dictENName)
}

func (s *DictDetailService) List(ctx context.Context, page *common.PageInfo, dictID *uint) ([]*sysManagement.DictDetailModel, int64, error) {
	return s.detailRepo.List(ctx, page, dictID)
}

func (s *DictDetailService) Create(ctx context.Context, detail *sysManagement.DictDetailModel) error {
	dict, err := s.dictRepo.FindOne(ctx, uint(detail.DictModelID))
	if err != nil {
		return err
	}
	if dict == nil {
		return errors.New("dictionary not found")
	}
	return s.detailRepo.Create(ctx, detail)
}

func (s *DictDetailService) Update(ctx context.Context, detail *sysManagement.DictDetailModel) error {
	return s.detailRepo.Update(ctx, detail)
}

func (s *DictDetailService) Delete(ctx context.Context, id uint) error {
	return s.detailRepo.Delete(ctx, id)
}

func (s *DictDetailService) FlatDictDetails(ctx context.Context, dictID uint) ([]*sysManagement.DictDetailModel, error) {
	return s.detailRepo.FindByDictID(ctx, dictID)
}
