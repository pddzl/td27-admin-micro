package authority

import (
	"context"
	"errors"

	"td27/rpc/basis/internal/model/authority"
	repoAuthority "td27/rpc/basis/internal/repository/authority"
)

type MenuService struct {
	menuRepo repoAuthority.MenuRepository
}

func NewMenuService(menuRepo repoAuthority.MenuRepository) *MenuService {
	return &MenuService{
		menuRepo: menuRepo,
	}
}

func (s *MenuService) GetByID(ctx context.Context, id uint) (*authority.MenuModel, error) {
	return s.menuRepo.FindOne(ctx, id)
}

func (s *MenuService) GetAll(ctx context.Context) ([]*authority.MenuModel, error) {
	return s.menuRepo.FindAll(ctx)
}

func (s *MenuService) GetByParentID(ctx context.Context, parentID uint) ([]*authority.MenuModel, error) {
	return s.menuRepo.FindByParentID(ctx, parentID)
}

func (s *MenuService) GetByRoleIDs(ctx context.Context, roleIDs []uint) ([]*authority.MenuModel, error) {
	return s.menuRepo.FindByRoleIDs(ctx, roleIDs)
}

func (s *MenuService) GetMenuTree(ctx context.Context) ([]*authority.MenuModel, error) {
	allMenus, err := s.menuRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(allMenus, 0), nil
}

func buildMenuTree(menus []*authority.MenuModel, parentID uint) []*authority.MenuModel {
	var tree []*authority.MenuModel
	for _, menu := range menus {
		if menu.ParentID == parentID {
			buildMenuTree(menus, menu.ID)
			tree = append(tree, menu)
		}
	}
	return tree
}

func (s *MenuService) Create(ctx context.Context, menu *authority.MenuModel) error {
	return s.menuRepo.Create(ctx, menu)
}

func (s *MenuService) Update(ctx context.Context, menu *authority.MenuModel) error {
	return s.menuRepo.Update(ctx, menu)
}

func (s *MenuService) Delete(ctx context.Context, id uint) error {
	children, err := s.menuRepo.FindByParentID(ctx, id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return errors.New("cannot delete menu with child menus")
	}
	return s.menuRepo.Delete(ctx, id)
}
