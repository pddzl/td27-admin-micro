package tool

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"td27/rpc/basis/internal/model/tool"
	"td27/rpc/basis/internal/model/common"
)

// FileRepository defines interface for file data operations
type FileRepository interface {
	FindOne(ctx context.Context, id uint) (*tool.FileModel, error)
	Create(ctx context.Context, file *tool.FileModel) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, page *common.PageInfo) ([]*tool.FileModel, int64, error)
}

type fileRepository struct {
	db *gorm.DB
}

// NewFileRepository creates a new file repository instance
func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) FindOne(ctx context.Context, id uint) (*tool.FileModel, error) {
	var file tool.FileModel
	if err := r.db.WithContext(ctx).First(&file, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) Create(ctx context.Context, file *tool.FileModel) error {
	return r.db.WithContext(ctx).Create(file).Error
}

func (r *fileRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&tool.FileModel{}, id).Error
}

func (r *fileRepository) List(ctx context.Context, page *common.PageInfo) ([]*tool.FileModel, int64, error) {
	var files []*tool.FileModel
	var total int64

	query := r.db.WithContext(ctx).Model(&tool.FileModel{})
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	if err := query.Order("created_at desc").Offset(offset).Limit(page.PageSize).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}
