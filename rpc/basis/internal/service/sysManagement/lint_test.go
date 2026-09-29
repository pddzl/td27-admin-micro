package sysManagement

import (
	"strings"
	"testing"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
)

func perm(id uint, domain sysManagement.PermissionDomain, resource string, action sysManagement.Action, effect sysManagement.Effect) *sysManagement.PermissionModel {
	if effect == "" {
		effect = sysManagement.EffectAllow
	}
	return &sysManagement.PermissionModel{
		Td27Model: common.Td27Model{ID: id},
		Domain:    domain,
		Resource:  resource,
		Action:    action,
		Effect:    effect,
	}
}

func findIssue(t *testing.T, issues []PolicyIssue, severity string, permissionID uint) PolicyIssue {
	t.Helper()
	for _, issue := range issues {
		if issue.Severity == severity && issue.PermissionID == permissionID {
			return issue
		}
	}
	t.Fatalf("expected %s issue for permission %d, got %+v", severity, permissionID, issues)
	return PolicyIssue{}
}

func TestLintConflictDenyKillsAllow(t *testing.T) {
	assignments := map[uint][]*sysManagement.PermissionModel{
		1: {
			perm(10, sysManagement.PermissionDomainAPI, "/user/*", sysManagement.ActionAll, sysManagement.EffectAllow),
			perm(11, sysManagement.PermissionDomainAPI, "/*", sysManagement.ActionAll, sysManagement.EffectDeny),
		},
	}

	issues := LintAssignments(assignments, nil)
	issue := findIssue(t, issues, SeverityConflict, 10)
	if issue.RoleID != 1 {
		t.Errorf("expected role 1, got %d", issue.RoleID)
	}
	if !strings.Contains(issue.Message, "deny rule 11") {
		t.Errorf("message should reference the denying rule, got: %s", issue.Message)
	}
}

func TestLintNoConflictWhenDenyIsNarrower(t *testing.T) {
	// a narrower deny (/user/delete) only overrides the matching requests of
	// the broader allow — the allow still grants everything else, so it is
	// not a (full) conflict.
	assignments := map[uint][]*sysManagement.PermissionModel{
		1: {
			perm(10, sysManagement.PermissionDomainAPI, "/*", sysManagement.ActionAll, sysManagement.EffectAllow),
			perm(11, sysManagement.PermissionDomainAPI, "/user/delete", sysManagement.ActionDelete, sysManagement.EffectDeny),
		},
	}

	for _, issue := range LintAssignments(assignments, nil) {
		if issue.Severity == SeverityConflict {
			t.Fatalf("unexpected conflict: %+v", issue)
		}
	}
}

func TestLintNoConflictWhenDenyIsNarrowerButNotCovering(t *testing.T) {
	// deny targets a different action — the allow on /user/create stays valid.
	assignments := map[uint][]*sysManagement.PermissionModel{
		1: {
			perm(10, sysManagement.PermissionDomainAPI, "/user/*", sysManagement.ActionCreate, sysManagement.EffectAllow),
			perm(11, sysManagement.PermissionDomainAPI, "/user/*", sysManagement.ActionDelete, sysManagement.EffectDeny),
		},
	}

	for _, issue := range LintAssignments(assignments, nil) {
		if issue.Severity == SeverityConflict {
			t.Fatalf("unexpected conflict: %+v", issue)
		}
	}
}

func TestLintShadowedRule(t *testing.T) {
	assignments := map[uint][]*sysManagement.PermissionModel{
		2: {
			perm(20, sysManagement.PermissionDomainAPI, "/*", sysManagement.ActionAll, sysManagement.EffectAllow),
			perm(21, sysManagement.PermissionDomainAPI, "/user/*", sysManagement.ActionAll, sysManagement.EffectAllow),
		},
	}

	issues := LintAssignments(assignments, nil)
	issue := findIssue(t, issues, SeverityShadowed, 21)
	if issue.RoleID != 2 {
		t.Errorf("expected role 2, got %d", issue.RoleID)
	}
}

func TestLintShadowedBySameResourceBroaderAction(t *testing.T) {
	assignments := map[uint][]*sysManagement.PermissionModel{
		3: {
			perm(30, sysManagement.PermissionDomainAPI, "/role/*", sysManagement.ActionAll, sysManagement.EffectAllow),
			perm(31, sysManagement.PermissionDomainAPI, "/role/*", sysManagement.ActionDelete, sysManagement.EffectAllow),
		},
	}

	findIssue(t, LintAssignments(assignments, nil), SeverityShadowed, 31)
}

func TestLintNoShadowAcrossDifferentEffects(t *testing.T) {
	// an allow and a deny with identical shape are a conflict for the allow,
	// not a shadowing pair.
	assignments := map[uint][]*sysManagement.PermissionModel{
		4: {
			perm(40, sysManagement.PermissionDomainAPI, "/dept/*", sysManagement.ActionAll, sysManagement.EffectAllow),
			perm(41, sysManagement.PermissionDomainAPI, "/dept/*", sysManagement.ActionAll, sysManagement.EffectDeny),
		},
	}

	findIssue(t, LintAssignments(assignments, nil), SeverityConflict, 40)
	for _, issue := range LintAssignments(assignments, nil) {
		if issue.Severity == SeverityShadowed {
			t.Fatalf("unexpected shadowed issue: %+v", issue)
		}
	}
}

func TestLintUnassignedAPIPermission(t *testing.T) {
	all := []*sysManagement.PermissionModel{
		perm(50, sysManagement.PermissionDomainAPI, "/user/list", sysManagement.ActionRead, ""),
		perm(51, sysManagement.PermissionDomainAPI, "/user/create", sysManagement.ActionCreate, ""),
		perm(52, sysManagement.PermissionDomainMenu, "sys.user", sysManagement.ActionView, ""),
	}
	assignments := map[uint][]*sysManagement.PermissionModel{
		5: {perm(50, sysManagement.PermissionDomainAPI, "/user/list", sysManagement.ActionRead, "")},
	}

	issues := LintAssignments(assignments, all)
	issue := findIssue(t, issues, SeverityUnassigned, 51)
	if issue.RoleID != 0 {
		t.Errorf("unassigned issues have no role, got %d", issue.RoleID)
	}
	// menu-domain permission 52 is exempt from the unassigned check
	for _, i := range issues {
		if i.PermissionID == 52 {
			t.Fatalf("menu permission should not be flagged unassigned: %+v", i)
		}
	}
}

func TestLintCleanPolicies(t *testing.T) {
	assignments := map[uint][]*sysManagement.PermissionModel{
		6: {
			perm(60, sysManagement.PermissionDomainAPI, "/user/*", sysManagement.ActionAll, sysManagement.EffectAllow),
			perm(61, sysManagement.PermissionDomainAPI, "/dept/*", sysManagement.ActionAll, sysManagement.EffectAllow),
		},
	}
	all := []*sysManagement.PermissionModel{
		perm(60, sysManagement.PermissionDomainAPI, "/user/*", sysManagement.ActionAll, ""),
		perm(61, sysManagement.PermissionDomainAPI, "/dept/*", sysManagement.ActionAll, ""),
	}

	if issues := LintAssignments(assignments, all); len(issues) != 0 {
		t.Fatalf("expected no issues, got %+v", issues)
	}
}
