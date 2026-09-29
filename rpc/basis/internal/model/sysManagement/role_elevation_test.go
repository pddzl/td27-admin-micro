package sysManagement

import (
	"testing"
	"time"
)

func TestRoleElevationActive(t *testing.T) {
	now := time.Now()
	future := now.Add(time.Hour)
	past := now.Add(-time.Hour)

	cases := []struct {
		name string
		e    RoleElevationModel
		want bool
	}{
		{"approved and unexpired", RoleElevationModel{Status: ElevationApproved, ExpiresAt: &future}, true},
		{"approved but expired", RoleElevationModel{Status: ElevationApproved, ExpiresAt: &past}, false},
		{"approved with no expiry", RoleElevationModel{Status: ElevationApproved}, false},
		{"pending never active", RoleElevationModel{Status: ElevationPending, ExpiresAt: &future}, false},
		{"rejected never active", RoleElevationModel{Status: ElevationRejected, ExpiresAt: &future}, false},
		{"revoked never active", RoleElevationModel{Status: ElevationRevoked, ExpiresAt: &future}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.e.Active(now); got != tc.want {
				t.Errorf("Active() = %v, want %v", got, tc.want)
			}
		})
	}
}
