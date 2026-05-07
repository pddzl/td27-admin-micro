package authority

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"td27/rpc/basis/internal/model/authority"
	"td27/rpc/basis/internal/model/common"
)

// RoleRepository defines interface for role data operations
type RoleRepository interface {
	FindOne(ctx context.Context, id uint) (*authority.RoleModel, error)
	FindOneWithChildren(ctx context.Context, id uint) (*authority.RoleModel, error)
	FindAll(ctx context.Context) ([]*authority.RoleModel, error)
	List(ctx context.Context, page *common.PageInfo) ([]*authority.RoleModel, int64, error)
	Create(ctx context.Context, role *authority.RoleModel) error
	Update(ctx context.Context, role *authority.RoleModel) error
	Delete(ctx context.Context, id uint) error
	AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
	GetPermissions(ctx context.Context, roleID uint) ([]uint, error)
	GetUserRoles(ctx context.Context, userID uint) ([]*authority.RoleModel, error)
	AssignUserRoles(ctx context.Context, userID uint, roleIDs []uint) error
}

type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new role repository instance
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindOne(ctx context.Context, id uint) (*authority.RoleModel, error) {
	var role authority.RoleModel
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindOneWithChildren(ctx context.Context, id uint) (*authority.RoleModel, error) {
	var role authority.RoleModel
	if err := r.db.WithContext(ctx).Preload("Children").First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindAll(ctx context.Context) ([]*authority.RoleModel, error) {
	var roles []*authority.RoleModel
	if err := r.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) List(ctx context.Context, page *common.PageInfo) ([]*authority.RoleModel, int64, error) {
	var roles []*authority.RoleModel
	var total int64

	query := r.db.WithContext(ctx).Model(&authority.RoleModel{})
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	if err := query.Offset(offset).Limit(page.PageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *roleRepository) Create(ctx context.Context, role *authority.RoleModel) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *roleRepository) Update(ctx context.Context, role *authority.RoleModel) error {
	return r.db.WithContext(ctx).Model(role).Updates(role).Error
}

func (r *roleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&authority.RolePermissionModel{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", id).Delete(&authority.UserRole{}).Error; err != nil {
			return err
		}
		return tx.Delete(&authority.RoleModel{}, id).Error
	})
}

func (r *roleRepository) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&authority.RolePermissionModel{}).Error; err != nil {
			return err
		}

		for _, permID := range permissionIDs {
			rp := &authority.RolePermissionModel{
				RoleID:       roleID,
				PermissionID: permID,
			}
			if err := tx.Create(rp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *roleRepository) GetPermissions(ctx context.Context, roleID uint) ([]uint, error) {
	var permissions []uint
	err := r.db.WithContext(ctx).Model(&authority.RolePermissionModel{}).
		Where("role_id = ?", roleID).
		Pluck("permission_id", &permissions).Error
	return permissions, err
}

func (r *roleRepository) GetUserRoles(ctx context.Context, userID uint) ([]*authority.RoleModel, error) {
	var roles []*authority.RoleModel
	err := r.db.WithContext(ctx).Joins("JOIN sys_management_user_roles ur ON ur.role_id = role_model.id").
		Where("ur.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *roleRepository) AssignUserRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&authority.UserRole{}).Error; err != nil {
			return err
		}

		for _, roleID := range roleIDs {
			ur := &authority.UserRole{
				UserID: userID,
				RoleID: roleID,
			}
			if err := tx.Create(ur).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
