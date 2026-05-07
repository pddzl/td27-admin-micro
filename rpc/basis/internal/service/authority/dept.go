package authority

import (
	"context"
	"errors"

	"td27/rpc/basis/internal/model/authority"
	"td27/rpc/basis/internal/model/common"
	repoAuthority "td27/rpc/basis/internal/repository/authority"
)

type DeptService struct {
	deptRepo repoAuthority.DeptRepository
}

func NewDeptService(deptRepo repoAuthority.DeptRepository) *DeptService {
	return &DeptService{
		deptRepo: deptRepo,
	}
}

func (s *DeptService) GetByID(ctx context.Context, id uint) (*authority.DeptModel, error) {
	return s.deptRepo.FindOne(ctx, id)
}

func (s *DeptService) GetAll(ctx context.Context) ([]*authority.DeptModel, error) {
	return s.deptRepo.FindAll(ctx)
}

func (s *DeptService) GetByParentID(ctx context.Context, parentID uint) ([]*authority.DeptModel, error) {
	return s.deptRepo.FindByParentID(ctx, parentID)
}

func (s *DeptService) GetDescendants(ctx context.Context, deptID uint) ([]*authority.DeptModel, error) {
	return s.deptRepo.FindDescendants(ctx, deptID)
}

func (s *DeptService) GetAncestors(ctx context.Context, path string) ([]*authority.DeptModel, error) {
	return s.deptRepo.FindAncestors(ctx, path)
}

func (s *DeptService) List(ctx context.Context, page *common.PageInfo) ([]*authority.DeptModel, int64, error) {
	return s.deptRepo.List(ctx, page)
}

func (s *DeptService) Create(ctx context.Context, dept *authority.DeptModel) error {
	if dept.ParentID != 0 {
		parent, err := s.deptRepo.FindOne(ctx, dept.ParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return errors.New("parent department not found")
		}
	}
	return s.deptRepo.Create(ctx, dept)
}

func (s *DeptService) Update(ctx context.Context, dept *authority.DeptModel) error {
	return s.deptRepo.Update(ctx, dept)
}

func (s *DeptService) Delete(ctx context.Context, id uint) error {
	descendants, err := s.deptRepo.FindDescendants(ctx, id)
	if err != nil {
		return err
	}
	if len(descendants) > 0 {
		return errors.New("cannot delete department with child departments")
	}
	return s.deptRepo.Delete(ctx, id)
}
