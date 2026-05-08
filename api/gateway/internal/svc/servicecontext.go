package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"td27/api/gateway/internal/config"
	"td27/rpc/basis/types/authority/api_pb"
	"td27/rpc/basis/types/authority/button_pb"
	"td27/rpc/basis/types/authority/dept_pb"
	"td27/rpc/basis/types/authority/dict_pb"
	"td27/rpc/basis/types/authority/menu_pb"
	"td27/rpc/basis/types/authority/permission_pb"
	"td27/rpc/basis/types/authority/role_pb"
	"td27/rpc/basis/types/authority/user_pb"
	"td27/rpc/basis/types/monitor/operation_log_pb"
	"td27/rpc/basis/types/tool/cache_pb"
	"td27/rpc/basis/types/tool/cron_pb"
	"td27/rpc/basis/types/tool/file_pb"
	"td27/rpc/basis/types/tool/service_token_pb"
)

type ServiceContext struct {
	Config config.Config

	UserClient         user_pb.UserClient
	RoleClient         role_pb.RoleClient
	PermissionClient   permission_pb.PermissionClient
	MenuClient         menu_pb.MenuClient
	DeptClient         dept_pb.DeptClient
	DictClient         dict_pb.DictClient
	APIClient          api_pb.APIClient
	ButtonClient       button_pb.ButtonClient
	FileClient         file_pb.FileClient
	CronClient         cron_pb.CronClient
	CacheClient        cache_pb.CacheClient
	ServiceTokenClient service_token_pb.ServiceTokenClient
	OperationLogClient operation_log_pb.OperationLogClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := zrpc.MustNewClient(c.BasisRpc)

	return &ServiceContext{
		Config:             c,
		UserClient:         user_pb.NewUserClient(conn.Conn()),
		RoleClient:         role_pb.NewRoleClient(conn.Conn()),
		PermissionClient:   permission_pb.NewPermissionClient(conn.Conn()),
		MenuClient:         menu_pb.NewMenuClient(conn.Conn()),
		DeptClient:         dept_pb.NewDeptClient(conn.Conn()),
		DictClient:         dict_pb.NewDictClient(conn.Conn()),
		APIClient:          api_pb.NewAPIClient(conn.Conn()),
		ButtonClient:       button_pb.NewButtonClient(conn.Conn()),
		FileClient:         file_pb.NewFileClient(conn.Conn()),
		CronClient:         cron_pb.NewCronClient(conn.Conn()),
		CacheClient:        cache_pb.NewCacheClient(conn.Conn()),
		ServiceTokenClient: service_token_pb.NewServiceTokenClient(conn.Conn()),
		OperationLogClient: operation_log_pb.NewOperationLogClient(conn.Conn()),
	}
}
