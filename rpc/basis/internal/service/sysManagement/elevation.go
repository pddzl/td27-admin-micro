package sysManagement

import (
	"context"
	"errors"
	"fmt"
	"time"

	"td27/rpc/basis/internal/model/sysManagement"
	sysManagementRepo "td27/rpc/basis/internal/repository/sysManagement"
)

// Elevation bounds: a request may propose any window, but the server clamps
// it — JIT grants are meant to be short-lived.
const (
	minElevationMinutes = 5
	maxElevationMinutes = 24 * 60
)

// ElevationService manages just-in-time role elevation requests.
type ElevationService struct {
	elevationRepo sysManagementRepo.RoleElevationRepository
	roleRepo      sysManagementRepo.RoleRepository
}

func NewElevationService(
	elevationRepo sysManagementRepo.RoleElevationRepository,
	roleRepo sysManagementRepo.RoleRepository,
) *ElevationService {
	return &ElevationService{elevationRepo: elevationRepo, roleRepo: roleRepo}
}

// CreateRequest files a new elevation request for the given role. The caller
// identity (userID) comes from the authenticated request, never the body.
func (s *ElevationService) CreateRequest(ctx context.Context, userID, roleID uint, reason string, durationMinutes int) (*sysManagement.RoleElevationModel, error) {
	if userID == 0 {
		return nil, errors.New("invalid user")
	}

	role, err := s.roleRepo.FindOne(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("role %d not found", roleID)
	}

	pending, err := s.elevationRepo.FindPendingByUserAndRole(ctx, userID, roleID)
	if err != nil {
		return nil, err
	}
	if pending != nil {
		return nil, errors.New("an elevation request for this role is already pending")
	}

	elevation := &sysManagement.RoleElevationModel{
		UserID:          userID,
		RoleID:          roleID,
		Status:          sysManagement.ElevationPending,
		Reason:          reason,
		DurationMinutes: clampElevationMinutes(durationMinutes),
		RequestedBy:     userID,
	}

	if err := s.elevationRepo.Create(ctx, elevation); err != nil {
		return nil, err
	}
	return elevation, nil
}

// Decide approves or rejects a pending request. Approved requests expire at
// now + duration (unless the approver overrides expiry).
func (s *ElevationService) Decide(ctx context.Context, elevationID, approverID uint, approve bool, expiresAt *time.Time) error {
	elevation, err := s.mustFind(ctx, elevationID)
	if err != nil {
		return err
	}
	if elevation.Status != sysManagement.ElevationPending {
		return fmt.Errorf("elevation %d is not pending (status: %s)", elevationID, elevation.Status)
	}

	if approverID != 0 && elevation.UserID == approverID {
		return errors.New("cannot approve your own elevation request")
	}

	if !approve {
		elevation.Status = sysManagement.ElevationRejected
		elevation.DecidedBy = &approverID
		return s.elevationRepo.Update(ctx, elevation)
	}

	if expiresAt == nil {
		t := time.Now().Add(time.Duration(elevation.DurationMinutes) * time.Minute)
		expiresAt = &t
	}
	elevation.Status = sysManagement.ElevationApproved
	elevation.ExpiresAt = expiresAt
	elevation.DecidedBy = &approverID

	return s.elevationRepo.Update(ctx, elevation)
}

// Revoke ends an approved elevation immediately.
func (s *ElevationService) Revoke(ctx context.Context, elevationID, revokerID uint) error {
	elevation, err := s.mustFind(ctx, elevationID)
	if err != nil {
		return err
	}
	if elevation.Status != sysManagement.ElevationApproved {
		return fmt.Errorf("elevation %d is not approved (status: %s)", elevationID, elevation.Status)
	}

	elevation.Status = sysManagement.ElevationRevoked
	elevation.DecidedBy = &revokerID
	return s.elevationRepo.Update(ctx, elevation)
}

func (s *ElevationService) List(ctx context.Context) ([]*sysManagement.RoleElevationModel, error) {
	return s.elevationRepo.List(ctx)
}

func (s *ElevationService) ListByUser(ctx context.Context, userID uint) ([]*sysManagement.RoleElevationModel, error) {
	return s.elevationRepo.ListByUser(ctx, userID)
}

// ActiveRoleIDs returns the role ids currently granted via elevations.
func (s *ElevationService) ActiveRoleIDs(ctx context.Context, userID uint) ([]uint, error) {
	return s.elevationRepo.ActiveRoleIDs(ctx, userID)
}

func (s *ElevationService) mustFind(ctx context.Context, id uint) (*sysManagement.RoleElevationModel, error) {
	if id == 0 {
		return nil, errors.New("invalid elevation id")
	}
	elevation, err := s.elevationRepo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	if elevation == nil {
		return nil, fmt.Errorf("elevation %d not found", id)
	}
	return elevation, nil
}

func clampElevationMinutes(minutes int) int {
	if minutes < minElevationMinutes {
		return minElevationMinutes
	}
	if minutes > maxElevationMinutes {
		return maxElevationMinutes
	}
	return minutes
}
