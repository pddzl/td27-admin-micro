package authority

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
	Name     string           `json:"name" gorm:"size:100;unique;not null;comment:权限名称"`
	Domain   PermissionDomain `gorm:"type:varchar(20);not null;check:domain IN ('menu','api','button','data')"`
	Resource string           `json:"resource" gorm:"size:200;unique;not null;comment:资源标识，如:/api/user"`
	Action   Action           `json:"action" gorm:"size:20;default:'all';comment:操作:all|view|create|update|delete"`
	DomainID uint             `json:"domainId" gorm:"index;comment:关联领域表ID(menu/api/button)"`
}

func (PermissionModel) TableName() string {
	return "sys_management_permission"
}

// RolePermissionModel Role-Permission join table
type RolePermissionModel struct {
	RoleID       uint `gorm:"column:role_id;primaryKey"`
	PermissionID uint `gorm:"column:permission_id;primaryKey"`
}

func (RolePermissionModel) TableName() string {
	return "sys_management_role_permissions"
}

// PermissionIdentity Permission identity for Casbin
type PermissionIdentity struct {
	Type     string `json:"type"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}
