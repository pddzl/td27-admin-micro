package authority

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AuthorityMenuModel = (*customAuthorityMenuModel)(nil)

type (
	// AuthorityMenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAuthorityMenuModel.
	AuthorityMenuModel interface {
		authorityMenuModel
		withSession(session sqlx.Session) AuthorityMenuModel
	}

	customAuthorityMenuModel struct {
		*defaultAuthorityMenuModel
	}
)

// NewAuthorityMenuModel returns a model for the database table.
func NewAuthorityMenuModel(conn sqlx.SqlConn) AuthorityMenuModel {
	return &customAuthorityMenuModel{
		defaultAuthorityMenuModel: newAuthorityMenuModel(conn),
	}
}

func (m *customAuthorityMenuModel) withSession(session sqlx.Session) AuthorityMenuModel {
	return NewAuthorityMenuModel(sqlx.NewSqlConnFromSession(session))
}
