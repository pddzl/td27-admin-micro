package svc

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"td27/rpc/basis/internal/config"
	"td27/rpc/basis/internal/initialization"
	sysManagementRepo "td27/rpc/basis/internal/repository/sysManagement"
	sysMonitorRepo "td27/rpc/basis/internal/repository/sysMonitor"
	sysToolRepo "td27/rpc/basis/internal/repository/sysTool"
	sysManagementService "td27/rpc/basis/internal/service/sysManagement"
	sysMonitorService "td27/rpc/basis/internal/service/sysMonitor"
	sysToolService "td27/rpc/basis/internal/service/sysTool"
)

type ServiceContext struct {
	Config         config.Config
	DB             *gorm.DB
	CasbinEnforcer *casbin.SyncedCachedEnforcer
	JWT            *JWTManager
	CronScheduler  *cron.Cron

	// Repositories
	UserRepo    sysManagementRepo.UserRepository
	RoleRepo    sysManagementRepo.RoleRepository
	PermRepo    sysManagementRepo.PermissionRepository
	MenuRepo    sysManagementRepo.MenuRepository
	DeptRepo    sysManagementRepo.DeptRepository
	DictRepo    sysManagementRepo.DictRepository
	DictDetRepo sysManagementRepo.DictDetailRepository
	APIRepo     sysManagementRepo.APIRepository
	ButtonRepo  sysManagementRepo.ButtonRepository

	FileRepo  sysToolRepo.FileRepository
	CronRepo  sysToolRepo.CronRepository
	CacheRepo sysToolRepo.CacheRepository
	TokenRepo sysToolRepo.ServiceTokenRepository

	LogRepo       sysMonitorRepo.OperationLogRepository
	DashboardRepo sysMonitorRepo.DashboardRepository

	// Services
	UserService   *sysManagementService.UserService
	RoleService   *sysManagementService.RoleService
	PermService   *sysManagementService.PermissionService
	MenuService   *sysManagementService.MenuService
	DeptService   *sysManagementService.DeptService
	DictService   *sysManagementService.DictService
	APIService    *sysManagementService.APIService
	ButtonService *sysManagementService.ButtonService

	FileService  *sysToolService.FileService
	CronService  *sysToolService.CronService
	CacheService *sysToolService.CacheService
	TokenService *sysToolService.ServiceTokenService

	LogService       *sysMonitorService.OperationLogService
	DashboardService *sysMonitorService.DashboardService
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

// CreateToken generates a JWT token with the given user claims
func (m *JWTManager) CreateToken(userId uint64, username string, roleIds []uint64) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(m.ExpiresTime) * time.Second)

	claims := jwt.MapClaims{
		"userId":   userId,
		"username": username,
		"roleIds":  roleIds,
		"exp":      expiresAt.Unix(),
		"iat":      now.Unix(),
		"iss":      m.Issuer,
		"nbf":      now.Unix() - 1,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.SigningKey)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return tokenString, nil
}

// ParseToken parses and validates a JWT token string
func (m *JWTManager) ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
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
	userRepo := sysManagementRepo.NewUserRepository(db)
	roleRepo := sysManagementRepo.NewRoleRepository(db)
	permRepo := sysManagementRepo.NewPermissionRepository(db)
	menuRepo := sysManagementRepo.NewMenuRepository(db)
	deptRepo := sysManagementRepo.NewDeptRepository(db)
	dictRepo := sysManagementRepo.NewDictRepository(db)
	dictDetRepo := sysManagementRepo.NewDictDetailRepository(db)
	apiRepo := sysManagementRepo.NewAPIRepository(db)
	buttonRepo := sysManagementRepo.NewButtonRepository(db)

	fileRepo := sysToolRepo.NewFileRepository(db)
	cronRepo := sysToolRepo.NewCronRepository(db)
	cacheRepo := sysToolRepo.NewCacheRepository(db)
	tokenRepo := sysToolRepo.NewServiceTokenRepository(db)

	logRepo := sysMonitorRepo.NewOperationLogRepository(db)
	dashboardRepo := sysMonitorRepo.NewDashboardRepository(db)

	// Initialize Services
	userService := sysManagementService.NewUserService(userRepo, roleRepo)
	roleService := sysManagementService.NewRoleService(roleRepo, permRepo, userRepo)
	permService := sysManagementService.NewPermissionService(permRepo, roleRepo, casbinEnforcer)
	menuService := sysManagementService.NewMenuService(menuRepo)
	deptService := sysManagementService.NewDeptService(deptRepo)
	dictService := sysManagementService.NewDictService(dictRepo)
	apiService := sysManagementService.NewAPIService(apiRepo)
	buttonService := sysManagementService.NewButtonService(buttonRepo)

	fileService := sysToolService.NewFileService(fileRepo, c.File)
	cronService := sysToolService.NewCronService(cronRepo, cronScheduler)
	cacheService := sysToolService.NewCacheService(cacheRepo)
	tokenService := sysToolService.NewServiceTokenService(tokenRepo)

	logService := sysMonitorService.NewOperationLogService(logRepo)
	dashboardService := sysMonitorService.NewDashboardService(dashboardRepo)

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

		LogRepo:       logRepo,
		DashboardRepo: dashboardRepo,

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

		LogService:       logService,
		DashboardService: dashboardService,
	}
}
