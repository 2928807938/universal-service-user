package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// API层
	configHTTP "github.com/2928807938/universal-service-user/api/config-api/http"
	configAppService "github.com/2928807938/universal-service-user/api/config-api/service"
	userHTTP "github.com/2928807938/universal-service-user/api/user-api/http"
	userAppService "github.com/2928807938/universal-service-user/api/user-api/service"

	// Auth模块
	authDomain "github.com/2928807938/universal-service-user/auth/domain/service"
	authEntity "github.com/2928807938/universal-service-user/auth/infrastructure/entity"
	authJWT "github.com/2928807938/universal-service-user/auth/infrastructure/jwt"
	authRepo "github.com/2928807938/universal-service-user/auth/infrastructure/repository"

	// Hook模块
	"github.com/2928807938/universal-service-user/hook"

	// Notification模块
	notificationDomain "github.com/2928807938/universal-service-user/notification/domain/service"
	notificationEmail "github.com/2928807938/universal-service-user/notification/infrastructure/email"

	// OAuth模块
	oauthEntity "github.com/2928807938/universal-service-user/oauth/infrastructure/entity"

	// ConfigCenter模块
	configDomainRepo "github.com/2928807938/universal-service-user/configcenter/domain/repository"
	configDomainService "github.com/2928807938/universal-service-user/configcenter/domain/service"
	configInfraEntity "github.com/2928807938/universal-service-user/configcenter/infrastructure/entity"
	configInfraRepo "github.com/2928807938/universal-service-user/configcenter/infrastructure/repository"

	// Share模块
	"github.com/2928807938/universal-service-user/share/config"
	"github.com/2928807938/universal-service-user/share/errors"
	shareRedis "github.com/2928807938/universal-service-user/share/redis"
	"github.com/2928807938/universal-service-user/share/types"

	userDomainRepo "github.com/2928807938/universal-service-user/user/domain/repository"
	// User模块
	userDomain "github.com/2928807938/universal-service-user/user/domain/service"
	userInfraEntity "github.com/2928807938/universal-service-user/user/infrastructure/entity"
	userInfraRepo "github.com/2928807938/universal-service-user/user/infrastructure/repository"

	// Verification模块
	verificationDomain "github.com/2928807938/universal-service-user/verification/domain/service"
	verificationInfraRepo "github.com/2928807938/universal-service-user/verification/infrastructure/repository"
)

func main() {
	log.Println("========== 通用用户服务启动中 ==========")

	// 1. 加载配置
	log.Println("[1/8] 加载配置文件...")
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	log.Printf("✓ 配置加载成功 (数据库:%s, Redis:%s:%d)",
		cfg.Database.Driver, cfg.Redis.Host, cfg.Redis.Port)

	// 2. 初始化数据库
	log.Println("[2/8] 初始化数据库...")
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	log.Println("✓ 数据库连接成功")

	// 3. 自动建表
	if cfg.Database.AutoMigrate {
		log.Println("[3/8] 执行数据库迁移...")
		if err := autoMigrate(db); err != nil {
			log.Fatalf("数据库迁移失败: %v", err)
		}
		log.Println("✓ 数据库迁移完成")
	} else {
		log.Println("[3/8] 跳过数据库迁移(auto_migrate=false)")
	}

	// 4. 初始化 Redis
	log.Println("[4/8] 初始化 Redis...")
	rdb, err := initRedis(cfg)
	if err != nil {
		log.Fatalf("初始化 Redis 失败: %v", err)
	}
	// 测试 Redis 连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis 连接测试失败: %v", err)
	}
	log.Println("✓ Redis 连接成功")

	// 5. 初始化领域服务
	log.Println("[5/8] 初始化领域服务...")
	domainServices := initDomainServices(cfg, db, rdb)
	log.Println("✓ 领域服务初始化完成")

	// 6. 初始化应用服务
	log.Println("[6/8] 初始化应用服务...")
	appServices := initAppServices(cfg, domainServices)
	log.Println("✓ 应用服务初始化完成")

	// 7. 初始化 HTTP 服务器
	log.Println("[7/8] 初始化 HTTP 服务器...")
	h := initHTTPServer(cfg, appServices, domainServices)
	log.Println("✓ HTTP 服务器初始化完成")

	// 8. 启动服务
	log.Printf("[8/8] 服务启动在端口 %s", cfg.Server.Port)
	log.Println("========================================")
	log.Printf("✓ 服务已启动: http://localhost:%s", cfg.Server.Port)
	log.Printf("✓ 健康检查: http://localhost:%s/health", cfg.Server.Port)
	log.Printf("✓ API 地址: http://localhost:%s/api/v1", cfg.Server.Port)
	log.Println("========================================")

	h.Spin()
}

