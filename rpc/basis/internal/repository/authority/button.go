package authority

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"td27/rpc/basis/internal/model/authority"
	"td27/rpc/basis/internal/model/common"
)

// ButtonRepository defines interface for button permission data operations
type ButtonRepository interface {
	FindOne(ctx context.Context, id uint) (*authority.ButtonModel, error)
	FindByCode(ctx context.Context, code string) (*authority.ButtonModel, error)
	FindByPagePath(ctx context.Context, pagePath string) ([]*authority.ButtonModel, error)
	FindByRoleIDs(ctx context.Context, roleIDs []uint) ([]*authority.ButtonModel, error)
	List(ctx context.Context, page *common.PageInfo, pagePath *string) ([]*authority.ButtonModel, int64, error)
	Create(ctx context.Context, button *authority.ButtonModel) error
	Update(ctx context.Context, button *authority.ButtonModel) error
	Delete(ctx context.Context, id uint) error
}

type buttonRepository struct {
	db *gorm.DB
}

// NewButtonRepository creates a new button repository instance
func NewButtonRepository(db *gorm.DB) ButtonRepository {
	return &buttonRepository{db: db}
}

func (r *buttonRepository) FindOne(ctx context.Context, id uint) (*authority.ButtonModel, error) {
	var button authority.ButtonModel
	if err := r.db.WithContext(ctx).First(&button, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &button, nil
}

func (r *buttonRepository) FindByCode(ctx context.Context, code string) (*authority.ButtonModel, error) {
	var button authority.ButtonModel
	if err := r.db.WithContext(ctx).Where("button_code = ?", code).First(&button).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &button, nil
}

func (r *buttonRepository) FindByPagePath(ctx context.Context, pagePath string) ([]*authority.ButtonModel, error) {
	var buttons []*authority.ButtonModel
	if err := r.db.WithContext(ctx).Where("page_path = ?", pagePath).Find(&buttons).Error; err != nil {
		return nil, err
	}
	return buttons, nil
}

func (r *buttonRepository) FindByRoleIDs(ctx context.Context, roleIDs []uint) ([]*authority.ButtonModel, error) {
	var buttons []*authority.ButtonModel
	err := r.db.WithContext(ctx).Distinct("button_model.*").
		Joins("JOIN sys_management_permission p ON p.domain_id = button_model.id AND p.domain = 'button'").
		Joins("JOIN sys_management_role_permissions rp ON rp.permission_id = p.id").
		Where("rp.role_id IN ?", roleIDs).
		Find(&buttons).Error
	return buttons, err
}

func (r *buttonRepository) List(ctx context.Context, page *common.PageInfo, pagePath *string) ([]*authority.ButtonModel, int64, error) {
	var buttons []*authority.ButtonModel
	var total int64

	query := r.db.WithContext(ctx).Model(&authority.ButtonModel{})
	
	if pagePath != nil {
		query = query.Where("page_path = ?", *pagePath)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	if err := query.Offset(offset).Limit(page.PageSize).Find(&buttons).Error; err != nil {
		return nil, 0, err
	}

	return buttons, total, nil
}

func (r *buttonRepository) Create(ctx context.Context, button *authority.ButtonModel) error {
	return r.db.WithContext(ctx).Create(button).Error
}

func (r *buttonRepository) Update(ctx context.Context, button *authority.ButtonModel) error {
	return r.db.WithContext(ctx).Model(button).Updates(button).Error
}

func (r *buttonRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&authority.ButtonModel{}, id).Error
}
