package sysManagement

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
)

// DictDetailRepository defines interface for dictionary detail data operations
type DictDetailRepository interface {
	FindOne(ctx context.Context, id uint) (*sysManagement.DictDetailModel, error)
	FindByDictID(ctx context.Context, dictID uint) ([]*sysManagement.DictDetailModel, error)
	FindByDictENName(ctx context.Context, dictENName string) ([]*sysManagement.DictDetailModel, error)
	List(ctx context.Context, page *common.PageInfo, dictID *uint) ([]*sysManagement.DictDetailModel, int64, error)
	Create(ctx context.Context, detail *sysManagement.DictDetailModel) error
	Update(ctx context.Context, detail *sysManagement.DictDetailModel) error
	Delete(ctx context.Context, id uint) error
}

type dictDetailRepository struct {
	db *sqlx.DB
}

// NewDictDetailRepository creates a new dictionary detail repository instance
func NewDictDetailRepository(db *sqlx.DB) DictDetailRepository {
	return &dictDetailRepository{db: db}
}

const dictDetailTable = "sys_management_dict_detail"
const dictDetailColumns = `id, created_at, updated_at, deleted_at, label, value, sort, dict_id, parent_id, description`

func (r *dictDetailRepository) FindOne(ctx context.Context, id uint) (*sysManagement.DictDetailModel, error) {
	var detail sysManagement.DictDetailModel
	err := r.db.GetContext(ctx, &detail, "SELECT "+dictDetailColumns+" FROM "+dictDetailTable+" WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &detail, nil
}

func (r *dictDetailRepository) FindByDictID(ctx context.Context, dictID uint) ([]*sysManagement.DictDetailModel, error) {
	var details []*sysManagement.DictDetailModel
	err := r.db.SelectContext(ctx, &details,
		"SELECT "+dictDetailColumns+" FROM "+dictDetailTable+" WHERE dict_id=$1 AND deleted_at IS NULL ORDER BY sort ASC", dictID)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func (r *dictDetailRepository) FindByDictENName(ctx context.Context, dictENName string) ([]*sysManagement.DictDetailModel, error) {
	var details []*sysManagement.DictDetailModel
	err := r.db.SelectContext(ctx, &details,
		`SELECT ddm.id, ddm.created_at, ddm.updated_at, ddm.deleted_at, ddm.label, ddm.value, ddm.sort, ddm.dict_id, ddm.parent_id, ddm.description FROM sys_management_dict_detail ddm
		 JOIN sys_management_dict d ON d.id = ddm.dict_id
		 WHERE d.en_name = $1 AND ddm.deleted_at IS NULL AND d.deleted_at IS NULL
		 ORDER BY ddm.sort ASC`, dictENName)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func (r *dictDetailRepository) List(ctx context.Context, page *common.PageInfo, dictID *uint) ([]*sysManagement.DictDetailModel, int64, error) {
	if dictID == nil {
		return nil, 0, nil
	}

	// Count only root nodes
	var total int64
	if err := r.db.GetContext(ctx, &total,
		"SELECT COUNT(*) FROM "+dictDetailTable+" WHERE dict_id=$1 AND parent_id IS NULL AND deleted_at IS NULL", *dictID); err != nil {
		return nil, 0, err
	}

	// Load paginated roots
	offset := (page.Page - 1) * page.PageSize
	var roots []*sysManagement.DictDetailModel
	if err := r.db.SelectContext(ctx, &roots,
		"SELECT "+dictDetailColumns+" FROM "+dictDetailTable+" WHERE dict_id=$1 AND parent_id IS NULL AND deleted_at IS NULL ORDER BY sort ASC LIMIT $2 OFFSET $3",
		*dictID, page.PageSize, offset); err != nil {
		return nil, 0, err
	}

	if len(roots) == 0 {
		return []*sysManagement.DictDetailModel{}, total, nil
	}

	// Load all descendants (single query, non-paginated)
	var all []*sysManagement.DictDetailModel
	if err := r.db.SelectContext(ctx, &all,
		"SELECT "+dictDetailColumns+" FROM "+dictDetailTable+" WHERE dict_id=$1 AND parent_id IS NOT NULL AND deleted_at IS NULL ORDER BY sort ASC", *dictID); err != nil {
		return nil, 0, err
	}

	// Build node map
	nodeMap := make(map[uint]*sysManagement.DictDetailModel, len(all)+len(roots))
	for _, r := range roots {
		r.Children = []*sysManagement.DictDetailModel{}
		nodeMap[r.ID] = r
	}
	for _, n := range all {
		n.Children = []*sysManagement.DictDetailModel{}
		nodeMap[n.ID] = n
	}

	// Attach children to parents
	for _, n := range all {
		if n.ParentID != nil {
			if parent, ok := nodeMap[uint(*n.ParentID)]; ok {
				parent.Children = append(parent.Children, n)
			}
		}
	}

	return roots, total, nil
}

func (r *dictDetailRepository) Create(ctx context.Context, detail *sysManagement.DictDetailModel) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO sys_management_dict_detail (created_at, updated_at, label, value, sort, dict_id, parent_id, description)
		 VALUES (NOW(), NOW(), $1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`,
		detail.Label, detail.Value, detail.Sort, detail.DictModelID, detail.ParentID, detail.Description).Scan(&detail.ID, &detail.CreatedAt, &detail.UpdatedAt)
	return err
}

func (r *dictDetailRepository) Update(ctx context.Context, detail *sysManagement.DictDetailModel) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE sys_management_dict_detail
		 SET label=$1, value=$2, sort=$3, dict_id=$4, parent_id=$5, description=$6, updated_at=NOW()
		 WHERE id=$7 AND deleted_at IS NULL`,
		detail.Label, detail.Value, detail.Sort, detail.DictModelID, detail.ParentID, detail.Description, detail.ID)
	return err
}

func (r *dictDetailRepository) Delete(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM sys_management_dict_detail WHERE id=$1", id)
	return err
}