// initDB 初始化数据库连接
func initDB(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Database.Driver {
	case "postgres":
		dialector = postgres.Open(cfg.Database.GetDSN())
	case "mysql":
		dialector = mysql.Open(cfg.Database.GetDSN())
	case "sqlite":
		dialector = sqlite.Open(cfg.Database.GetDSN())
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", cfg.Database.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	return db, nil
}

// autoMigrate 自动迁移数据库表
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// User 模块表
		&userInfraEntity.UserPO{},
		&userInfraEntity.UserProfilePO{},
		&userInfraEntity.UserLoginLogPO{},

		// Auth 模块表
		&authEntity.SessionPO{},

		// OAuth 模块表
		&oauthEntity.OAuthBindingPO{},

		// ConfigCenter 模块表
		&configInfraEntity.AppPO{},
		&configInfraEntity.AppConfigPO{},
		&configInfraEntity.ConfigHistoryPO{},
	)
}

// initRedis 初始化 Redis 连接
func initRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	return rdb, nil
}

// DomainServices 领域服务集合
type DomainServices struct {
	AppRepository        configDomainRepo.AppRepository
	AppConfigRepository  configDomainRepo.AppConfigRepository
	ConfigHistoryRepo    configDomainRepo.ConfigHistoryRepository
	UserRepo             userDomainRepo.UserRepository
	UserDomainService    *userDomain.UserDomainService
	LoginAttemptService  *userDomain.LoginAttemptService
	VerificationService  *verificationDomain.VerificationService
	NotificationService  *notificationDomain.NotificationService
	AuthService          *authDomain.AuthService
	ConfigCacheService   *configDomainService.ConfigCacheService
}

