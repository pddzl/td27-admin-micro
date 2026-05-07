package authority

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"td27/rpc/basis/internal/model/authority"
	"td27/rpc/basis/internal/model/common"
)

// DictDetailRepository defines interface for dictionary detail data operations
type DictDetailRepository interface {
	FindOne(ctx context.Context, id uint) (*authority.DictDetailModel, error)
	FindByDictID(ctx context.Context, dictID uint) ([]*authority.DictDetailModel, error)
	FindByDictENName(ctx context.Context, dictENName string) ([]*authority.DictDetailModel, error)
	List(ctx context.Context, page *common.PageInfo, dictID *uint) ([]*authority.DictDetailModel, int64, error)
	Create(ctx context.Context, detail *authority.DictDetailModel) error
	Update(ctx context.Context, detail *authority.DictDetailModel) error
	Delete(ctx context.Context, id uint) error
}

type dictDetailRepository struct {
	db *gorm.DB
}

// NewDictDetailRepository creates a new dictionary detail repository instance
func NewDictDetailRepository(db *gorm.DB) DictDetailRepository {
	return &dictDetailRepository{db: db}
}

func (r *dictDetailRepository) FindOne(ctx context.Context, id uint) (*authority.DictDetailModel, error) {
	var detail authority.DictDetailModel
	if err := r.db.WithContext(ctx).First(&detail, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &detail, nil
}

func (r *dictDetailRepository) FindByDictID(ctx context.Context, dictID uint) ([]*authority.DictDetailModel, error) {
	var details []*authority.DictDetailModel
	if err := r.db.WithContext(ctx).Where("dict_id = ?", dictID).Order("sort asc").Find(&details).Error; err != nil {
		return nil, err
	}
	return details, nil
}

func (r *dictDetailRepository) FindByDictENName(ctx context.Context, dictENName string) ([]*authority.DictDetailModel, error) {
	var details []*authority.DictDetailModel
	err := r.db.WithContext(ctx).Joins("JOIN sys_management_dict d ON d.id = dict_detail_model.dict_id").
		Where("d.en_name = ?", dictENName).
		Order("dict_detail_model.sort asc").
		Find(&details).Error
	return details, err
}

func (r *dictDetailRepository) List(ctx context.Context, page *common.PageInfo, dictID *uint) ([]*authority.DictDetailModel, int64, error) {
	var details []*authority.DictDetailModel
	var total int64

	query := r.db.WithContext(ctx).Model(&authority.DictDetailModel{})
	
	if dictID != nil {
		query = query.Where("dict_id = ?", *dictID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	if err := query.Order("sort asc").Offset(offset).Limit(page.PageSize).Find(&details).Error; err != nil {
		return nil, 0, err
	}

	return details, total, nil
}

func (r *dictDetailRepository) Create(ctx context.Context, detail *authority.DictDetailModel) error {
	return r.db.WithContext(ctx).Create(detail).Error
}

func (r *dictDetailRepository) Update(ctx context.Context, detail *authority.DictDetailModel) error {
	return r.db.WithContext(ctx).Model(detail).Updates(detail).Error
}

func (r *dictDetailRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&authority.DictDetailModel{}, id).Error
}
