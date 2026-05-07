package authority

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"td27/rpc/basis/internal/model/authority"
)

// MenuRepository defines interface for menu data operations
type MenuRepository interface {
	FindOne(ctx context.Context, id uint) (*authority.MenuModel, error)
	FindAll(ctx context.Context) ([]*authority.MenuModel, error)
	FindByParentID(ctx context.Context, parentID uint) ([]*authority.MenuModel, error)
	FindByRoleIDs(ctx context.Context, roleIDs []uint) ([]*authority.MenuModel, error)
	Create(ctx context.Context, menu *authority.MenuModel) error
	Update(ctx context.Context, menu *authority.MenuModel) error
	Delete(ctx context.Context, id uint) error
}

type menuRepository struct {
	db *gorm.DB
}

// NewMenuRepository creates a new menu repository instance
func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) FindOne(ctx context.Context, id uint) (*authority.MenuModel, error) {
	var menu authority.MenuModel
	if err := r.db.WithContext(ctx).First(&menu, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) FindAll(ctx context.Context) ([]*authority.MenuModel, error) {
	var menus []*authority.MenuModel
	if err := r.db.WithContext(ctx).Order("sort asc").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *menuRepository) FindByParentID(ctx context.Context, parentID uint) ([]*authority.MenuModel, error) {
	var menus []*authority.MenuModel
	if err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Order("sort asc").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *menuRepository) FindByRoleIDs(ctx context.Context, roleIDs []uint) ([]*authority.MenuModel, error) {
	var menus []*authority.MenuModel
	err := r.db.WithContext(ctx).Distinct("menu_model.*").
		Joins("JOIN sys_management_permission p ON p.domain_id = menu_model.id AND p.domain = 'menu'").
		Joins("JOIN sys_management_role_permissions rp ON rp.permission_id = p.id").
		Where("rp.role_id IN ?", roleIDs).
		Order("menu_model.sort asc").
		Find(&menus).Error
	return menus, err
}

func (r *menuRepository) Create(ctx context.Context, menu *authority.MenuModel) error {
	return r.db.WithContext(ctx).Create(menu).Error
}

func (r *menuRepository) Update(ctx context.Context, menu *authority.MenuModel) error {
	return r.db.WithContext(ctx).Model(menu).Updates(menu).Error
}

func (r *menuRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("parent_id = ?", id).Delete(&authority.MenuModel{}).Error; err != nil {
			return err
		}
		return tx.Delete(&authority.MenuModel{}, id).Error
	})
}
