package authority

import (
	"gorm.io/gorm"
)

var _ AuthorityRoleModel = (*customAuthorityRoleModel)(nil)

type (
	// AuthorityRoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAuthorityRoleModel.
	AuthorityRoleModel interface {
		authorityRoleModel
	}

	customAuthorityRoleModel struct {
		*defaultAuthorityRoleModel
	}
)

// NewAuthorityRoleModel returns a model for the database table.
func NewAuthorityRoleModel(conn *gorm.DB) AuthorityRoleModel {
	return &customAuthorityRoleModel{
		defaultAuthorityRoleModel: newAuthorityRoleModel(conn),
	}
}
