package sysTool

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysTool"
)

// ServiceTokenRepository defines interface for service token data operations
type ServiceTokenRepository interface {
	FindOne(ctx context.Context, id uint) (*sysTool.ServiceToken, error)
	FindByTokenHash(ctx context.Context, tokenHash string) (*sysTool.ServiceToken, error)
	FindAllValid(ctx context.Context) ([]*sysTool.ServiceToken, error)
	List(ctx context.Context, page *common.PageInfo) ([]*sysTool.ServiceToken, int64, error)
	Create(ctx context.Context, token *sysTool.ServiceToken) error
	Update(ctx context.Context, token *sysTool.ServiceToken) error
	ToggleStatus(ctx context.Context, id uint, status bool) error
	Delete(ctx context.Context, id uint) error
	AssignPermissions(ctx context.Context, tokenID uint, permissionIDs []uint) error
	GetPermissions(ctx context.Context, tokenID uint) ([]uint, error)
}

type serviceTokenRepository struct {
	db *sqlx.DB
}

// NewServiceTokenRepository creates a new service token repository instance
func NewServiceTokenRepository(db *sqlx.DB) ServiceTokenRepository {
	return &serviceTokenRepository{db: db}
}

const serviceTokenColumns = `id, COALESCE(created_at, NOW()) as created_at, COALESCE(updated_at, NOW()) as updated_at, deleted_at, name, token_hash, status, expires_at`

func (r *serviceTokenRepository) FindOne(ctx context.Context, id uint) (*sysTool.ServiceToken, error) {
	var token sysTool.ServiceToken
	err := sqlx.GetContext(ctx, r.db, &token,
		"SELECT "+serviceTokenColumns+" FROM sys_tool_service_token WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *serviceTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*sysTool.ServiceToken, error) {
	var token sysTool.ServiceToken
	err := sqlx.GetContext(ctx, r.db, &token,
		"SELECT "+serviceTokenColumns+" FROM sys_tool_service_token WHERE token_hash=$1 AND status=$2 AND (expires_at IS NULL OR expires_at > $3) AND deleted_at IS NULL",
		tokenHash, true, time.Now().Unix())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *serviceTokenRepository) FindAllValid(ctx context.Context) ([]*sysTool.ServiceToken, error) {
	var tokens []*sysTool.ServiceToken
	err := sqlx.SelectContext(ctx, r.db, &tokens,
		"SELECT "+serviceTokenColumns+" FROM sys_tool_service_token WHERE status=$1 AND (expires_at IS NULL OR expires_at > $2) AND deleted_at IS NULL",
		true, time.Now().Unix())
	return tokens, err
}

func (r *serviceTokenRepository) List(ctx context.Context, page *common.PageInfo) ([]*sysTool.ServiceToken, int64, error) {
	var total int64
	err := sqlx.GetContext(ctx, r.db, &total,
		"SELECT COUNT(*) FROM sys_tool_service_token WHERE deleted_at IS NULL")
	if err != nil {
		return nil, 0, err
	}

	offset := (page.Page - 1) * page.PageSize
	var tokens []*sysTool.ServiceToken
	err = sqlx.SelectContext(ctx, r.db, &tokens,
		"SELECT "+serviceTokenColumns+" FROM sys_tool_service_token WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		page.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return tokens, total, nil
}

func (r *serviceTokenRepository) Create(ctx context.Context, token *sysTool.ServiceToken) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO sys_tool_service_token (name, token_hash, status, expires_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		token.Name, token.TokenHash, token.Status, token.ExpiresAt, token.CreatedAt, token.UpdatedAt)
	return err
}

func (r *serviceTokenRepository) Update(ctx context.Context, token *sysTool.ServiceToken) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sys_tool_service_token SET name=$1, token_hash=$2, status=$3, expires_at=$4, updated_at=NOW() WHERE id=$5 AND deleted_at IS NULL",
		token.Name, token.TokenHash, token.Status, token.ExpiresAt, token.ID)
	return err
}

func (r *serviceTokenRepository) ToggleStatus(ctx context.Context, id uint, status bool) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sys_tool_service_token SET status=$1, updated_at=NOW() WHERE id=$2 AND deleted_at IS NULL",
		status, id)
	return err
}

func (r *serviceTokenRepository) Delete(ctx context.Context, id uint) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, "DELETE FROM sys_tool_token_permission WHERE token_id=$1", id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE sys_tool_service_token SET deleted_at=NOW() WHERE id=$1", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *serviceTokenRepository) AssignPermissions(ctx context.Context, tokenID uint, permissionIDs []uint) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, "DELETE FROM sys_tool_token_permission WHERE token_id=$1", tokenID)
	if err != nil {
		return err
	}

	for _, permID := range permissionIDs {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO sys_tool_token_permission (token_id, permission_id) VALUES ($1, $2)",
			tokenID, permID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *serviceTokenRepository) GetPermissions(ctx context.Context, tokenID uint) ([]uint, error) {
	var permissions []uint
	err := sqlx.SelectContext(ctx, r.db, &permissions,
		"SELECT permission_id FROM sys_tool_token_permission WHERE token_id=$1", tokenID)
	return permissions, err
}