// initDomainServices 初始化领域服务
func initDomainServices(cfg *config.Config, db *gorm.DB, rdb *redis.Client) *DomainServices {
	// ConfigCenter repositories
	appRepo := configInfraRepo.NewAppRepositoryImpl(db)
	appConfigRepo := configInfraRepo.NewAppConfigRepositoryImpl(db)
	configHistoryRepo := configInfraRepo.NewConfigHistoryRepositoryImpl(db)
	configCacheService := configDomainService.NewConfigCacheService(appConfigRepo)

	// 1. User 仓储
	userRepo := userInfraRepo.NewUserRepositoryImpl(db)

	// 2. 创建登录尝试服务（防止暴力破解）
	// 将 go-redis 客户端包装为项目专用的 Redis 客户端
	wrappedRedisClient := &shareRedis.Client{
		Client: rdb,
	}
	var loginAttemptService *userDomain.LoginAttemptService
	if cfg.LoginRateLimit.Enabled {
		loginAttemptService = userDomain.NewLoginAttemptService(
			wrappedRedisClient,
			cfg.LoginRateLimit.Account.MaxFailures,
			time.Duration(cfg.LoginRateLimit.Account.LockDuration)*time.Second,
		)
		log.Printf("✓ 登录防暴力破解已启用 (最大失败次数: %d, 锁定时长: %ds)",
			cfg.LoginRateLimit.Account.MaxFailures,
			cfg.LoginRateLimit.Account.LockDuration)
	} else {
		log.Println("⚠ 登录防暴力破解未启用")
	}

	// 3. User 领域服务
	userDomainSvc := userDomain.NewUserDomainService(userRepo, loginAttemptService)

	// 4. Verification 仓储
	verificationRepo := verificationInfraRepo.NewRedisVerificationRepository(rdb)

	// 5. Verification 领域服务
	verificationConfig := &verificationDomain.VerificationConfig{
		CodeLength:     cfg.Verification.CodeLength,
		ExpireDuration: time.Duration(cfg.Verification.Expire) * time.Second,
		SendInterval:   time.Duration(cfg.Verification.RateLimit) * time.Second,
	}
	verificationSvc := verificationDomain.NewVerificationService(verificationRepo, verificationConfig)

	// 6. Notification 服务
	var emailProvider *notificationEmail.SMTPProvider
	if cfg.Email.Enabled {
		emailCfg := &notificationEmail.SMTPConfig{
			Host:      cfg.Email.SMTP.Host,
			Port:      cfg.Email.SMTP.Port,
			Username:  cfg.Email.SMTP.Username,
			Password:  cfg.Email.SMTP.Password,
			From:      cfg.Email.SMTP.From,
			Templates: cfg.Email.Templates,
		}
		var err error
		emailProvider, err = notificationEmail.NewSMTPProvider(emailCfg)
		if err != nil {
			log.Printf("警告: 邮件服务初始化失败: %v", err)
		}
	}

	notificationSvc := notificationDomain.NewNotificationService(emailProvider, nil)

	// 7. Auth 服务
	sessionRepo := authRepo.NewSessionRepositoryImpl(db)
	jwtConfig := &authJWT.Config{
		Secret:             cfg.JWT.Secret,
		AccessTokenExpire:  cfg.JWT.GetAccessTokenDuration(),
		RefreshTokenExpire: cfg.JWT.GetRefreshTokenDuration(),
		Issuer:             cfg.JWT.Issuer,
		EnableBlacklist:    true,
	}
	blacklist := authJWT.NewRedisBlacklist(rdb)
	jwtProvider, err := authJWT.NewProvider(jwtConfig, blacklist)
	if err != nil {
		log.Fatalf("初始化 JWT Provider 失败: %v", err)
	}
	authSvc := authDomain.NewAuthService(jwtProvider, sessionRepo, userDomainSvc)

	return &DomainServices{
		AppRepository:       appRepo,
		AppConfigRepository: appConfigRepo,
		ConfigHistoryRepo:   configHistoryRepo,
		UserRepo:            userRepo,
		UserDomainService:   userDomainSvc,
		LoginAttemptService: loginAttemptService,
		VerificationService: verificationSvc,
		NotificationService: notificationSvc,
		AuthService:         authSvc,
		ConfigCacheService:  configCacheService,
	}
}

// AppServices 应用服务集合
type AppServices struct {
	AppAppService          *configAppService.AppService
	UserAppService         *userAppService.UserAppService
	AuthAppService         *userAppService.AuthAppService
	VerificationAppService *userAppService.VerificationAppService
}

// initAppServices 初始化应用服务
func initAppServices(cfg *config.Config, ds *DomainServices) *AppServices {
	appAppSvc := configAppService.NewAppService(ds.AppRepository)

	userAppSvc := userAppService.NewUserAppService(
		ds.UserRepo,
		ds.VerificationService,
		ds.LoginAttemptService,
	)

	authAppSvc := userAppService.NewAuthAppService(
		cfg,
		ds.UserDomainService,
		ds.VerificationService,
		ds.NotificationService,
		ds.AuthService,
	)

	verificationAppSvc := userAppService.NewVerificationAppService(
		ds.VerificationService,
		ds.NotificationService,
	)

	return &AppServices{
		AppAppService:          appAppSvc,
		UserAppService:         userAppSvc,
		AuthAppService:         authAppSvc,
		VerificationAppService: verificationAppSvc,
	}
}

