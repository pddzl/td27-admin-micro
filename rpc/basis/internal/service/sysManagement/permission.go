package sysManagement

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/casbin/casbin/v2"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
	sysManagementRepo "td27/rpc/basis/internal/repository/sysManagement"
)

type PermissionService struct {
	permRepo sysManagementRepo.PermissionRepository
	roleRepo sysManagementRepo.RoleRepository
	enforcer *casbin.SyncedCachedEnforcer
}

func NewPermissionService(
	permRepo sysManagementRepo.PermissionRepository,
	roleRepo sysManagementRepo.RoleRepository,
	enforcer *casbin.SyncedCachedEnforcer,
) *PermissionService {
	return &PermissionService{
		permRepo: permRepo,
		roleRepo: roleRepo,
		enforcer: enforcer,
	}
}

func (s *PermissionService) GetByID(ctx context.Context, id uint) (*sysManagement.PermissionModel, error) {
	return s.permRepo.FindOne(ctx, id)
}

func (s *PermissionService) GetAll(ctx context.Context) ([]*sysManagement.PermissionModel, error) {
	return s.permRepo.FindAll(ctx)
}

func (s *PermissionService) GetByDomain(ctx context.Context, domain sysManagement.PermissionDomain) ([]*sysManagement.PermissionModel, error) {
	return s.permRepo.FindByDomain(ctx, domain)
}

func (s *PermissionService) GetByRoleID(ctx context.Context, roleID uint) ([]*sysManagement.PermissionModel, error) {
	return s.permRepo.FindByRoleID(ctx, roleID)
}

func (s *PermissionService) GetByResourceAndAction(ctx context.Context, resource string, action sysManagement.Action) (*sysManagement.PermissionModel, error) {
	return s.permRepo.FindByResourceAndAction(ctx, resource, action)
}

func (s *PermissionService) List(ctx context.Context, page *common.PageInfo, domain *sysManagement.PermissionDomain) ([]*sysManagement.PermissionModel, int64, error) {
	perms, err := s.permRepo.FindAll(ctx)
	if err != nil {
		return nil, 0, err
	}
	return perms, 0, nil
}

func (s *PermissionService) Create(ctx context.Context, perm *sysManagement.PermissionModel) error {
	if perm.Effect == "" {
		perm.Effect = sysManagement.EffectAllow
	}
	existing, err := s.permRepo.FindByResourceAndAction(ctx, perm.Resource, perm.Action)
	if err != nil {
		return err
	}
	if existing != nil && existing.Effect == perm.Effect {
		return errors.New("permission with same resource, action and effect already exists")
	}

	return s.permRepo.Create(ctx, perm)
}

func (s *PermissionService) Update(ctx context.Context, perm *sysManagement.PermissionModel) error {
	return s.permRepo.Update(ctx, perm)
}

func (s *PermissionService) Delete(ctx context.Context, id uint) error {
	return s.permRepo.Delete(ctx, id)
}

func (s *PermissionService) CheckPermission(ctx context.Context, roleIDs []uint, resource string, action sysManagement.Action) (bool, error) {
	for _, roleID := range roleIDs {
		sub := fmt.Sprintf("%d", roleID)
		allowed, err := s.enforcer.Enforce(sub, resource, action.String())
		if err != nil {
			return false, err
		}
		if allowed {
			return true, nil
		}
	}
	return false, nil
}

// GetRolesForUser resolves a user's role ids live from the database, so
// authorization reflects role assignments made after the user's token was
// issued.
func (s *PermissionService) GetRolesForUser(ctx context.Context, userID uint) ([]uint, error) {
	roles, err := s.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	roleIDs := make([]uint, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}
	return roleIDs, nil
}

// RebuildRolePolicies 重建指定角色的 Casbin in-memory 策略
// 先删除该角色的所有现有策略，再批量添加新策略。
// 适配器是只读的（策略由统一权限表管理），这里临时关闭 auto-save
// 使变更仅在内存生效。
func (s *PermissionService) RebuildRolePolicies(ctx context.Context, roleID uint, permissions []*sysManagement.PermissionModel) error {
	roleIDStr := strconv.Itoa(int(roleID))

	s.enforcer.EnableAutoSave(false)
	defer s.enforcer.EnableAutoSave(true)

	_, err := s.enforcer.RemoveFilteredPolicy(0, roleIDStr)
	if err != nil {
		return fmt.Errorf("remove old policies failed: %w", err)
	}

	if len(permissions) == 0 {
		return nil
	}

	newPolicies := make([][]string, 0, len(permissions))
	for _, perm := range permissions {
		effect := perm.Effect
		if effect == "" {
			effect = sysManagement.EffectAllow
		}
		newPolicies = append(newPolicies, []string{
			roleIDStr,
			perm.Resource,
			string(perm.Action),
			string(effect),
		})
	}

	_, err = s.enforcer.AddPolicies(newPolicies)
	if err != nil {
		return fmt.Errorf("add new policies failed: %w", err)
	}

	return nil
}

func (s *PermissionService) ReloadPolicy(ctx context.Context) error {
	return s.enforcer.LoadPolicy()
}

// LintPolicies analyzes all role-permission assignments and reports policy
// quality problems (dead allows, redundant rules, unassigned permissions).
func (s *PermissionService) LintPolicies(ctx context.Context) ([]PolicyIssue, error) {
	allPerms, err := s.permRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	roles, err := s.roleRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	assignments := make(map[uint][]*sysManagement.PermissionModel, len(roles))
	for _, role := range roles {
		perms, err := s.permRepo.FindByRoleID(ctx, role.ID)
		if err != nil {
			return nil, err
		}
		assignments[role.ID] = perms
	}

	return LintAssignments(assignments, allPerms), nil
}
