package main

import (
	"flag"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"td27/rpc/basis/internal/config"
	"td27/rpc/basis/internal/server/authority"
	"td27/rpc/basis/internal/server/tool"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/authority/user_pb"
	"td27/rpc/basis/types/authority/role_pb"
	"td27/rpc/basis/types/authority/permission_pb"
	"td27/rpc/basis/types/authority/menu_pb"
	"td27/rpc/basis/types/authority/dept_pb"
	"td27/rpc/basis/types/authority/dict_pb"
	"td27/rpc/basis/types/authority/api_pb"
	"td27/rpc/basis/types/authority/button_pb"
	"td27/rpc/basis/types/tool/file_pb"
	"td27/rpc/basis/types/tool/cron_pb"
	"td27/rpc/basis/types/tool/cache_pb"
	"td27/rpc/basis/types/tool/service_token_pb"
	"td27/rpc/basis/types/monitor/operation_log_pb"
	"td27/rpc/basis/internal/server/monitor"
)

var configFile = flag.String("f", "etc/basis.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user_pb.RegisterUserServer(grpcServer, authority.NewUserServer(ctx))
		role_pb.RegisterRoleServer(grpcServer, authority.NewRoleServer(ctx))
		permission_pb.RegisterPermissionServer(grpcServer, authority.NewPermissionServer(ctx))
		menu_pb.RegisterMenuServer(grpcServer, authority.NewMenuServer(ctx))
		dept_pb.RegisterDeptServer(grpcServer, authority.NewDeptServer(ctx))
		dict_pb.RegisterDictServer(grpcServer, authority.NewDictServer(ctx))
		api_pb.RegisterAPIServer(grpcServer, authority.NewAPIServer(ctx))
		button_pb.RegisterButtonServer(grpcServer, authority.NewButtonServer(ctx))
		file_pb.RegisterFileServer(grpcServer, tool.NewFileServer(ctx))
		cron_pb.RegisterCronServer(grpcServer, tool.NewCronServer(ctx))
		cache_pb.RegisterCacheServer(grpcServer, tool.NewCacheServer(ctx))
		service_token_pb.RegisterServiceTokenServer(grpcServer, tool.NewServiceTokenServer(ctx))
		operation_log_pb.RegisterOperationLogServer(grpcServer, monitor.NewOperationLogServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	logx.Infof("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
