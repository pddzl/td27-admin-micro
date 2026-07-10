package sysManagement

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
)

// DeptRepository defines interface for department data operations
type DeptRepository interface {
	FindOne(ctx context.Context, id uint) (*sysManagement.DeptModel, error)
	FindAll(ctx context.Context) ([]*sysManagement.DeptModel, error)
	FindByParentID(ctx context.Context, parentID uint) ([]*sysManagement.DeptModel, error)
	FindDescendants(ctx context.Context, deptID uint) ([]*sysManagement.DeptModel, error)
	FindAncestors(ctx context.Context, path string) ([]*sysManagement.DeptModel, error)
	List(ctx context.Context, page *common.PageInfo) ([]*sysManagement.DeptModel, int64, error)
	Create(ctx context.Context, dept *sysManagement.DeptModel) error
	Update(ctx context.Context, dept *sysManagement.DeptModel) error
	Delete(ctx context.Context, id uint) error
}

type deptRepository struct {
	db *sqlx.DB
}

// NewDeptRepository creates a new department repository instance
func NewDeptRepository(db *sqlx.DB) DeptRepository {
	return &deptRepository{db: db}
}

const deptTable = "sys_management_dept"
const deptColumns = `id, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at, deleted_at, dept_name, parent_id, path, level, sort, status`

func (r *deptRepository) FindOne(ctx context.Context, id uint) (*sysManagement.DeptModel, error) {
	var dept sysManagement.DeptModel
	if err := r.db.GetContext(ctx, &dept, "SELECT "+deptColumns+" FROM "+deptTable+" WHERE id=$1 AND deleted_at IS NULL", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &dept, nil
}

func (r *deptRepository) FindAll(ctx context.Context) ([]*sysManagement.DeptModel, error) {
	var depts []*sysManagement.DeptModel
	if err := r.db.SelectContext(ctx, &depts, "SELECT "+deptColumns+" FROM "+deptTable+" WHERE deleted_at IS NULL ORDER BY level asc, sort asc"); err != nil {
		return nil, err
	}
	return depts, nil
}

func (r *deptRepository) FindByParentID(ctx context.Context, parentID uint) ([]*sysManagement.DeptModel, error) {
	var depts []*sysManagement.DeptModel
	if err := r.db.SelectContext(ctx, &depts, "SELECT "+deptColumns+" FROM "+deptTable+" WHERE parent_id=$1 AND deleted_at IS NULL ORDER BY sort asc", parentID); err != nil {
		return nil, err
	}
	return depts, nil
}

func (r *deptRepository) FindDescendants(ctx context.Context, deptID uint) ([]*sysManagement.DeptModel, error) {
	dept, err := r.FindOne(ctx, deptID)
	if err != nil || dept == nil {
		return nil, err
	}

	pathPrefix := fmt.Sprintf("%s%d/", dept.Path, deptID)
	var depts []*sysManagement.DeptModel
	if err = r.db.SelectContext(ctx, &depts, "SELECT "+deptColumns+" FROM "+deptTable+" WHERE path LIKE $1 AND deleted_at IS NULL ORDER BY level asc, sort asc", pathPrefix+"%"); err != nil {
		return nil, err
	}
	return depts, nil
}

func (r *deptRepository) FindAncestors(ctx context.Context, path string) ([]*sysManagement.DeptModel, error) {
	ancestorPaths := splitPathToAncestors(path)
	query, args, err := sqlx.In("SELECT "+deptColumns+" FROM "+deptTable+" WHERE path IN (?) AND deleted_at IS NULL ORDER BY level asc", ancestorPaths)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	var depts []*sysManagement.DeptModel
	if err = r.db.SelectContext(ctx, &depts, query, args...); err != nil {
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

func (r *deptRepository) List(ctx context.Context, page *common.PageInfo) ([]*sysManagement.DeptModel, int64, error) {
	var total int64
	var depts []*sysManagement.DeptModel

	if err := r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM "+deptTable+" WHERE deleted_at IS NULL"); err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	if err := r.db.SelectContext(ctx, &depts,
		"SELECT "+deptColumns+" FROM "+deptTable+" WHERE deleted_at IS NULL ORDER BY level asc, sort asc LIMIT $1 OFFSET $2",
		page.PageSize, offset); err != nil {
		return nil, 0, err
	}

	return depts, total, nil
}

func (r *deptRepository) Create(ctx context.Context, dept *sysManagement.DeptModel) error {
	// Parent lookup (outside transaction, same as original GORM pattern)
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

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // no-op if already committed

	err = tx.QueryRowContext(ctx,
		`INSERT INTO `+deptTable+` (dept_name, parent_id, path, level, sort, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		dept.DeptName, dept.ParentID, dept.Path, dept.Level, dept.Sort, dept.Status, dept.CreatedAt, dept.UpdatedAt).Scan(&dept.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *deptRepository) Update(ctx context.Context, dept *sysManagement.DeptModel) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE `+deptTable+` SET dept_name=$1, parent_id=$2, path=$3, level=$4, sort=$5, status=$6, updated_at=$7 WHERE id=$8`,
		dept.DeptName, dept.ParentID, dept.Path, dept.Level, dept.Sort, dept.Status, dept.UpdatedAt, dept.ID)
	return err
}

func (r *deptRepository) Delete(ctx context.Context, id uint) error {
	descendants, err := r.FindDescendants(ctx, id)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // no-op if already committed

	for _, d := range descendants {
		if _, err = tx.ExecContext(ctx, "UPDATE "+deptTable+" SET deleted_at=NOW() WHERE id=$1", d.ID); err != nil {
			return err
		}
	}

	if _, err = tx.ExecContext(ctx, "UPDATE "+deptTable+" SET deleted_at=NOW() WHERE id=$1", id); err != nil {
		return err
	}

	return tx.Commit()
}
