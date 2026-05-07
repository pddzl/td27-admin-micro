package authority

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"td27/rpc/basis/internal/model/authority"
)

// PermissionRepository defines interface for permission data operations
type PermissionRepository interface {
	FindOne(ctx context.Context, id uint) (*authority.PermissionModel, error)
	FindAll(ctx context.Context) ([]*authority.PermissionModel, error)
	FindByDomain(ctx context.Context, domain authority.PermissionDomain) ([]*authority.PermissionModel, error)
	FindByRoleID(ctx context.Context, roleID uint) ([]*authority.PermissionModel, error)
	FindByResourceAndAction(ctx context.Context, resource string, action authority.Action) (*authority.PermissionModel, error)
	Create(ctx context.Context, permission *authority.PermissionModel) error
	Update(ctx context.Context, permission *authority.PermissionModel) error
	Delete(ctx context.Context, id uint) error
}

type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository creates a new permission repository instance
func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) FindOne(ctx context.Context, id uint) (*authority.PermissionModel, error) {
	var perm authority.PermissionModel
	if err := r.db.WithContext(ctx).First(&perm, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &perm, nil
}

func (r *permissionRepository) FindAll(ctx context.Context) ([]*authority.PermissionModel, error) {
	var perms []*authority.PermissionModel
	if err := r.db.WithContext(ctx).Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

func (r *permissionRepository) FindByDomain(ctx context.Context, domain authority.PermissionDomain) ([]*authority.PermissionModel, error) {
	var perms []*authority.PermissionModel
	if err := r.db.WithContext(ctx).Where("domain = ?", domain).Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

func (r *permissionRepository) FindByRoleID(ctx context.Context, roleID uint) ([]*authority.PermissionModel, error) {
	var perms []*authority.PermissionModel
	err := r.db.WithContext(ctx).Joins("JOIN sys_management_role_permissions rp ON rp.permission_id = permission_model.id").
		Where("rp.role_id = ?", roleID).
		Find(&perms).Error
	return perms, err
}

func (r *permissionRepository) FindByResourceAndAction(ctx context.Context, resource string, action authority.Action) (*authority.PermissionModel, error) {
	var perm authority.PermissionModel
	err := r.db.WithContext(ctx).Where("resource = ? AND action = ?", resource, action).First(&perm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &perm, nil
}

func (r *permissionRepository) Create(ctx context.Context, permission *authority.PermissionModel) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

func (r *permissionRepository) Update(ctx context.Context, permission *authority.PermissionModel) error {
	return r.db.WithContext(ctx).Model(permission).Updates(permission).Error
}

func (r *permissionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("permission_id = ?", id).Delete(&authority.RolePermissionModel{}).Error; err != nil {
			return err
		}
		return tx.Delete(&authority.PermissionModel{}, id).Error
	})
}
