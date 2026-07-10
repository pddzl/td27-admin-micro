package sysManagement

import (
	"td27/rpc/basis/internal/model/common"
)

// RoleModel Role entity with hierarchical inheritance support
type RoleModel struct {
	common.Td27Model
	RoleName       string       `json:"roleName" db:"role_name"`
	ParentID       *uint        `json:"parentId" db:"parent_id"`
	Parent         *RoleModel   `json:"parent,omitempty" gorm:"foreignKey:ParentID" db:"-"`
	Children       []*RoleModel `json:"children,omitempty" gorm:"foreignKey:ParentID" db:"-"`
	PermissionHash string       `json:"-" db:"permission_hash"`
}

func (RoleModel) TableName() string {
	return "sys_management_role"
}

// UserRole User-Role join table for multi-role support
type UserRole struct {
	UserID uint `gorm:"column:user_id;primaryKey" db:"user_id"`
	RoleID uint `gorm:"column:role_id;primaryKey" db:"role_id"`
}

func (UserRole) TableName() string {
	return "sys_management_user_roles"
}