// initHTTPServer 初始化 HTTP 服务器
func initHTTPServer(cfg *config.Config, appServices *AppServices, ds *DomainServices) *server.Hertz {
	// 创建 Hertz 服务器
	h := server.Default(server.WithHostPorts(":" + cfg.Server.Port))

	// 全局中间件
	h.Use(corsMiddleware())
	h.Use(loggerMiddleware())

	// 健康检查
	h.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, types.Success(map[string]interface{}{
			"status":  "ok",
			"version": "1.0.0",
			"time":    time.Now().Format(time.RFC3339),
		}))
	})

	// API 路由
	api := h.Group("/api/v1")
	api.Use(tenantMiddleware(ds.ConfigCacheService))
	{
		// 应用注册路由
		appHandler := configHTTP.NewAppHandler(appServices.AppAppService)
		apps := api.Group("/apps")
		{
			apps.POST("/register", appHandler.Register)
		}

		// 认证路由
		authHandler := userHTTP.NewAuthHandler(appServices.AuthAppService)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)          // 登录
			auth.POST("/logout", authHandler.Logout)        // 登出
			auth.POST("/refresh", authHandler.RefreshToken) // 刷新令牌
		}

		// 用户路由
		userHandler := userHTTP.NewUserHandler(appServices.UserAppService)
		users := api.Group("/users")
		{
			users.POST("/register", userHandler.Register)           // 用户注册
			users.GET("/:id", userHandler.GetUser)                  // 获取用户
			users.PUT("/:id", userHandler.UpdateUser)               // 更新用户
			users.POST("/password/change", userHandler.ChangePassword) // 修改密码
			users.POST("/password/reset", userHandler.ResetPassword)   // 重置密码(忘记密码)
		}

		// 验证码路由
		verificationHandler := userHTTP.NewVerificationHandler(
			appServices.VerificationAppService,
		)
		verification := api.Group("/verification")
		{
			verification.POST("/code/send", verificationHandler.SendCode)     // 发送验证码
			verification.POST("/code/verify", verificationHandler.VerifyCode) // 验证验证码
		}
	}

	// 注册钩子(示例)
	hook.Register(hook.AfterCreate, func(ctx *hook.Context) error {
		if ctx.User != nil {
			log.Printf("钩子触发: 用户创建 - UserID: %d", ctx.User.ID)
		}
		return nil
	})

	return h
}

// corsMiddleware CORS 中间件
func corsMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-Id, X-App-Environment")
		c.Header("Access-Control-Max-Age", "86400")

		if string(c.Method()) == "OPTIONS" {
			c.AbortWithStatus(consts.StatusNoContent)
			return
		}

		c.Next(ctx)
	}
}

// loggerMiddleware 日志中间件
func loggerMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		c.Next(ctx)
		duration := time.Since(start)

		log.Printf("[%s] %s %s - %d (%v)",
			c.Method(),
			c.Path(),
			c.ClientIP(),
			c.Response.StatusCode(),
			duration,
		)
	}
}

// tenantMiddleware validates tenant headers and loads config.
func tenantMiddleware(cfgService *configDomainService.ConfigCacheService) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if cfgService == nil {
			errors.HandleError(ctx, c, errors.ErrInternal("config service not initialized", fmt.Errorf("nil config cache service")))
			return
		}
		path := string(c.Path())
		if path == "/api/v1/apps/register" {
			c.Next(ctx)
			return
		}

		tenantID := string(c.GetHeader("X-Tenant-Id"))
		if tenantID == "" {
			errors.HandleError(ctx, c, errors.ErrBadRequest("X-Tenant-Id is required"))
			return
		}

		environment := string(c.GetHeader("X-App-Environment"))
		if strings.TrimSpace(environment) == "" {
			environment = "prod"
		}

		cached, err := cfgService.GetConfig(ctx, tenantID, environment)
		if err != nil {
			errors.HandleError(ctx, c, err)
			return
		}

		c.Set("tenant_id", tenantID)
		c.Set("tenant_env", environment)
		c.Set("tenant_config", cached)
		c.Next(ctx)
	}
}
