package middleware

import (
	"testing"

	"td27/rpc/basis/types/sysManagement/permission_pb"
)

func TestHTTPMethodToAction(t *testing.T) {
	cases := []struct {
		name   string
		method string
		path   string
		want   permission_pb.Action
	}{
		// GET requests are always read
		{"get maps to read", "GET", "/user/get-user-info", permission_pb.Action_ACTION_READ},
		{"get path param maps to read", "GET", "/api/dept/123", permission_pb.Action_ACTION_READ},

		// POST query endpoints map to read
		{"post list maps to read", "POST", "/user/list", permission_pb.Action_ACTION_READ},
		{"post flat maps to read", "POST", "/dict-detail/flat", permission_pb.Action_ACTION_READ},
		{"post el-tree maps to read", "POST", "/menu/el-tree", permission_pb.Action_ACTION_READ},
		{"post by-page maps to read", "POST", "/button/by-page", permission_pb.Action_ACTION_READ},
		{"post cache get maps to read", "POST", "/cache/get", permission_pb.Action_ACTION_READ},
		{"post permission check maps to read", "POST", "/permission/check", permission_pb.Action_ACTION_READ},
		{"post token validate maps to read", "POST", "/service-token/validate", permission_pb.Action_ACTION_READ},
		{"post batch-check maps to read", "POST", "/button/batch-check", permission_pb.Action_ACTION_READ},
		{"post user-menus maps to read", "POST", "/menu/user-menus", permission_pb.Action_ACTION_READ},
		{"post user-buttons maps to read", "POST", "/button/user-buttons", permission_pb.Action_ACTION_READ},
		{"post get-el-tree-depts maps to read", "POST", "/dept/get-el-tree-depts", permission_pb.Action_ACTION_READ},

		// POST mutations map by verb suffix
		{"post create maps to create", "POST", "/user/create", permission_pb.Action_ACTION_CREATE},
		{"post delete maps to delete", "POST", "/role/delete", permission_pb.Action_ACTION_DELETE},
		{"post delete-by-ids maps to delete", "POST", "/api/apis/delete-by-ids", permission_pb.Action_ACTION_DELETE},
		{"post execute maps to execute", "POST", "/cron/execute", permission_pb.Action_ACTION_EXECUTE},
		{"post cleanup maps to execute", "POST", "/operation-log/cleanup", permission_pb.Action_ACTION_EXECUTE},
		{"post reload-policy maps to execute", "POST", "/permission/reload-policy", permission_pb.Action_ACTION_EXECUTE},
		{"post update maps to update", "POST", "/dict/update", permission_pb.Action_ACTION_UPDATE},
		{"post toggle-status maps to update", "POST", "/cron/toggle-status", permission_pb.Action_ACTION_UPDATE},
		{"post assign-permissions maps to update", "POST", "/role/assign-permissions", permission_pb.Action_ACTION_UPDATE},
		{"post modify-password maps to update", "POST", "/user/modify-password", permission_pb.Action_ACTION_UPDATE},
		{"post switch-active maps to update", "POST", "/user/switch-active", permission_pb.Action_ACTION_UPDATE},
		{"post cache set maps to update", "POST", "/cache/set", permission_pb.Action_ACTION_UPDATE},
		{"post upload maps to update", "POST", "/file/upload", permission_pb.Action_ACTION_UPDATE},

		// suffix precedence: delete wins over substring collisions
		{"delete suffix inside detail path", "POST", "/dict-detail/delete", permission_pb.Action_ACTION_DELETE},
		{"cleanup is not confused with list", "POST", "/cache/cleanup", permission_pb.Action_ACTION_EXECUTE},
		{"nested module paths are matched on last segment", "POST", "/service-token/toggle-status", permission_pb.Action_ACTION_UPDATE},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := httpMethodToAction(tc.method, tc.path)
			if got != tc.want {
				t.Fatalf("httpMethodToAction(%q, %q) = %v, want %v", tc.method, tc.path, got, tc.want)
			}
		})
	}
}

func TestRoleIdsFromClaims(t *testing.T) {
	cases := []struct {
		name  string
		claim interface{}
		want  []int64
	}{
		{
			name:  "json decoded float64 slice",
			claim: []interface{}{float64(1), float64(2), float64(3)},
			want:  []int64{1, 2, 3},
		},
		{
			name:  "empty slice",
			claim: []interface{}{},
			want:  []int64{},
		},
		{
			name:  "nil claim",
			claim: nil,
			want:  nil,
		},
		{
			name:  "non-slice claim",
			claim: "1,2,3",
			want:  nil,
		},
		{
			name:  "slice with non-numeric entries skips them",
			claim: []interface{}{float64(7), "admin", nil},
			want:  []int64{7},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := roleIdsFromClaims(tc.claim)
			if len(got) != len(tc.want) {
				t.Fatalf("roleIdsFromClaims(%v) = %v, want %v", tc.claim, got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("roleIdsFromClaims(%v) = %v, want %v", tc.claim, got, tc.want)
				}
			}
		})
	}
}

func TestIsSelfService(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"/user/modify-password", true},
		{"/role/elevation/request", true},
		{"/role/elevation/my", true},
		{"/user/list", false},
		{"/role/list", false},
		// elevation decisions are privileged, never self-service
		{"/role/elevation/decide", false},
		{"/role/elevation/revoke", false},
		{"/role/elevation/list", false},
		{"", false},
	}
	for _, tc := range cases {
		if got := isSelfService(tc.path); got != tc.want {
			t.Errorf("isSelfService(%q) = %v, want %v", tc.path, got, tc.want)
		}
	}
}
