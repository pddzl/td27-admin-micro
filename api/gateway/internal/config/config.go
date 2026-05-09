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
