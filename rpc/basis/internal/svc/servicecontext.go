package svc

import (
	"log"

	"td27/rpc/basis/internal/config"
	"td27/rpc/basis/internal/initialization"
	"td27/rpc/basis/internal/model/authority"
)

type ServiceContext struct {
	Config            config.Config
	AuthorityUserRepo authority.AuthorityUserModel
	AuthorityRoleRepo authority.AuthorityRoleModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn, err := initialization.Gorm(
		c.DataBase.Username,
		c.DataBase.Password,
		c.DataBase.Host,
		c.DataBase.Port,
		c.DataBase.DBName,
		c.DataBase.Config,
		c.DataBase.MaxOpenConn,
		c.DataBase.MaxIdleConn,
		c.DataBase.ConnMaxLifetime)
	if err != nil {
		log.Panicf("init mysql err: %v", err)
		return nil
	} else {
		log.Println("init mysql success")
	}

	return &ServiceContext{
		Config:            c,
		AuthorityUserRepo: authority.NewAuthorityUserModel(conn),
		AuthorityRoleRepo: authority.NewAuthorityRoleModel(conn),
	}
}
