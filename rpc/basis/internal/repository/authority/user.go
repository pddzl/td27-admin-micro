package authority

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"td27/rpc/basis/internal/model/authority"
	"td27/rpc/basis/internal/model/common"
)

// UserRepository defines interface for user data operations
type UserRepository interface {
	FindOne(ctx context.Context, id uint) (*authority.UserModel, error)
	FindOneByUsername(ctx context.Context, username string) (*authority.UserModel, error)
	FindOneWithRoles(ctx context.Context, id uint) (*authority.UserModel, error)
	List(ctx context.Context, page *common.PageInfo, deptID *uint) ([]*authority.UserModel, int64, error)
	Create(ctx context.Context, user *authority.UserModel) error
	Update(ctx context.Context, user *authority.UserModel) error
	Delete(ctx context.Context, id uint) error
	UpdatePassword(ctx context.Context, id uint, newPasswordHash string) error
	ToggleActive(ctx context.Context, id uint, active bool) error
	CountByRoleID(ctx context.Context, roleID uint) (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindOne(ctx context.Context, id uint) (*authority.UserModel, error) {
	var user authority.UserModel
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindOneByUsername(ctx context.Context, username string) (*authority.UserModel, error) {
	var user authority.UserModel
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindOneWithRoles(ctx context.Context, id uint) (*authority.UserModel, error) {
	var user authority.UserModel
	if err := r.db.WithContext(ctx).Preload("Roles").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) List(ctx context.Context, page *common.PageInfo, deptID *uint) ([]*authority.UserModel, int64, error) {
	var users []*authority.UserModel
	var total int64

	query := r.db.WithContext(ctx).Model(&authority.UserModel{})
	
	if deptID != nil {
		query = query.Where("dept_id = ?", *deptID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	if err := query.Offset(offset).Limit(page.PageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Create(ctx context.Context, user *authority.UserModel) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *authority.UserModel) error {
	return r.db.WithContext(ctx).Model(user).Updates(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&authority.UserModel{}, id).Error
}

func (r *userRepository) UpdatePassword(ctx context.Context, id uint, newPasswordHash string) error {
	return r.db.WithContext(ctx).Model(&authority.UserModel{}).Where("id = ?", id).Update("password", newPasswordHash).Error
}

func (r *userRepository) ToggleActive(ctx context.Context, id uint, active bool) error {
	return r.db.WithContext(ctx).Model(&authority.UserModel{}).Where("id = ?", id).Update("active", active).Error
}

func (r *userRepository) CountByRoleID(ctx context.Context, roleID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&authority.UserRole{}).Where("role_id = ?", roleID).Count(&count).Error
	return count, err
}
