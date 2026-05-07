package tool

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"td27/rpc/basis/internal/model/tool"
	"td27/rpc/basis/internal/model/common"
)

// ServiceTokenRepository defines interface for service token data operations
type ServiceTokenRepository interface {
	FindOne(ctx context.Context, id uint) (*tool.ServiceToken, error)
	FindByTokenHash(ctx context.Context, tokenHash string) (*tool.ServiceToken, error)
	FindAllValid(ctx context.Context) ([]*tool.ServiceToken, error)
	List(ctx context.Context, page *common.PageInfo) ([]*tool.ServiceToken, int64, error)
	Create(ctx context.Context, token *tool.ServiceToken) error
	Update(ctx context.Context, token *tool.ServiceToken) error
	ToggleStatus(ctx context.Context, id uint, status bool) error
	Delete(ctx context.Context, id uint) error
	AssignPermissions(ctx context.Context, tokenID uint, permissionIDs []uint) error
	GetPermissions(ctx context.Context, tokenID uint) ([]uint, error)
}

type serviceTokenRepository struct {
	db *gorm.DB
}

// NewServiceTokenRepository creates a new service token repository instance
func NewServiceTokenRepository(db *gorm.DB) ServiceTokenRepository {
	return &serviceTokenRepository{db: db}
}

func (r *serviceTokenRepository) FindOne(ctx context.Context, id uint) (*tool.ServiceToken, error) {
	var token tool.ServiceToken
	if err := r.db.WithContext(ctx).First(&token, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *serviceTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*tool.ServiceToken, error) {
	var token tool.ServiceToken
	err := r.db.WithContext(ctx).Where("token_hash = ? AND status = ?", tokenHash, true).
		Where("expires_at IS NULL OR expires_at > ?", time.Now().Unix()).
		First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *serviceTokenRepository) FindAllValid(ctx context.Context) ([]*tool.ServiceToken, error) {
	var tokens []*tool.ServiceToken
	err := r.db.WithContext(ctx).Where("status = ?", true).
		Where("expires_at IS NULL OR expires_at > ?", time.Now().Unix()).
		Find(&tokens).Error
	return tokens, err
}

func (r *serviceTokenRepository) List(ctx context.Context, page *common.PageInfo) ([]*tool.ServiceToken, int64, error) {
	var tokens []*tool.ServiceToken
	var total int64

	query := r.db.WithContext(ctx).Model(&tool.ServiceToken{})
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	if err := query.Order("created_at desc").Offset(offset).Limit(page.PageSize).Find(&tokens).Error; err != nil {
		return nil, 0, err
	}

	return tokens, total, nil
}

func (r *serviceTokenRepository) Create(ctx context.Context, token *tool.ServiceToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *serviceTokenRepository) Update(ctx context.Context, token *tool.ServiceToken) error {
	return r.db.WithContext(ctx).Model(token).Updates(token).Error
}

func (r *serviceTokenRepository) ToggleStatus(ctx context.Context, id uint, status bool) error {
	return r.db.WithContext(ctx).Model(&tool.ServiceToken{}).Where("id = ?", id).Update("status", status).Error
}

func (r *serviceTokenRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("token_id = ?", id).Delete(&tool.TokenPermission{}).Error; err != nil {
			return err
		}
		return tx.Delete(&tool.ServiceToken{}, id).Error
	})
}

func (r *serviceTokenRepository) AssignPermissions(ctx context.Context, tokenID uint, permissionIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("token_id = ?", tokenID).Delete(&tool.TokenPermission{}).Error; err != nil {
			return err
		}

		for _, permID := range permissionIDs {
			tp := &tool.TokenPermission{
				TokenID:      tokenID,
				PermissionID: permID,
			}
			if err := tx.Create(tp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *serviceTokenRepository) GetPermissions(ctx context.Context, tokenID uint) ([]uint, error) {
	var permissions []uint
	err := r.db.WithContext(ctx).Model(&tool.TokenPermission{}).
		Where("token_id = ?", tokenID).
		Pluck("permission_id", &permissions).Error
	return permissions, err
}
