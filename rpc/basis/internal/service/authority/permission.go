package authority

import (
	"context"
	"errors"
	"fmt"

	"github.com/casbin/casbin/v2"

	"td27/rpc/basis/internal/model/authority"
	"td27/rpc/basis/internal/model/common"
	repoAuthority "td27/rpc/basis/internal/repository/authority"
)

type PermissionService struct {
	permRepo   repoAuthority.PermissionRepository
	roleRepo   repoAuthority.RoleRepository
	enforcer   *casbin.SyncedCachedEnforcer
}

func NewPermissionService(
	permRepo repoAuthority.PermissionRepository,
	roleRepo repoAuthority.RoleRepository,
	enforcer *casbin.SyncedCachedEnforcer,
) *PermissionService {
	return &PermissionService{
		permRepo:   permRepo,
		roleRepo:   roleRepo,
		enforcer:   enforcer,
	}
}

func (s *PermissionService) GetByID(ctx context.Context, id uint) (*authority.PermissionModel, error) {
	return s.permRepo.FindOne(ctx, id)
}

func (s *PermissionService) GetAll(ctx context.Context) ([]*authority.PermissionModel, error) {
	return s.permRepo.FindAll(ctx)
}

func (s *PermissionService) GetByDomain(ctx context.Context, domain authority.PermissionDomain) ([]*authority.PermissionModel, error) {
	return s.permRepo.FindByDomain(ctx, domain)
}

func (s *PermissionService) GetByRoleID(ctx context.Context, roleID uint) ([]*authority.PermissionModel, error) {
	return s.permRepo.FindByRoleID(ctx, roleID)
}

func (s *PermissionService) GetByResourceAndAction(ctx context.Context, resource string, action authority.Action) (*authority.PermissionModel, error) {
	return s.permRepo.FindByResourceAndAction(ctx, resource, action)
}

func (s *PermissionService) List(ctx context.Context, page *common.PageInfo, domain *authority.PermissionDomain) ([]*authority.PermissionModel, int64, error) {
	perms, err := s.permRepo.FindAll(ctx)
	if err != nil {
		return nil, 0, err
	}
	return perms, 0, nil
}

func (s *PermissionService) Create(ctx context.Context, perm *authority.PermissionModel) error {
	existing, err := s.permRepo.FindByResourceAndAction(ctx, perm.Resource, perm.Action)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("permission with same resource and action already exists")
	}

	return s.permRepo.Create(ctx, perm)
}

func (s *PermissionService) Update(ctx context.Context, perm *authority.PermissionModel) error {
	return s.permRepo.Update(ctx, perm)
}

func (s *PermissionService) Delete(ctx context.Context, id uint) error {
	return s.permRepo.Delete(ctx, id)
}

func (s *PermissionService) CheckPermission(ctx context.Context, roleIDs []uint, resource string, action authority.Action) (bool, error) {
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

func (s *PermissionService) ReloadPolicy(ctx context.Context) error {
	return s.enforcer.LoadPolicy()
}
