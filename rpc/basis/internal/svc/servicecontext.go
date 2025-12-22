package svc

import (
	"basis/internal/model/authority"
	"log"

	"basis/internal/config"
	"basis/internal/initialization"
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
		log.Panicf("init rawDB err: %v", err)
		return nil
	} else {
		log.Println("init rawDB success")
	}

	return &ServiceContext{
		Config:            c,
		AuthorityUserRepo: authority.NewAuthorityUserModel(conn),
	}
}
