package authority

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"td27/rpc/basis/internal/model/authority"
	"td27/rpc/basis/internal/model/common"
)

// APIRepository defines interface for API data operations
type APIRepository interface {
	FindOne(ctx context.Context, id uint) (*authority.ApiModel, error)
	FindAll(ctx context.Context) ([]*authority.ApiModel, error)
	FindByPathAndMethod(ctx context.Context, path string, method string) (*authority.ApiModel, error)
	FindByGroup(ctx context.Context, groupEN string) ([]*authority.ApiModel, error)
	List(ctx context.Context, page *common.PageInfo, groupEN *string) ([]*authority.ApiModel, int64, error)
	Create(ctx context.Context, api *authority.ApiModel) error
	Update(ctx context.Context, api *authority.ApiModel) error
	Delete(ctx context.Context, id uint) error
}

type apiRepository struct {
	db *gorm.DB
}

// NewAPIRepository creates a new API repository instance
func NewAPIRepository(db *gorm.DB) APIRepository {
	return &apiRepository{db: db}
}

func (r *apiRepository) FindOne(ctx context.Context, id uint) (*authority.ApiModel, error) {
	var api authority.ApiModel
	if err := r.db.WithContext(ctx).First(&api, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &api, nil
}

func (r *apiRepository) FindAll(ctx context.Context) ([]*authority.ApiModel, error) {
	var apis []*authority.ApiModel
	if err := r.db.WithContext(ctx).Find(&apis).Error; err != nil {
		return nil, err
	}
	return apis, nil
}

func (r *apiRepository) FindByPathAndMethod(ctx context.Context, path string, method string) (*authority.ApiModel, error) {
	var api authority.ApiModel
	if err := r.db.WithContext(ctx).Where("path = ? AND method = ?", path, method).First(&api).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &api, nil
}

func (r *apiRepository) FindByGroup(ctx context.Context, groupEN string) ([]*authority.ApiModel, error) {
	var apis []*authority.ApiModel
	if err := r.db.WithContext(ctx).Where("group_en = ?", groupEN).Find(&apis).Error; err != nil {
		return nil, err
	}
	return apis, nil
}

func (r *apiRepository) List(ctx context.Context, page *common.PageInfo, groupEN *string) ([]*authority.ApiModel, int64, error) {
	var apis []*authority.ApiModel
	var total int64

	query := r.db.WithContext(ctx).Model(&authority.ApiModel{})
	
	if groupEN != nil {
		query = query.Where("group_en = ?", *groupEN)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	if err := query.Offset(offset).Limit(page.PageSize).Find(&apis).Error; err != nil {
		return nil, 0, err
	}

	return apis, total, nil
}

func (r *apiRepository) Create(ctx context.Context, api *authority.ApiModel) error {
	return r.db.WithContext(ctx).Create(api).Error
}

func (r *apiRepository) Update(ctx context.Context, api *authority.ApiModel) error {
	return r.db.WithContext(ctx).Model(api).Updates(api).Error
}

func (r *apiRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&authority.ApiModel{}, id).Error
}
