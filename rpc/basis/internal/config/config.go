package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	//Log      logx.LogConf
	DataBase DataBase
}
