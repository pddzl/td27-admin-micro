package sysManagement

import (
	"time"

	"td27/rpc/basis/internal/model/common"
)

// ElevationStatus lifecycle of a JIT role elevation request.
type ElevationStatus string

const (
	ElevationPending  ElevationStatus = "pending"
	ElevationApproved ElevationStatus = "approved"
	ElevationRejected ElevationStatus = "rejected"
	ElevationRevoked  ElevationStatus = "revoked"
)

// RoleElevationModel Temporary role grant for just-in-time access.
type RoleElevationModel struct {
	common.Td27Model
	UserID          uint            `json:"userId" db:"user_id"`
	RoleID          uint            `json:"roleId" db:"role_id"`
	Status          ElevationStatus `json:"status" db:"status"`
	Reason          string          `json:"reason" db:"reason"`
	DurationMinutes int             `json:"durationMinutes" db:"duration_minutes"`
	ExpiresAt       *time.Time      `json:"expiresAt" db:"expires_at"`
	RequestedBy     uint            `json:"requestedBy" db:"requested_by"`
	DecidedBy       *uint           `json:"decidedBy" db:"decided_by"`
}

func (RoleElevationModel) TableName() string {
	return "sys_management_role_elevations"
}

// Active reports whether the elevation currently grants its role. Expiry is
// evaluated at check time so authorization stops the moment the window ends.
func (e *RoleElevationModel) Active(now time.Time) bool {
	return e.Status == ElevationApproved && e.ExpiresAt != nil && e.ExpiresAt.After(now)
}
