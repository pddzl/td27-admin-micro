package config

type DataBase struct {
	Host            string
	Port            uint
	Config          string
	DBName          string
	Username        string
	Password        string
	MaxOpenConn     int
	MaxIdleConn     int
	ConnMaxLifetime int64
}
