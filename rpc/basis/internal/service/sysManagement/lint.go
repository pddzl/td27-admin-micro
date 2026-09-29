package sysManagement

import (
	"fmt"

	"github.com/casbin/casbin/v2/util"

	"td27/rpc/basis/internal/model/sysManagement"
)

// PolicyIssueSeverity levels reported by LintAssignments.
const (
	SeverityConflict  = "conflict"  // allow rule is dead: a deny rule covers it
	SeverityShadowed  = "shadowed"  // rule made redundant by a broader rule of the same role
	SeverityUnassigned = "unassigned" // api permission not attached to any role
)

// PolicyIssue is a lint finding for a role-permission assignment.
type PolicyIssue struct {
	Severity     string
	Message      string
	RoleID       uint
	PermissionID uint
}

// LintAssignments analyzes role→permission assignments and reports policy
// quality problems. assignments maps role id to its granted permissions;
// allPerms is the full permission catalogue (api-domain entries not present in
// any role's assignment are reported as unassigned).
func LintAssignments(assignments map[uint][]*sysManagement.PermissionModel, allPerms []*sysManagement.PermissionModel) []PolicyIssue {
	var issues []PolicyIssue

	assigned := make(map[uint]bool)

	for roleID, perms := range assignments {
		for i, perm := range perms {
			assigned[perm.ID] = true

			for j, other := range perms {
				if i == j {
					continue
				}

				// A deny rule covering an allow rule makes the allow dead.
				if perm.Effect != sysManagement.EffectDeny &&
					other.Effect == sysManagement.EffectDeny &&
					ruleCovers(other, perm) {
					issues = append(issues, PolicyIssue{
						Severity:     SeverityConflict,
						Message:      fmt.Sprintf("allow rule is dead: deny rule %d (id=%d) covers it", other.ID, other.ID),
						RoleID:       roleID,
						PermissionID: perm.ID,
					})
					continue
				}

				// A broader rule with the same effect makes the narrower one redundant.
				if perm.Effect == other.Effect && ruleCovers(other, perm) {
					issues = append(issues, PolicyIssue{
						Severity:     SeverityShadowed,
						Message:      fmt.Sprintf("rule is redundant: rule %d (id=%d) already covers it", other.ID, other.ID),
						RoleID:       roleID,
						PermissionID: perm.ID,
					})
				}
			}
		}
	}

	for _, perm := range allPerms {
		if perm.Domain == sysManagement.PermissionDomainAPI && !assigned[perm.ID] {
			issues = append(issues, PolicyIssue{
				Severity:     SeverityUnassigned,
				Message:      "api permission is not attached to any role",
				PermissionID: perm.ID,
			})
		}
	}

	return issues
}

// ruleCovers reports whether rule b's resource pattern and action subsume
// rule a's. Resource patterns use the same keyMatch2 semantics as the enforcer.
func ruleCovers(b, a *sysManagement.PermissionModel) bool {
	if b.Resource != a.Resource && !util.KeyMatch2(a.Resource, b.Resource) {
		return false
	}
	return b.Action == sysManagement.ActionAll || b.Action == a.Action
}
