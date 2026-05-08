package svc

import (
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"td27/rpc/basis/internal/config"
	"td27/rpc/basis/internal/initialization"
	authRepo "td27/rpc/basis/internal/repository/authority"
	monitorRepo "td27/rpc/basis/internal/repository/monitor"
	toolRepo "td27/rpc/basis/internal/repository/tool"
	authService "td27/rpc/basis/internal/service/authority"
	monitorService "td27/rpc/basis/internal/service/monitor"
	toolService "td27/rpc/basis/internal/service/tool"
)

type ServiceContext struct {
	Config         config.Config
	DB             *gorm.DB
	CasbinEnforcer *casbin.SyncedCachedEnforcer
	JWT            *JWTManager
	CronScheduler  *cron.Cron

	// Repositories
	UserRepo    authRepo.UserRepository
	RoleRepo    authRepo.RoleRepository
	PermRepo    authRepo.PermissionRepository
	MenuRepo    authRepo.MenuRepository
	DeptRepo    authRepo.DeptRepository
	DictRepo    authRepo.DictRepository
	DictDetRepo authRepo.DictDetailRepository
	APIRepo     authRepo.APIRepository
	ButtonRepo  authRepo.ButtonRepository

	FileRepo  toolRepo.FileRepository
	CronRepo  toolRepo.CronRepository
	CacheRepo toolRepo.CacheRepository
	TokenRepo toolRepo.ServiceTokenRepository

	LogRepo monitorRepo.OperationLogRepository

	// Services
	UserService   *authService.UserService
	RoleService   *authService.RoleService
	PermService   *authService.PermissionService
	MenuService   *authService.MenuService
	DeptService   *authService.DeptService
	DictService   *authService.DictService
	APIService    *authService.APIService
	ButtonService *authService.ButtonService

	FileService  *toolService.FileService
	CronService  *toolService.CronService
	CacheService *toolService.CacheService
	TokenService *toolService.ServiceTokenService

	LogService *monitorService.OperationLogService
}

// JWTManager handles JWT token operations
type JWTManager struct {
	SigningKey  []byte
	ExpiresTime int64
	BufferTime  int64
	Issuer      string
}

// NewJWTManager creates a new JWT manager instance
func NewJWTManager(cfg config.JWT) *JWTManager {
	return &JWTManager{
		SigningKey:  []byte(cfg.SigningKey),
		ExpiresTime: cfg.ExpiresTime,
		BufferTime:  cfg.BufferTime,
		Issuer:      cfg.Issuer,
	}
}

// getCasbinModel returns the Casbin RBAC model
func getCasbinModel(enableRoleHierarchy bool) (casbinModel.Model, error) {
	var modelText string

	if enableRoleHierarchy {
		modelText = `
		[request_definition]
		r = sub, obj, act

		[policy_definition]
		p = sub, obj, act

		[role_definition]
		g = _, _
		g2 = _, _

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == '*')
		`
	} else {
		modelText = `
		[request_definition]
		r = sub, obj, act

		[policy_definition]
		p = sub, obj, act

		[role_definition]
		g = _, _
		g2 = _, _

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == '*')
		`
	}

	return casbinModel.NewModelFromString(modelText)
}

