package sysManagement

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"td27/rpc/basis/internal/model/sysManagement"
)

// RoleElevationRepository defines interface for JIT elevation data operations
type RoleElevationRepository interface {
	FindOne(ctx context.Context, id uint) (*sysManagement.RoleElevationModel, error)
	FindPendingByUserAndRole(ctx context.Context, userID, roleID uint) (*sysManagement.RoleElevationModel, error)
	Create(ctx context.Context, elevation *sysManagement.RoleElevationModel) error
	Update(ctx context.Context, elevation *sysManagement.RoleElevationModel) error
	List(ctx context.Context) ([]*sysManagement.RoleElevationModel, error)
	ListByUser(ctx context.Context, userID uint) ([]*sysManagement.RoleElevationModel, error)
	ActiveRoleIDs(ctx context.Context, userID uint) ([]uint, error)
}

type roleElevationRepository struct {
	db *sqlx.DB
}

func NewRoleElevationRepository(db *sqlx.DB) RoleElevationRepository {
	return &roleElevationRepository{db: db}
}

const elevationColumns = `id, created_at, updated_at, deleted_at, user_id, role_id, status, reason, duration_minutes, expires_at, requested_by, decided_by`

func (r *roleElevationRepository) FindOne(ctx context.Context, id uint) (*sysManagement.RoleElevationModel, error) {
	var e sysManagement.RoleElevationModel
	err := r.db.GetContext(ctx, &e,
		"SELECT "+elevationColumns+" FROM sys_management_role_elevations WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *roleElevationRepository) FindPendingByUserAndRole(ctx context.Context, userID, roleID uint) (*sysManagement.RoleElevationModel, error) {
	var e sysManagement.RoleElevationModel
	err := r.db.GetContext(ctx, &e,
		"SELECT "+elevationColumns+" FROM sys_management_role_elevations WHERE user_id=$1 AND role_id=$2 AND status=$3 AND deleted_at IS NULL LIMIT 1",
		userID, roleID, sysManagement.ElevationPending)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *roleElevationRepository) Create(ctx context.Context, elevation *sysManagement.RoleElevationModel) error {
	query := `INSERT INTO sys_management_role_elevations (created_at, updated_at, user_id, role_id, status, reason, duration_minutes, requested_by)
	           VALUES (NOW(), NOW(), :user_id, :role_id, :status, :reason, :duration_minutes, :requested_by)
	           RETURNING id`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.GetContext(ctx, &elevation.ID, elevation)
}

func (r *roleElevationRepository) Update(ctx context.Context, elevation *sysManagement.RoleElevationModel) error {
	query := `UPDATE sys_management_role_elevations SET
	           updated_at = NOW(),
	           status = :status,
	           expires_at = :expires_at,
	           decided_by = :decided_by
	         WHERE id = :id AND deleted_at IS NULL`
	_, err := r.db.NamedExecContext(ctx, query, elevation)
	return err
}

func (r *roleElevationRepository) List(ctx context.Context) ([]*sysManagement.RoleElevationModel, error) {
	var elevations []*sysManagement.RoleElevationModel
	err := r.db.SelectContext(ctx, &elevations,
		"SELECT "+elevationColumns+" FROM sys_management_role_elevations WHERE deleted_at IS NULL ORDER BY id DESC")
	return elevations, err
}

func (r *roleElevationRepository) ListByUser(ctx context.Context, userID uint) ([]*sysManagement.RoleElevationModel, error) {
	var elevations []*sysManagement.RoleElevationModel
	err := r.db.SelectContext(ctx, &elevations,
		"SELECT "+elevationColumns+" FROM sys_management_role_elevations WHERE user_id=$1 AND deleted_at IS NULL ORDER BY id DESC", userID)
	return elevations, err
}

// ActiveRoleIDs returns the distinct role ids currently granted through
// approved, unexpired elevations. This is the query live role resolution uses.
func (r *roleElevationRepository) ActiveRoleIDs(ctx context.Context, userID uint) ([]uint, error) {
	var roleIDs []uint
	err := r.db.SelectContext(ctx, &roleIDs, `
		SELECT DISTINCT role_id FROM sys_management_role_elevations
		WHERE user_id=$1 AND status=$2 AND expires_at > NOW() AND deleted_at IS NULL`,
		userID, sysManagement.ElevationApproved)
	return roleIDs, err
}
