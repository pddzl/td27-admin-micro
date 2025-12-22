package authority

import (
	"gorm.io/gorm"
)

var _ AuthorityMenuModel = (*customAuthorityMenuModel)(nil)

type (
	// AuthorityMenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAuthorityMenuModel.
	AuthorityMenuModel interface {
		authorityMenuModel
	}

	customAuthorityMenuModel struct {
		*defaultAuthorityMenuModel
	}
)

// NewAuthorityMenuModel returns a model for the database table.
func NewAuthorityMenuModel(conn *gorm.DB) AuthorityMenuModel {
	return &customAuthorityMenuModel{
		defaultAuthorityMenuModel: newAuthorityMenuModel(conn),
	}
}
