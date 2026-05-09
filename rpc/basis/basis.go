package main

import (
	"flag"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"

	"td27/rpc/basis/internal/config"
	"td27/rpc/basis/internal/svc"

	basisServer "td27/rpc/basis/internal/server/basis"
	"td27/rpc/basis/internal/server/sysManagement"
	"td27/rpc/basis/internal/server/sysMonitor"
	"td27/rpc/basis/internal/server/sysTool"

	"td27/rpc/basis/types/basis_pb"
	"td27/rpc/basis/types/sysManagement/api_pb"
	"td27/rpc/basis/types/sysManagement/button_pb"
	"td27/rpc/basis/types/sysManagement/dept_pb"
	"td27/rpc/basis/types/sysManagement/dict_pb"
	"td27/rpc/basis/types/sysManagement/menu_pb"
	"td27/rpc/basis/types/sysManagement/permission_pb"
	"td27/rpc/basis/types/sysManagement/role_pb"
	"td27/rpc/basis/types/sysManagement/user_pb"
	"td27/rpc/basis/types/sysMonitor/dashboard_pb"
	"td27/rpc/basis/types/sysMonitor/operation_log_pb"
	"td27/rpc/basis/types/sysTool/cache_pb"
	"td27/rpc/basis/types/sysTool/cron_pb"
	"td27/rpc/basis/types/sysTool/file_pb"
	"td27/rpc/basis/types/sysTool/service_token_pb"
)

var configFile = flag.String("f", "etc/basis.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		basis_pb.RegisterBasisServer(grpcServer, basisServer.NewBasisServer(ctx))
		user_pb.RegisterUserServer(grpcServer, sysManagement.NewUserServer(ctx))
		role_pb.RegisterRoleServer(grpcServer, sysManagement.NewRoleServer(ctx))
		permission_pb.RegisterPermissionServer(grpcServer, sysManagement.NewPermissionServer(ctx))
		menu_pb.RegisterMenuServer(grpcServer, sysManagement.NewMenuServer(ctx))
		dept_pb.RegisterDeptServer(grpcServer, sysManagement.NewDeptServer(ctx))
		dict_pb.RegisterDictServer(grpcServer, sysManagement.NewDictServer(ctx))
		api_pb.RegisterAPIServer(grpcServer, sysManagement.NewAPIServer(ctx))
		button_pb.RegisterButtonServer(grpcServer, sysManagement.NewButtonServer(ctx))
		file_pb.RegisterFileServer(grpcServer, sysTool.NewFileServer(ctx))
		cron_pb.RegisterCronServer(grpcServer, sysTool.NewCronServer(ctx))
		cache_pb.RegisterCacheServer(grpcServer, sysTool.NewCacheServer(ctx))
		service_token_pb.RegisterServiceTokenServer(grpcServer, sysTool.NewServiceTokenServer(ctx))
		operation_log_pb.RegisterOperationLogServer(grpcServer, sysMonitor.NewOperationLogServer(ctx))
		dashboard_pb.RegisterDashboardServer(grpcServer, sysMonitor.NewDashboardServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	logx.Infof("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
