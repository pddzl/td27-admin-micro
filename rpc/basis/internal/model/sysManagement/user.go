package sysManagement

import (
	"td27/rpc/basis/internal/model/common"
)

// UserModel User entity with multi-role support
type UserModel struct {
	common.Td27Model
	Username string       `json:"username" db:"username"`
	Password string       `json:"-" db:"password"`
	Phone    string       `json:"phone" db:"phone"`
	Email    string       `json:"email" db:"email"`
	Active   bool         `json:"active" db:"active"`
	DeptID   uint         `json:"deptId" db:"dept_id"`
	Roles    []*RoleModel `json:"roles" gorm:"many2many:sys_management_user_roles;joinForeignKey:user_id;joinReferences:role_id" db:"-"`
}

func (UserModel) TableName() string {
	return "sys_management_user"
}

// GetPrimaryRoleID returns the first role ID for backward compatibility
func (u *UserModel) GetPrimaryRoleID() uint {
	if len(u.Roles) > 0 {
		return u.Roles[0].ID
	}
	return 0
}

// GetAllRoleIDs returns all role IDs associated with the user
func (u *UserModel) GetAllRoleIDs() []uint {
	ids := make([]uint, 0, len(u.Roles))
	for _, role := range u.Roles {
		ids = append(ids, role.ID)
	}
	return ids
}

// HasRole checks if user has the specified role
func (u *UserModel) HasRole(roleID uint) bool {
	for _, role := range u.Roles {
		if role.ID == roleID {
			return true
		}
	}
	return false
}
