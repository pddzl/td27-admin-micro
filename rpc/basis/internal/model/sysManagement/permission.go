package sysManagement

import (
	"td27/rpc/basis/internal/model/common"
)

// PermissionDomain Permission type
type PermissionDomain string

const (
	PermissionDomainMenu   PermissionDomain = "menu"
	PermissionDomainAPI    PermissionDomain = "api"
	PermissionDomainButton PermissionDomain = "button"
	PermissionDomainData   PermissionDomain = "data"
)

type Action string

func (a Action) String() string {
	return string(a)
}

const (
	ActionAll     Action = "all"
	ActionView    Action = "view"    // menu
	ActionRead    Action = "read"    // api
	ActionCreate  Action = "create"  // api
	ActionUpdate  Action = "update"  // api
	ActionDelete  Action = "delete"  // api
	ActionExecute Action = "execute" // button
)

// HTTPMethodToAction maps HTTP methods to CRUD actions
func HTTPMethodToAction(method string) Action {
	switch method {
	case "GET":
		return ActionRead
	case "POST":
		return ActionCreate
	case "PUT", "PATCH":
		return ActionUpdate
	case "DELETE":
		return ActionDelete
	default:
		return ActionAll
	}
}

// PermissionModel Unified permission table for RBAC authorization
type PermissionModel struct {
	common.Td27Model
	Name     string           `json:"name" db:"name"`
	Domain   PermissionDomain `gorm:"type:varchar(20);not null;check:domain IN ('menu','api','button','data')" db:"domain"`
	Resource string           `json:"resource" db:"resource"`
	Action   Action           `json:"action" db:"action"`
	DomainID uint             `json:"domainId" db:"domain_id"`
}

func (PermissionModel) TableName() string {
	return "sys_management_permission"
}

// RolePermissionModel Role-Permission join table
type RolePermissionModel struct {
	RoleID       uint `gorm:"column:role_id;primaryKey" db:"role_id"`
	PermissionID uint `gorm:"column:permission_id;primaryKey" db:"permission_id"`
}

func (RolePermissionModel) TableName() string {
	return "sys_management_role_permissions"
}

// PermissionIdentity Permission identity for Casbin
type PermissionIdentity struct {
	Type     string `json:"type"`
	Resource string `json:"resource" db:"resource"`
	Action   string `json:"action" db:"action"`
}
