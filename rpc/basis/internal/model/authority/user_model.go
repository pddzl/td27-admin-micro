package authority

import (
	"gorm.io/gorm"
)

var _ AuthorityUserModel = (*customAuthorityUserModel)(nil)

type (
	// AuthorityUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAuthorityUserModel.
	AuthorityUserModel interface {
		authorityUserModel
	}

	customAuthorityUserModel struct {
		*defaultAuthorityUserModel
	}
)

// NewAuthorityUserModel returns a model for the database table.
func NewAuthorityUserModel(conn *gorm.DB) AuthorityUserModel {
	return &customAuthorityUserModel{
		defaultAuthorityUserModel: newAuthorityUserModel(conn),
	}
}
