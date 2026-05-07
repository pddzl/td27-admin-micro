package authority

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"td27/rpc/basis/internal/model/authority"
	"td27/rpc/basis/internal/model/common"
)

// DeptRepository defines interface for department data operations
type DeptRepository interface {
	FindOne(ctx context.Context, id uint) (*authority.DeptModel, error)
	FindAll(ctx context.Context) ([]*authority.DeptModel, error)
	FindByParentID(ctx context.Context, parentID uint) ([]*authority.DeptModel, error)
	FindDescendants(ctx context.Context, deptID uint) ([]*authority.DeptModel, error)
	FindAncestors(ctx context.Context, path string) ([]*authority.DeptModel, error)
	List(ctx context.Context, page *common.PageInfo) ([]*authority.DeptModel, int64, error)
	Create(ctx context.Context, dept *authority.DeptModel) error
	Update(ctx context.Context, dept *authority.DeptModel) error
	Delete(ctx context.Context, id uint) error
}

type deptRepository struct {
	db *gorm.DB
}

// NewDeptRepository creates a new department repository instance
func NewDeptRepository(db *gorm.DB) DeptRepository {
	return &deptRepository{db: db}
}

func (r *deptRepository) FindOne(ctx context.Context, id uint) (*authority.DeptModel, error) {
	var dept authority.DeptModel
	if err := r.db.WithContext(ctx).First(&dept, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dept, nil
}

func (r *deptRepository) FindAll(ctx context.Context) ([]*authority.DeptModel, error) {
	var depts []*authority.DeptModel
	if err := r.db.WithContext(ctx).Order("level asc, sort asc").Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}

func (r *deptRepository) FindByParentID(ctx context.Context, parentID uint) ([]*authority.DeptModel, error) {
	var depts []*authority.DeptModel
	if err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Order("sort asc").Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}

func (r *deptRepository) FindDescendants(ctx context.Context, deptID uint) ([]*authority.DeptModel, error) {
	dept, err := r.FindOne(ctx, deptID)
	if err != nil || dept == nil {
		return nil, err
	}

	pathPrefix := fmt.Sprintf("%s%d/", dept.Path, deptID)
	var depts []*authority.DeptModel
	if err := r.db.WithContext(ctx).Where("path LIKE ?", pathPrefix+"%").Order("level asc, sort asc").Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}

func (r *deptRepository) FindAncestors(ctx context.Context, path string) ([]*authority.DeptModel, error) {
	var depts []*authority.DeptModel
	if err := r.db.WithContext(ctx).Where("path IN ?", splitPathToAncestors(path)).Order("level asc").Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}

func splitPathToAncestors(path string) []string {
	// Path format: "/1/2/3/" -> returns ["/1/", "/1/2/"]
	var paths []string
	current := ""
	for _, c := range path {
		if c == '/' && current != "" {
			paths = append(paths, current+"/")
		}
		current += string(c)
	}
	return paths
}

func (r *deptRepository) List(ctx context.Context, page *common.PageInfo) ([]*authority.DeptModel, int64, error) {
	var depts []*authority.DeptModel
	var total int64

	query := r.db.WithContext(ctx).Model(&authority.DeptModel{})
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	if err := query.Order("level asc, sort asc").Offset(offset).Limit(page.PageSize).Find(&depts).Error; err != nil {
		return nil, 0, err
	}

	return depts, total, nil
}

func (r *deptRepository) Create(ctx context.Context, dept *authority.DeptModel) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if dept.ParentID != 0 {
			parent, err := r.FindOne(ctx, dept.ParentID)
			if err != nil || parent == nil {
				return fmt.Errorf("parent department not found")
			}
			dept.Path = fmt.Sprintf("%s%d/", parent.Path, parent.ID)
			dept.Level = parent.Level + 1
		} else {
			dept.Path = "/"
			dept.Level = 1
		}
		return tx.Create(dept).Error
	})
}

func (r *deptRepository) Update(ctx context.Context, dept *authority.DeptModel) error {
	return r.db.WithContext(ctx).Model(dept).Updates(dept).Error
}

func (r *deptRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		descendants, err := r.FindDescendants(ctx, id)
		if err != nil {
			return err
		}

		for _, d := range descendants {
			if err := tx.Delete(d).Error; err != nil {
				return err
			}
		}

		return tx.Delete(&authority.DeptModel{}, id).Error
	})
}
