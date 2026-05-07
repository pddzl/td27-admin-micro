package authority

import (
	"td27/rpc/basis/internal/model/common"
)

// RoleModel Role entity with hierarchical inheritance support
type RoleModel struct {
	common.Td27Model
	RoleName       string       `json:"roleName" gorm:"unique;size:191" binding:"required"`
	ParentID       *uint        `json:"parentId" gorm:"index;comment:父角色ID"`
	Parent         *RoleModel   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children       []*RoleModel `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	PermissionHash string       `json:"-" gorm:"comment:权限哈希，用于缓存失效判断"`
}

func (RoleModel) TableName() string {
	return "sys_management_role"
}

// UserRole User-Role join table for multi-role support
type UserRole struct {
	UserID uint `gorm:"column:user_id;primaryKey"`
	RoleID uint `gorm:"column:role_id;primaryKey"`
}

func (UserRole) TableName() string {
	return "sys_management_user_roles"
}
