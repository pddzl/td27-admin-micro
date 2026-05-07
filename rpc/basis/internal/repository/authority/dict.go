package authority

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"td27/rpc/basis/internal/model/authority"
	"td27/rpc/basis/internal/model/common"
)

// DictRepository defines interface for dictionary data operations
type DictRepository interface {
	FindOne(ctx context.Context, id uint) (*authority.DictModel, error)
	FindByENName(ctx context.Context, enName string) (*authority.DictModel, error)
	FindAll(ctx context.Context) ([]*authority.DictModel, error)
	List(ctx context.Context, page *common.PageInfo) ([]*authority.DictModel, int64, error)
	Create(ctx context.Context, dict *authority.DictModel) error
	Update(ctx context.Context, dict *authority.DictModel) error
	Delete(ctx context.Context, id uint) error
}

type dictRepository struct {
	db *gorm.DB
}

// NewDictRepository creates a new dictionary repository instance
func NewDictRepository(db *gorm.DB) DictRepository {
	return &dictRepository{db: db}
}

func (r *dictRepository) FindOne(ctx context.Context, id uint) (*authority.DictModel, error) {
	var dict authority.DictModel
	if err := r.db.WithContext(ctx).Preload("DictDetails").First(&dict, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dict, nil
}

func (r *dictRepository) FindByENName(ctx context.Context, enName string) (*authority.DictModel, error) {
	var dict authority.DictModel
	if err := r.db.WithContext(ctx).Preload("DictDetails").Where("en_name = ?", enName).First(&dict).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dict, nil
}

func (r *dictRepository) FindAll(ctx context.Context) ([]*authority.DictModel, error) {
	var dicts []*authority.DictModel
	if err := r.db.WithContext(ctx).Preload("DictDetails").Find(&dicts).Error; err != nil {
		return nil, err
	}
	return dicts, nil
}

func (r *dictRepository) List(ctx context.Context, page *common.PageInfo) ([]*authority.DictModel, int64, error) {
	var dicts []*authority.DictModel
	var total int64

	query := r.db.WithContext(ctx).Model(&authority.DictModel{})
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	if err := query.Preload("DictDetails").Offset(offset).Limit(page.PageSize).Find(&dicts).Error; err != nil {
		return nil, 0, err
	}

	return dicts, total, nil
}

func (r *dictRepository) Create(ctx context.Context, dict *authority.DictModel) error {
	return r.db.WithContext(ctx).Create(dict).Error
}

func (r *dictRepository) Update(ctx context.Context, dict *authority.DictModel) error {
	return r.db.WithContext(ctx).Model(dict).Updates(dict).Error
}

func (r *dictRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("dict_id = ?", id).Delete(&authority.DictDetailModel{}).Error; err != nil {
			return err
		}
		return tx.Delete(&authority.DictModel{}, id).Error
	})
}
