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

	// RateLimit controls per-IP request rate limiting and global concurrency
	// throttling. Defaults apply when Enabled is true but values are zero:
	// Rate 50 req/s per IP (burst 100), SensitiveRate 0.5 req/s per IP
	// (burst 10) for auth endpoints, MaxConcurrent 200 in-flight requests.
	RateLimit struct {
		Enabled        bool    `json:",optional"`
		Rate           float64 `json:",optional"`
		Burst          int     `json:",optional"`
		SensitiveRate  float64 `json:",optional"`
		SensitiveBurst int     `json:",optional"`
		MaxConcurrent  int     `json:",optional"`
	}
}
