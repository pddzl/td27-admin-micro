package authority

import (
	"context"
	"errors"
	"fmt"

	"td27/rpc/basis/internal/model/authority"
	"td27/rpc/basis/internal/model/common"
	repoAuthority "td27/rpc/basis/internal/repository/authority"
)

type RoleService struct {
	roleRepo       repoAuthority.RoleRepository
	permissionRepo repoAuthority.PermissionRepository
	userRepo       repoAuthority.UserRepository
}

func NewRoleService(
	roleRepo repoAuthority.RoleRepository,
	permissionRepo repoAuthority.PermissionRepository,
	userRepo repoAuthority.UserRepository,
) *RoleService {
	return &RoleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		userRepo:       userRepo,
	}
}

func (s *RoleService) GetByID(ctx context.Context, id uint) (*authority.RoleModel, error) {
	return s.roleRepo.FindOne(ctx, id)
}

func (s *RoleService) GetByIDWithChildren(ctx context.Context, id uint) (*authority.RoleModel, error) {
	return s.roleRepo.FindOneWithChildren(ctx, id)
}

func (s *RoleService) GetAll(ctx context.Context) ([]*authority.RoleModel, error) {
	return s.roleRepo.FindAll(ctx)
}

func (s *RoleService) List(ctx context.Context, page *common.PageInfo) ([]*authority.RoleModel, int64, error) {
	return s.roleRepo.List(ctx, page)
}

func (s *RoleService) Create(ctx context.Context, role *authority.RoleModel) error {
	return s.roleRepo.Create(ctx, role)
}

func (s *RoleService) Update(ctx context.Context, role *authority.RoleModel) error {
	return s.roleRepo.Update(ctx, role)
}

func (s *RoleService) Delete(ctx context.Context, id uint) error {
	userCount, err := s.userRepo.CountByRoleID(ctx, id)
	if err != nil {
		return err
	}
	if userCount > 0 {
		return errors.New("cannot delete role with assigned users")
	}

	return s.roleRepo.Delete(ctx, id)
}

func (s *RoleService) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	role, err := s.roleRepo.FindOne(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	for _, permID := range permissionIDs {
		perm, err := s.permissionRepo.FindOne(ctx, permID)
		if err != nil {
			return err
		}
		if perm == nil {
			return fmt.Errorf("permission with id %d not found", permID)
		}
	}

	return s.roleRepo.AssignPermissions(ctx, roleID, permissionIDs)
}

func (s *RoleService) GetPermissions(ctx context.Context, roleID uint) ([]uint, error) {
	return s.roleRepo.GetPermissions(ctx, roleID)
}

func (s *RoleService) GetPermissionDetails(ctx context.Context, roleID uint) ([]*authority.PermissionModel, error) {
	return s.permissionRepo.FindByRoleID(ctx, roleID)
}
