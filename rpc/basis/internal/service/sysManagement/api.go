package sysManagement

import (
	"context"
	"errors"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
	sysManagementRepo "td27/rpc/basis/internal/repository/sysManagement"
)

type APIService struct {
	apiRepo sysManagementRepo.APIRepository
}

func NewAPIService(apiRepo sysManagementRepo.APIRepository) *APIService {
	return &APIService{
		apiRepo: apiRepo,
	}
}

func (s *APIService) GetByID(ctx context.Context, id uint) (*sysManagement.ApiModel, error) {
	return s.apiRepo.FindOne(ctx, id)
}

func (s *APIService) GetAll(ctx context.Context) ([]*sysManagement.ApiModel, error) {
	return s.apiRepo.FindAll(ctx)
}

func (s *APIService) GetByPathAndMethod(ctx context.Context, path string, method string) (*sysManagement.ApiModel, error) {
	return s.apiRepo.FindByPathAndMethod(ctx, path, method)
}

func (s *APIService) GetByGroup(ctx context.Context, groupEN string) ([]*sysManagement.ApiModel, error) {
	return s.apiRepo.FindByGroup(ctx, groupEN)
}

func (s *APIService) List(ctx context.Context, page *common.PageInfo, groupEN *string) ([]*sysManagement.ApiModel, int64, error) {
	return s.apiRepo.List(ctx, page, groupEN)
}

func (s *APIService) Create(ctx context.Context, api *sysManagement.ApiModel) error {
	existing, err := s.apiRepo.FindByPathAndMethod(ctx, api.Path, api.Method)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("API with same path and method already exists")
	}
	return s.apiRepo.Create(ctx, api)
}

func (s *APIService) Update(ctx context.Context, api *sysManagement.ApiModel) error {
	return s.apiRepo.Update(ctx, api)
}

func (s *APIService) Delete(ctx context.Context, id uint) error {
	return s.apiRepo.Delete(ctx, id)
}

func (s *APIService) DeleteByIds(ctx context.Context, ids []uint) error {
	return s.apiRepo.DeleteByIds(ctx, ids)
}
