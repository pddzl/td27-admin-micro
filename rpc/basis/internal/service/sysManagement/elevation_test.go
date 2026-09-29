package sysManagement

import "testing"

func TestClampElevationMinutes(t *testing.T) {
	cases := []struct {
		in   int
		want int
	}{
		{0, minElevationMinutes},
		{-10, minElevationMinutes},
		{4, minElevationMinutes},
		{5, 5},
		{60, 60},
		{maxElevationMinutes, maxElevationMinutes},
		{maxElevationMinutes + 1, maxElevationMinutes},
	}

	for _, tc := range cases {
		if got := clampElevationMinutes(tc.in); got != tc.want {
			t.Errorf("clampElevationMinutes(%d) = %d, want %d", tc.in, got, tc.want)
		}
	}
}