func NewServiceContext(c config.Config) *ServiceContext {
	// Initialize database connection
	db, err := initialization.Gorm(c.Pgsql)
	if err != nil {
		logx.Errorf("init postgresql err: %v", err)
		panic(err)
	}
	logx.Infof("init postgresql success")

	// Initialize Casbin enforcer
	casbinModel, err := getCasbinModel(c.Casbin.EnableRoleHierarchy)
	if err != nil {
		logx.Errorf("init casbin model err: %v", err)
		return nil
	}

	// TODO: Implement PermissionAdapter when repository layer is done
	// For now, use a dummy adapter, will replace with proper adapter later
	casbinEnforcer, err := casbin.NewSyncedCachedEnforcer(casbinModel)
	if err != nil {
		logx.Errorf("init casbin enforcer err: %v", err)
		panic(err)
	}

	// Configure Casbin cache
	cacheTTL := c.Casbin.CacheTTL
	if cacheTTL <= 0 {
		cacheTTL = 3600
	}
	casbinEnforcer.SetExpireTime(time.Duration(cacheTTL) * time.Second)

	// Start auto load policy if enabled
	if c.Casbin.AutoLoadInterval > 0 {
		casbinEnforcer.StartAutoLoadPolicy(time.Duration(c.Casbin.AutoLoadInterval) * time.Second)
	}

	logx.Infof("init casbin enforcer success")

	// Initialize JWT manager
	jwtManager := NewJWTManager(c.JWT)
	logx.Infof("init jwt manager success")

	// Initialize cron scheduler
	cronScheduler := cron.New()
	cronScheduler.Start()
	logx.Infof("init cron scheduler success")

	// Initialize Repositories
	userRepo := authRepo.NewUserRepository(db)
	roleRepo := authRepo.NewRoleRepository(db)
	permRepo := authRepo.NewPermissionRepository(db)
	menuRepo := authRepo.NewMenuRepository(db)
	deptRepo := authRepo.NewDeptRepository(db)
	dictRepo := authRepo.NewDictRepository(db)
	dictDetRepo := authRepo.NewDictDetailRepository(db)
	apiRepo := authRepo.NewAPIRepository(db)
	buttonRepo := authRepo.NewButtonRepository(db)

	fileRepo := toolRepo.NewFileRepository(db)
	cronRepo := toolRepo.NewCronRepository(db)
	cacheRepo := toolRepo.NewCacheRepository(db)
	tokenRepo := toolRepo.NewServiceTokenRepository(db)

	logRepo := monitorRepo.NewOperationLogRepository(db)

	// Initialize Services
	userService := authService.NewUserService(userRepo, roleRepo)
	roleService := authService.NewRoleService(roleRepo, permRepo, userRepo)
	permService := authService.NewPermissionService(permRepo, roleRepo, casbinEnforcer)
	menuService := authService.NewMenuService(menuRepo)
	deptService := authService.NewDeptService(deptRepo)
	dictService := authService.NewDictService(dictRepo)
	apiService := authService.NewAPIService(apiRepo)
	buttonService := authService.NewButtonService(buttonRepo)

	fileService := toolService.NewFileService(fileRepo, c.File)
	cronService := toolService.NewCronService(cronRepo, cronScheduler)
	cacheService := toolService.NewCacheService(cacheRepo)
	tokenService := toolService.NewServiceTokenService(tokenRepo)

	logService := monitorService.NewOperationLogService(logRepo)

	return &ServiceContext{
		Config:         c,
		DB:             db,
		CasbinEnforcer: casbinEnforcer,
		JWT:            jwtManager,
		CronScheduler:  cronScheduler,

		// Repositories
		UserRepo:    userRepo,
		RoleRepo:    roleRepo,
		PermRepo:    permRepo,
		MenuRepo:    menuRepo,
		DeptRepo:    deptRepo,
		DictRepo:    dictRepo,
		DictDetRepo: dictDetRepo,
		APIRepo:     apiRepo,
		ButtonRepo:  buttonRepo,

		FileRepo:  fileRepo,
		CronRepo:  cronRepo,
		CacheRepo: cacheRepo,
		TokenRepo: tokenRepo,

		LogRepo: logRepo,

		// Services
		UserService:   userService,
		RoleService:   roleService,
		PermService:   permService,
		MenuService:   menuService,
		DeptService:   deptService,
		DictService:   dictService,
		APIService:    apiService,
		ButtonService: buttonService,

		FileService:  fileService,
		CronService:  cronService,
		CacheService: cacheService,
		TokenService: tokenService,

		LogService: logService,
	}
}
