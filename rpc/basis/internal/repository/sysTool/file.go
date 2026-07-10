package sysTool

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysTool"
)

// FileRepository defines interface for file data operations
type FileRepository interface {
	FindOne(ctx context.Context, id uint) (*sysTool.FileModel, error)
	Create(ctx context.Context, file *sysTool.FileModel) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, page *common.PageInfo) ([]*sysTool.FileModel, int64, error)
}

type fileRepository struct {
	db *sqlx.DB
}

// NewFileRepository creates a new file repository instance
func NewFileRepository(db *sqlx.DB) FileRepository {
	return &fileRepository{db: db}
}

const fileColumns = `id, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at, deleted_at, file_name, full_path, mime`

func (r *fileRepository) FindOne(ctx context.Context, id uint) (*sysTool.FileModel, error) {
	var file sysTool.FileModel
	err := sqlx.GetContext(ctx, r.db, &file,
		"SELECT "+fileColumns+" FROM sys_tool_file WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) Create(ctx context.Context, file *sysTool.FileModel) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO sys_tool_file (file_name, full_path, mime, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		file.FileName, file.FullPath, file.Mime, file.CreatedAt, file.UpdatedAt)
	return err
}

func (r *fileRepository) Delete(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sys_tool_file SET deleted_at=NOW() WHERE id=$1", id)
	return err
}

func (r *fileRepository) List(ctx context.Context, page *common.PageInfo) ([]*sysTool.FileModel, int64, error) {
	var total int64
	err := sqlx.GetContext(ctx, r.db, &total,
		"SELECT COUNT(*) FROM sys_tool_file WHERE deleted_at IS NULL")
	if err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	var files []*sysTool.FileModel
	err = sqlx.SelectContext(ctx, r.db, &files,
		"SELECT "+fileColumns+" FROM sys_tool_file WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		page.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return files, total, nil
}
