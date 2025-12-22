package authority

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AuthorityRoleModel = (*customAuthorityRoleModel)(nil)

type (
	// AuthorityRoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAuthorityRoleModel.
	AuthorityRoleModel interface {
		authorityRoleModel
		withSession(session sqlx.Session) AuthorityRoleModel
	}

	customAuthorityRoleModel struct {
		*defaultAuthorityRoleModel
	}
)

// NewAuthorityRoleModel returns a model for the database table.
func NewAuthorityRoleModel(conn sqlx.SqlConn) AuthorityRoleModel {
	return &customAuthorityRoleModel{
		defaultAuthorityRoleModel: newAuthorityRoleModel(conn),
	}
}

func (m *customAuthorityRoleModel) withSession(session sqlx.Session) AuthorityRoleModel {
	return NewAuthorityRoleModel(sqlx.NewSqlConnFromSession(session))
}
