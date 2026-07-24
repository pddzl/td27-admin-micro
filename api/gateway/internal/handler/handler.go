package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"

	"td27/api/gateway/internal/handler/basis/sysManagement"
	"td27/api/gateway/internal/handler/basis/sysMonitor"
	"td27/api/gateway/internal/handler/basis/sysTool"
	"td27/api/gateway/internal/middleware"
	"td27/api/gateway/internal/svc"
)

func RegisterHandlers(server *rest.Server, svcCtx *svc.ServiceContext) {
	jwtMiddleware := middleware.NewJwtMiddleware(svcCtx)
	opRecordMiddleware := middleware.NewOperationRecordMiddleware(svcCtx)

	loginHandler := sysManagement.NewLoginHandler(svcCtx)
	captchaHandler := sysManagement.NewCaptchaHandler(svcCtx)
	logoutHandler := sysManagement.NewLogoutHandler(svcCtx)
	deptHandler := sysManagement.NewDeptHandler(svcCtx)
	dictHandler := sysManagement.NewDictHandler(svcCtx)
	dictDetailHandler := sysManagement.NewDictDetailHandler(svcCtx)
	apiHandler := sysManagement.NewAPIHandler(svcCtx)
	buttonHandler := sysManagement.NewButtonHandler(svcCtx)
	userHandler := sysManagement.NewUserHandler(svcCtx)
	roleHandler := sysManagement.NewRoleHandler(svcCtx)
	menuHandler := sysManagement.NewMenuHandler(svcCtx)
	permissionHandler := sysManagement.NewPermissionHandler(svcCtx)
	fileHandler := sysTool.NewFileHandler(svcCtx)
	cronHandler := sysTool.NewCronHandler(svcCtx)
	cacheHandler := sysTool.NewCacheHandler(svcCtx)
	serviceTokenHandler := sysTool.NewServiceTokenHandler(svcCtx)
	operationLogHandler := sysMonitor.NewOperationLogHandler(svcCtx)
	dashboardHandler := sysMonitor.NewDashboardHandler(svcCtx)

	// Public routes
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/health",
		Handler: loginHandler.Health,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/login",
		Handler: loginHandler.Login,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/captcha",
		Handler: captchaHandler.GenerateCaptcha,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/logout",
		Handler: logoutHandler.Logout,
	})

	// Dept routes
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/dept/list",
		Handler: jwtMiddleware.Handle(deptHandler.ListDept),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/dept/getElTreeDepts",
		Handler: jwtMiddleware.Handle(deptHandler.GetElTreeDepts),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/dept/getTree",
		Handler: jwtMiddleware.Handle(deptHandler.GetDeptTree),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/dept/:id",
		Handler: jwtMiddleware.Handle(deptHandler.GetDept),
	})
	//server.AddRoute(rest.Route{
	//	Method:  http.MethodPost,
	//	Path:    "/dept/descendants",
	//	Handler: jwtMiddleware.Handle(deptHandler.GetDeptDescendants),
	//})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/dept/create",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(deptHandler.CreateDept)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/dept/update",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(deptHandler.UpdateDept)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/dept/delete",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(deptHandler.DeleteDept)),
	})

	// Dict routes
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/dict/list",
		Handler: jwtMiddleware.Handle(dictHandler.ListDict),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/dict/:id",
		Handler: jwtMiddleware.Handle(dictHandler.GetDict),
	})
	//server.AddRoute(rest.Route{
	//	Method:  http.MethodGet,
	//	Path:    "/api/dict/en/:en_name",
	//	Handler: jwtMiddleware.Handle(dictHandler.GetDictByENName),
	//})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/dict/create",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(dictHandler.CreateDict)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/dict/update",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(dictHandler.UpdateDict)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/dict/delete",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(dictHandler.DeleteDict)),
	})

	// DictDetail route
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/dict-detail/create",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(dictDetailHandler.CreateDictDetail)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/dict-detail/update",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(dictDetailHandler.UpdateDictDetail)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/dict-detail/delete",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(dictDetailHandler.DeleteDictDetail)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/dict-detail/flat",
		Handler: jwtMiddleware.Handle(dictDetailHandler.FlatDictDetails),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/dict-detail/list",
		Handler: jwtMiddleware.Handle(dictDetailHandler.ListDictDetail),
	})

	// API routes
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/list",
		Handler: jwtMiddleware.Handle(apiHandler.ListAPI),
	})
	//server.AddRoute(rest.Route{
	//	Method:  http.MethodGet,
	//	Path:    "/api/apis/:id",
	//	Handler: jwtMiddleware.Handle(apiHandler.GetAPI),
	//})
	//server.AddRoute(rest.Route{
	//	Method:  http.MethodPost,
	//	Path:    "/api/apis/by-group",
	//	Handler: jwtMiddleware.Handle(apiHandler.GetAPIsByGroup),
	//})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/create",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(apiHandler.CreateAPI)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/update",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(apiHandler.UpdateAPI)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/delete",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(apiHandler.DeleteAPI)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/apis/delete-by-ids",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(apiHandler.DeleteByIds)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/elTree",
		Handler: jwtMiddleware.Handle(apiHandler.GetAPITree),
	})

	// Button routes
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/button/list",
		Handler: jwtMiddleware.Handle(buttonHandler.ListButton),
	})
	//server.AddRoute(rest.Route{
	//	Method:  http.MethodGet,
	//	Path:    "/api/button/:id",
	//	Handler: jwtMiddleware.Handle(buttonHandler.GetButton),
	//})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/button/by-page",
		Handler: jwtMiddleware.Handle(buttonHandler.GetButtonsByPagePath),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/button/user-buttons",
		Handler: jwtMiddleware.Handle(buttonHandler.GetUserButtons),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/button/create",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(buttonHandler.CreateButton)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/button/update",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(buttonHandler.UpdateButton)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/button/delete",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(buttonHandler.DeleteButton)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/button/batch-check",
		Handler: jwtMiddleware.Handle(buttonHandler.BatchCheckPermission),
	})

	// User routes
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/user/getUserInfo",
		Handler: jwtMiddleware.Handle(userHandler.GetUserInfo),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/user/list",
		Handler: jwtMiddleware.Handle(userHandler.ListUser),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/user/create",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(userHandler.CreateUser)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/user/update",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(userHandler.UpdateUser)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/user/delete",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(userHandler.DeleteUser)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/user/modify-password",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(userHandler.ModifyPassword)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/user/switch-active",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(userHandler.SwitchActive)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/user/assign-roles",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(userHandler.AssignRoles)),
	})

	// Role routes
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/role/info",
		Handler: jwtMiddleware.Handle(roleHandler.GetRole),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/role/children",
		Handler: jwtMiddleware.Handle(roleHandler.GetRoleWithChildren),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/role/list",
		Handler: jwtMiddleware.Handle(roleHandler.ListRole),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/role/all",
		Handler: jwtMiddleware.Handle(roleHandler.GetAllRoles),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/role/create",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(roleHandler.CreateRole)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/role/update",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(roleHandler.UpdateRole)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/role/delete",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(roleHandler.DeleteRole)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/role/assign-permissions",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(roleHandler.AssignPermissions)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/role/permissions",
		Handler: jwtMiddleware.Handle(roleHandler.GetRolePermissions),
	})

	// Menu routes
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/menu/list",
		Handler: jwtMiddleware.Handle(menuHandler.ListMenu),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/menu/info",
		Handler: jwtMiddleware.Handle(menuHandler.GetMenu),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/menu/elTree",
		Handler: jwtMiddleware.Handle(menuHandler.GetElTreeMenus),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/menu/user-menus",
		Handler: jwtMiddleware.Handle(menuHandler.GetUserMenus),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/menu/create",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(menuHandler.CreateMenu)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/menu/update",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(menuHandler.UpdateMenu)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/menu/delete",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(menuHandler.DeleteMenu)),
	})

	// Permission routes
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/permission/info",
		Handler: jwtMiddleware.Handle(permissionHandler.GetPermission),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/permission/list",
		Handler: jwtMiddleware.Handle(permissionHandler.ListPermission),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/permission/all",
		Handler: jwtMiddleware.Handle(permissionHandler.GetAllPermissions),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/permission/by-role",
		Handler: jwtMiddleware.Handle(permissionHandler.GetPermissionsByRoleId),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/permission/create",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(permissionHandler.CreatePermission)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/permission/update",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(permissionHandler.UpdatePermission)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/permission/delete",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(permissionHandler.DeletePermission)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/permission/check",
		Handler: jwtMiddleware.Handle(permissionHandler.CheckPermission),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/permission/reload-policy",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(permissionHandler.ReloadPolicy)),
	})

	// File routes
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/file/upload",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(fileHandler.UploadFile)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/file/:id",
		Handler: jwtMiddleware.Handle(fileHandler.GetFile),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/file/download/:id",
		Handler: jwtMiddleware.Handle(fileHandler.DownloadFile),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/file/list",
		Handler: jwtMiddleware.Handle(fileHandler.ListFile),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/file/delete",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(fileHandler.DeleteFile)),
	})

	// Cron routes
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/cron/:id",
		Handler: jwtMiddleware.Handle(cronHandler.GetCron),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/cron/list",
		Handler: jwtMiddleware.Handle(cronHandler.ListCron),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/cron/create",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(cronHandler.CreateCron)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/cron/update",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(cronHandler.UpdateCron)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/cron/toggle-status",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(cronHandler.ToggleCronStatus)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/cron/execute",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(cronHandler.ExecuteCronNow)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/cron/delete",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(cronHandler.DeleteCron)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/cron/delete-by-ids",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(cronHandler.DeleteByIds)),
	})

	// Cache routes
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/cache/get",
		Handler: jwtMiddleware.Handle(cacheHandler.GetCache),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/cache/set",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(cacheHandler.SetCache)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/cache/delete",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(cacheHandler.DeleteCache)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/cache/list",
		Handler: jwtMiddleware.Handle(cacheHandler.ListCache),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/cache/cleanup",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(cacheHandler.CleanupExpired)),
	})

	// Service token routes
	// ServiceToken - camelCase routes matching frontend
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/serviceToken/create",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(serviceTokenHandler.CreateServiceToken)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/serviceToken/list",
		Handler: jwtMiddleware.Handle(serviceTokenHandler.ListServiceToken),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/serviceToken/update",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(serviceTokenHandler.UpdateServiceToken)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/serviceToken/delete",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(serviceTokenHandler.DeleteServiceToken)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/service-token/:id",
		Handler: jwtMiddleware.Handle(serviceTokenHandler.GetServiceToken),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/service-token/toggle-status",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(serviceTokenHandler.ToggleTokenStatus)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/service-token/assign-permissions",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(serviceTokenHandler.AssignTokenPermissions)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/service-token/permissions/:id",
		Handler: jwtMiddleware.Handle(serviceTokenHandler.GetTokenPermissions),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/service-token/validate",
		Handler: jwtMiddleware.Handle(serviceTokenHandler.ValidateToken),
	})

	// Dashboard routes
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/dashboard/statistics",
		Handler: jwtMiddleware.Handle(dashboardHandler.GetStatistics),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/dashboard/recent-operations",
		Handler: jwtMiddleware.Handle(dashboardHandler.GetRecentOperations),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/dashboard/system-info",
		Handler: jwtMiddleware.Handle(dashboardHandler.GetSystemInfo),
	})

	// Operation log routes
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/operation-log/list",
		Handler: jwtMiddleware.Handle(operationLogHandler.ListOperationLog),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/opl/list",
		Handler: jwtMiddleware.Handle(operationLogHandler.ListOperationLog),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/operation-log/delete",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(operationLogHandler.DeleteOperationLog)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/opl/delete",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(operationLogHandler.DeleteOperationLog)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/operation-log/delete-by-ids",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(operationLogHandler.DeleteOperationLogByIds)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/opl/deleteByIds",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(operationLogHandler.DeleteOperationLogByIds)),
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/operation-log/cleanup",
		Handler: opRecordMiddleware.Handle(jwtMiddleware.Handle(operationLogHandler.CleanupExpiredLogs)),
	})
}
