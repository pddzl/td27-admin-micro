package sysManagement

import (
	"context"
	"errors"
	"fmt"

	"github.com/casbin/casbin/v2"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
	sysManagementRepo "td27/rpc/basis/internal/repository/sysManagement"
)

type RoleService struct {
	roleRepo       sysManagementRepo.RoleRepository
	permissionRepo sysManagementRepo.PermissionRepository
	userRepo       sysManagementRepo.UserRepository
	enforcer       *casbin.SyncedCachedEnforcer
}

func NewRoleService(
	roleRepo sysManagementRepo.RoleRepository,
	permissionRepo sysManagementRepo.PermissionRepository,
	userRepo sysManagementRepo.UserRepository,
	enforcer *casbin.SyncedCachedEnforcer,
) *RoleService {
	return &RoleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		userRepo:       userRepo,
		enforcer:       enforcer,
	}
}

func (s *RoleService) GetByID(ctx context.Context, id uint) (*sysManagement.RoleModel, error) {
	return s.roleRepo.FindOne(ctx, id)
}

func (s *RoleService) GetByIDWithChildren(ctx context.Context, id uint) (*sysManagement.RoleModel, error) {
	return s.roleRepo.FindOneWithChildren(ctx, id)
}

func (s *RoleService) GetAll(ctx context.Context) ([]*sysManagement.RoleModel, error) {
	return s.roleRepo.FindAll(ctx)
}

func (s *RoleService) List(ctx context.Context, page *common.PageInfo) ([]*sysManagement.RoleModel, int64, error) {
	return s.roleRepo.List(ctx, page)
}

func (s *RoleService) Create(ctx context.Context, role *sysManagement.RoleModel) error {
	return s.roleRepo.Create(ctx, role)
}

func (s *RoleService) Update(ctx context.Context, role *sysManagement.RoleModel) error {
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

	if err := s.roleRepo.Delete(ctx, id); err != nil {
		return err
	}

	// 角色删除后刷新 Casbin 策略
	if err := s.enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("reload casbin policy failed: %w", err)
	}

	return nil
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

	if err := s.roleRepo.AssignPermissions(ctx, roleID, permissionIDs); err != nil {
		return err
	}

	// 权限变更后刷新 Casbin 策略
	if err := s.enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("reload casbin policy failed: %w", err)
	}

	return nil
}

func (s *RoleService) GetPermissions(ctx context.Context, roleID uint) ([]uint, error) {
	return s.roleRepo.GetPermissions(ctx, roleID)
}

func (s *RoleService) GetPermissionDetails(ctx context.Context, roleID uint) ([]*sysManagement.PermissionModel, error) {
	return s.permissionRepo.FindByRoleID(ctx, roleID)
}
