package svc

import (
	"testing"

	"github.com/casbin/casbin/v2"
)

// newTestEnforcer builds the production Casbin model (getCasbinModel) with
// in-memory policies — no database involved. Policies mirror the shape loaded
// by the PostgresAdapter: sub = role id as string, obj = resource pattern,
// act = action string, eft = "allow" | "deny".
func newTestEnforcer(t *testing.T, enableRoleHierarchy bool, policies [][]string, grouping [][]string) *casbin.SyncedCachedEnforcer {
	t.Helper()

	m, err := getCasbinModel(enableRoleHierarchy)
	if err != nil {
		t.Fatalf("getCasbinModel(%v) failed: %v", enableRoleHierarchy, err)
	}

	e, err := casbin.NewSyncedCachedEnforcer(m)
	if err != nil {
		t.Fatalf("NewSyncedCachedEnforcer failed: %v", err)
	}

	if len(policies) > 0 {
		if _, err := e.AddPolicies(policies); err != nil {
			t.Fatalf("AddPolicies failed: %v", err)
		}
	}
	for _, rule := range grouping {
		if _, err := e.AddGroupingPolicy(rule); err != nil {
			t.Fatalf("AddGroupingPolicy(%v) failed: %v", rule, err)
		}
	}
	return e
}

// check runs one Enforce call and reports a failure with a readable message.
func check(t *testing.T, e *casbin.SyncedCachedEnforcer, sub, obj, act string, want bool) {
	t.Helper()

	got, err := e.Enforce(sub, obj, act)
	if err != nil {
		t.Fatalf("Enforce(%s, %s, %s) failed: %v", sub, obj, act, err)
	}
	if got != want {
		t.Errorf("Enforce(%s, %s, %s) = %v, want %v", sub, obj, act, got, want)
	}
}

// TestEnforceSuperAdminWildcard mirrors seed row: role 1, resource /*, action all.
func TestEnforceSuperAdminWildcard(t *testing.T) {
	e := newTestEnforcer(t, false, [][]string{{"1", "/*", "all", "allow"}}, nil)

	check(t, e, "1", "/user/create", "create", true)
	check(t, e, "1", "/user/list", "read", true)
	check(t, e, "1", "/dept/delete", "delete", true)
	check(t, e, "1", "/cron/execute", "execute", true)
	check(t, e, "1", "/api/dept/123", "read", true)

	// other roles are not granted by role 1 policies
	check(t, e, "2", "/user/list", "read", false)
}

// TestEnforceAuditorReadOnly mirrors seed row: role 2, resource /*, action read.
func TestEnforceAuditorReadOnly(t *testing.T) {
	e := newTestEnforcer(t, false, [][]string{{"2", "/*", "read", "allow"}}, nil)

	check(t, e, "2", "/user/list", "read", true)
	check(t, e, "2", "/menu/el-tree", "read", true)
	check(t, e, "2", "/user/create", "create", false)
	check(t, e, "2", "/role/delete", "delete", false)
	check(t, e, "2", "/permission/reload-policy", "execute", false)
}

// TestEnforceUserManager mirrors seed rows for role 3:
// read on /* plus all on /user/*, /role/*, /dept/*.
func TestEnforceUserManager(t *testing.T) {
	e := newTestEnforcer(t, false, [][]string{
		{"3", "/*", "read", "allow"},
		{"3", "/user/*", "all", "allow"},
		{"3", "/role/*", "all", "allow"},
		{"3", "/dept/*", "all", "allow"},
	}, nil)

	check(t, e, "3", "/user/create", "create", true)
	check(t, e, "3", "/user/switch-active", "update", true)
	check(t, e, "3", "/role/delete", "delete", true)
	check(t, e, "3", "/dept/list", "read", true)

	// outside scoped modules only read is granted
	check(t, e, "3", "/file/list", "read", true)
	check(t, e, "3", "/file/delete", "delete", false)
	check(t, e, "3", "/cron/execute", "execute", false)
}

// TestEnforceKeyMatch2Params verifies keyMatch2 ':param' patterns used for
// detail routes such as GET /api/dept/:id.
func TestEnforceKeyMatch2Params(t *testing.T) {
	e := newTestEnforcer(t, false, [][]string{{"5", "/api/dept/:id", "read", "allow"}}, nil)

	check(t, e, "5", "/api/dept/123", "read", true)
	check(t, e, "5", "/api/dept/abc", "read", true)
	check(t, e, "5", "/api/dept/123/nested", "read", false)
	check(t, e, "5", "/api/dict/123", "read", false)
}

