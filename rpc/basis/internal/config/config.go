package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	
	// Database configuration
	Pgsql Pgsql `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`
	
	// JWT authentication configuration
	JWT JWT `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	
	// Casbin RBAC configuration
	Casbin Casbin `mapstructure:"casbin" json:"casbin" yaml:"casbin"`
	
	// File upload configuration
	File File `mapstructure:"file" json:"file" yaml:"file"`
	
	// System configuration
	System System `mapstructure:"system" json:"system" yaml:"system"`
}

// Pgsql PostgreSQL database configuration
type Pgsql struct {
	Host         string `mapstructure:"host" json:"host" yaml:"host"`
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`
	Username     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	Dbname       string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"`
	LogLevel     string `mapstructure:"log-level" json:"log-level" yaml:"log-level"`
	LogMode      bool   `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`
}

// JWT authentication configuration
type JWT struct {
	SigningKey      string `mapstructure:"signing-key" json:"signing-key" yaml:"signing-key"`
	ExpiresTime     int64  `mapstructure:"expires-time" json:"expires-time" yaml:"expires-time"`
	BufferTime      int64  `mapstructure:"buffer-time" json:"buffer-time" yaml:"buffer-time"`
	Issuer          string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`
	MultiLogin      bool   `mapstructure:"multi-login" json:"multi-login" yaml:"multi-login"`
	MultiLoginLimit int    `mapstructure:"multi-login-limit" json:"multi-login-limit" yaml:"multi-login-limit"`
}

// Casbin RBAC configuration
type Casbin struct {
	CacheTTL              int  `mapstructure:"cache-ttl" json:"cache-ttl" yaml:"cache-ttl"`
	AutoLoadInterval      int  `mapstructure:"auto-load-interval" json:"auto-load-interval" yaml:"auto-load-interval"`
	EnableRoleHierarchy   bool `mapstructure:"enable-role-hierarchy" json:"enable-role-hierarchy" yaml:"enable-role-hierarchy"`
	EnableDataPermission  bool `mapstructure:"enable-data-permission" json:"enable-data-permission" yaml:"enable-data-permission"`
}

// File upload configuration
type File struct {
	UploadPath  string `mapstructure:"upload-path" json:"upload-path" yaml:"upload-path"`
	MaxSize     int    `mapstructure:"max-size" json:"max-size" yaml:"max-size"`
	AllowedExt  []string `mapstructure:"allowed-ext" json:"allowed-ext" yaml:"allowed-ext"`
}

// System configuration
type System struct {
	Env           string `mapstructure:"env" json:"env" yaml:"env"`
	Addr          int    `mapstructure:"addr" json:"addr" yaml:"addr"`
	AppName       string `mapstructure:"app-name" json:"app-name" yaml:"app-name"`
	OperationLog  bool   `mapstructure:"operation-log" json:"operation-log" yaml:"operation-log"`
	MaxUploadSize int    `mapstructure:"max-upload-size" json:"max-upload-size" yaml:"max-upload-size"`
}
