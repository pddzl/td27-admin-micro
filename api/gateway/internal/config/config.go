package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	Auth struct {
		AccessSecret string
		AccessExpire int64
		// EnforceRbac enables per-request authorization against the rpc Casbin
		// enforcer. Requires api-domain permissions to be populated in
		// sys_management_permission; disabled by default.
		EnforceRbac bool
	}

	BasisRpc zrpc.RpcClientConf

	Cors struct {
		AllowedOrigins []string
		AllowedMethods []string
		AllowedHeaders []string
	}

	Captcha struct {
		KeyLong   int `json:"key-long"`
		ImgWidth  int `json:"img-width"`
		ImgHeight int `json:"img-height"`
	}
}