// TestEnforceAllActionIsWildcard verifies the matcher treats p.act == 'all'
// as an action wildcard, not a literal match.
func TestEnforceAllActionIsWildcard(t *testing.T) {
	e := newTestEnforcer(t, false, [][]string{{"7", "/file/*", "all", "allow"}}, nil)

	check(t, e, "7", "/file/upload", "update", true)
	check(t, e, "7", "/file/list", "read", true)
	check(t, e, "7", "/file/delete", "delete", true)
	check(t, e, "7", "/file/download/1", "read", true)
}

// TestEnforceRoleHierarchy verifies the role-hierarchy model variant
// (g(r.sub, p.sub)) grants inherited roles and the flat variant does not.
func TestEnforceRoleHierarchy(t *testing.T) {
	policies := [][]string{{"1", "/user/*", "all", "allow"}}
	grouping := [][]string{{"10", "1"}} // role 10 inherits role 1

	hierarchical := newTestEnforcer(t, true, policies, grouping)
	check(t, hierarchical, "10", "/user/create", "create", true)
	check(t, hierarchical, "10", "/role/create", "create", false)

	flat := newTestEnforcer(t, false, policies, grouping)
	check(t, flat, "10", "/user/create", "create", false)
	check(t, flat, "1", "/user/create", "create", true)
}

// TestEnforceSubjectIsolation verifies one role's policies never leak to others.
func TestEnforceSubjectIsolation(t *testing.T) {
	e := newTestEnforcer(t, false, [][]string{
		{"3", "/user/*", "all", "allow"},
		{"4", "/role/*", "all", "allow"},
	}, nil)

	check(t, e, "3", "/user/create", "create", true)
	check(t, e, "3", "/role/create", "create", false)
	check(t, e, "4", "/role/create", "create", true)
	check(t, e, "4", "/user/create", "create", false)
}

// TestDenyOverridesAllow verifies the deny-override policy effect: a matching
// deny rule wins over any allow rule, including the super-admin wildcard.
func TestDenyOverridesAllow(t *testing.T) {
	e := newTestEnforcer(t, false, [][]string{
		{"1", "/*", "all", "allow"},
		{"1", "/user/delete", "delete", "deny"},
	}, nil)

	// deny beats the /* all wildcard
	check(t, e, "1", "/user/delete", "delete", false)
	// sibling actions and resources remain granted
	check(t, e, "1", "/user/create", "create", true)
	check(t, e, "1", "/user/switch-active", "update", true)
	check(t, e, "1", "/role/delete", "delete", true)
}

// TestDenyWildcardScopedToPattern verifies a deny rule only blocks requests
// its own resource pattern and action match.
func TestDenyWildcardScopedToPattern(t *testing.T) {
	e := newTestEnforcer(t, false, [][]string{
		{"2", "/*", "all", "allow"},
		{"2", "/cron/*", "execute", "deny"},
	}, nil)

	check(t, e, "2", "/cron/execute", "execute", false)
	check(t, e, "2", "/cron/add", "update", true)     // different action, deny not matched
	check(t, e, "2", "/file/delete", "delete", true)  // different resource
}

// TestDenyOnlyRoleIsFullyDenied verifies a role with only deny rules has no
// implicit access (allow-effect still required).
func TestDenyOnlyRoleIsFullyDenied(t *testing.T) {
	e := newTestEnforcer(t, false, [][]string{{"3", "/user/*", "all", "deny"}}, nil)

	check(t, e, "3", "/user/delete", "delete", false)
	check(t, e, "3", "/user/list", "read", false)
}

// TestDenyOverridesAllowWithHierarchy verifies deny rules override inherited
// allows in the hierarchical model.
func TestDenyOverridesAllowWithHierarchy(t *testing.T) {
	policies := [][]string{
		{"1", "/*", "all", "allow"},
		{"10", "/role/delete", "delete", "deny"},
	}
	grouping := [][]string{{"10", "1"}}

	e := newTestEnforcer(t, true, policies, grouping)

	// inherited allow works...
	check(t, e, "10", "/user/create", "create", true)
	// ...but the role's own deny overrides the inherited wildcard
	check(t, e, "10", "/role/delete", "delete", false)
	check(t, e, "1", "/role/delete", "delete", true) // deny is subject-scoped
}
