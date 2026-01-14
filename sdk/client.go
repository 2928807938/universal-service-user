package sdk

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	// API 层
	"universal-service-user/api/user-api/dto/request"
	userAppService "universal-service-user/api/user-api/service"

	// Auth 模块
	authDomain "universal-service-user/auth/domain/service"
	authEntity "universal-service-user/auth/infrastructure/entity"
	authJWT "universal-service-user/auth/infrastructure/jwt"
	authRepo "universal-service-user/auth/infrastructure/repository"

	// Hook 模块
	"universal-service-user/hook"

	// Notification 模块
	notificationDomain "universal-service-user/notification/domain/service"
	notificationEmail "universal-service-user/notification/infrastructure/email"

	// OAuth 模块
	oauthEntity "universal-service-user/oauth/infrastructure/entity"

	// Share 模块
	"universal-service-user/share/config"
	shareRedis "universal-service-user/share/redis"

	// User 模块
	userDomain "universal-service-user/user/domain/service"
	userInfraEntity "universal-service-user/user/infrastructure/entity"
	userInfraRepo "universal-service-user/user/infrastructure/repository"

	// Verification 模块
	verificationEnum "universal-service-user/verification/domain/enum"
	verificationDomain "universal-service-user/verification/domain/service"
	verificationInfraRepo "universal-service-user/verification/infrastructure/repository"
)

// Client SDK 客户端
// 提供统一的用户服务接口
type Client struct {
	// 基础依赖
	db    *gorm.DB
	redis *redis.Client

	// 配置
	config      *config.Config
	autoMigrate bool

	// 邮箱和短信配置
	emailConfig *notificationEmail.SMTPConfig
	smsProvider string
	smsConfig   interface{}

	// 日志
	logger    Logger
	enableLog bool

	// 应用服务
	authAppService *userAppService.AuthAppService
	userAppService *userAppService.UserAppService

	// 领域服务（暴露给钩子使用）
	userDomainService    *userDomain.UserDomainService
	loginAttemptService  *userDomain.LoginAttemptService
	verificationService  *verificationDomain.VerificationService
	notificationService  *notificationDomain.NotificationService
	authService          *authDomain.AuthService

	// 钩子执行器
	hookExecutor *hook.Executor
	hookRegistry *hook.Registry
}

// defaultLogger 默认日志实现
type defaultLogger struct{}

func (l *defaultLogger) Info(msg string, fields ...interface{}) {
	log.Printf("[SDK INFO] %s %v", msg, fields)
}

func (l *defaultLogger) Error(msg string, fields ...interface{}) {
	log.Printf("[SDK ERROR] %s %v", msg, fields)
}

func (l *defaultLogger) Debug(msg string, fields ...interface{}) {
	log.Printf("[SDK DEBUG] %s %v", msg, fields)
}

func (l *defaultLogger) Warn(msg string, fields ...interface{}) {
	log.Printf("[SDK WARN] %s %v", msg, fields)
}

// New 创建新的 SDK 客户端
// 参数:
//
//	opts: 配置选项，使用 WithXXX 函数进行配置
//
// 返回:
//
//	*Client: SDK 客户端实例
//	error: 初始化过程中的错误
//
// 示例:
//
//	client, err := sdk.New(
//	    sdk.WithDatabase(db),
//	    sdk.WithRedis(rdb),
//	    sdk.WithAutoMigrate(true),
//	)
func New(opts ...Option) (*Client, error) {
	// 创建客户端实例
	c := &Client{
		autoMigrate:  true,
		logger:       &defaultLogger{},
		enableLog:    false,
		hookRegistry: hook.NewRegistry(),
	}

	// 应用配置选项
	for _, opt := range opts {
		opt(c)
	}

	// 验证必需的配置
	if err := c.validate(); err != nil {
		return nil, fmt.Errorf("SDK 配置验证失败: %w", err)
	}

	// 设置默认配置
	c.setDefaults()

	// 初始化钩子执行器
	c.hookExecutor = hook.NewExecutor(c.hookRegistry)
	if c.enableLog {
		c.hookExecutor.EnableLog(true)
	}

	// 自动建表
	if c.autoMigrate {
		if err := c.migrate(); err != nil {
			return nil, fmt.Errorf("数据库迁移失败: %w", err)
		}
		if c.enableLog {
			c.logger.Info("数据库迁移完成")
		}
	}

	// 初始化领域服务
	if err := c.initDomainServices(); err != nil {
		return nil, fmt.Errorf("初始化领域服务失败: %w", err)
	}

	// 初始化应用服务
	if err := c.initAppServices(); err != nil {
		return nil, fmt.Errorf("初始化应用服务失败: %w", err)
	}

	if c.enableLog {
		c.logger.Info("SDK 初始化完成")
	}

	return c, nil
}

// validate 验证必需的配置
func (c *Client) validate() error {
	if c.db == nil {
		return fmt.Errorf("数据库连接不能为空，请使用 WithDatabase() 配置")
	}
	if c.redis == nil {
		return fmt.Errorf("Redis 连接不能为空，请使用 WithRedis() 配置")
	}
	return nil
}

// setDefaults 设置默认配置
func (c *Client) setDefaults() {
	if c.config == nil {
		c.config = &config.Config{}
	}

	// JWT 默认配置
	if c.config.JWT.Secret == "" {
		c.config.JWT.Secret = "default-jwt-secret-please-change-in-production"
		if c.enableLog {
			c.logger.Warn("使用默认 JWT Secret，生产环境请务必修改")
		}
	}
	if c.config.JWT.Expire == 0 {
		c.config.JWT.Expire = 7200 // 2 小时
	}
	if c.config.JWT.RefreshExpire == 0 {
		c.config.JWT.RefreshExpire = 604800 // 7 天
	}
	if c.config.JWT.Issuer == "" {
		c.config.JWT.Issuer = "user-service-sdk"
	}

	// 验证码默认配置
	if c.config.Verification.CodeLength == 0 {
		c.config.Verification.CodeLength = 6
	}
	if c.config.Verification.Expire == 0 {
		c.config.Verification.Expire = 300 // 5 分钟
	}
	if c.config.Verification.RateLimit == 0 {
		c.config.Verification.RateLimit = 60 // 60 秒
	}
}

// migrate 执行数据库迁移
func (c *Client) migrate() error {
	return c.db.AutoMigrate(
		// User 模块表
		&userInfraEntity.UserPO{},
		&userInfraEntity.UserProfilePO{},
		&userInfraEntity.UserLoginLogPO{},

		// Auth 模块表
		&authEntity.SessionPO{},

		// OAuth 模块表
		&oauthEntity.OAuthBindingPO{},
	)
}

// initDomainServices 初始化领域服务
func (c *Client) initDomainServices() error {
	// 1. User 仓储
	userRepo := userInfraRepo.NewUserRepositoryImpl(c.db)

	// 2. LoginAttempt 服务（用于登录防刷）
	// 需要包装 Redis 客户端为项目专用类型
	wrappedRedis := &shareRedis.Client{
		Client: c.redis,
	}
	c.loginAttemptService = userDomain.NewLoginAttemptService(
		wrappedRedis,
		5,              // 最大失败次数
		15*time.Minute, // 锁定时长
	)

	// 3. User 领域服务
	c.userDomainService = userDomain.NewUserDomainService(userRepo, c.loginAttemptService)

	// 3. Verification 仓储
	verificationRepo := verificationInfraRepo.NewRedisVerificationRepository(c.redis)

	// 4. Verification 领域服务
	verificationConfig := &verificationDomain.VerificationConfig{
		CodeLength:     c.config.Verification.CodeLength,
		ExpireDuration: time.Duration(c.config.Verification.Expire) * time.Second,
		SendInterval:   time.Duration(c.config.Verification.RateLimit) * time.Second,
	}
	c.verificationService = verificationDomain.NewVerificationService(verificationRepo, verificationConfig)

	// 5. Notification 服务
	var emailProvider *notificationEmail.SMTPProvider
	if c.emailConfig != nil {
		var err error
		emailProvider, err = notificationEmail.NewSMTPProvider(c.emailConfig)
		if err != nil {
			return fmt.Errorf("初始化邮件服务失败: %w", err)
		}
	}

	// TODO: 支持短信服务提供者
	c.notificationService = notificationDomain.NewNotificationService(emailProvider, nil)

	// 6. Auth 服务
	sessionRepo := authRepo.NewSessionRepositoryImpl(c.db)
	jwtConfig := &authJWT.Config{
		Secret:             c.config.JWT.Secret,
		AccessTokenExpire:  c.config.JWT.GetAccessTokenDuration(),
		RefreshTokenExpire: c.config.JWT.GetRefreshTokenDuration(),
		Issuer:             c.config.JWT.Issuer,
	}
	blacklist := authJWT.NewRedisBlacklist(c.redis)
	jwtProvider, err := authJWT.NewProvider(jwtConfig, blacklist)
	if err != nil {
		return fmt.Errorf("初始化 JWT Provider 失败: %w", err)
	}
	c.authService = authDomain.NewAuthService(jwtProvider, sessionRepo, c.userDomainService)

	return nil
}

// initAppServices 初始化应用服务
func (c *Client) initAppServices() error {
	// 获取 User 仓储（重新创建以确保一致性）
	userRepo := userInfraRepo.NewUserRepositoryImpl(c.db)

	// 1. User 应用服务
	c.userAppService = userAppService.NewUserAppService(
		userRepo,
		c.verificationService,
		c.loginAttemptService,
	)

	// 2. Auth 应用服务
	c.authAppService = userAppService.NewAuthAppService(
		c.config,
		c.userDomainService,
		c.verificationService,
		c.notificationService,
		c.authService,
	)

	return nil
}

// GetDB 获取数据库连接
// 返回 GORM 数据库实例，供高级用户直接操作数据库
func (c *Client) GetDB() *gorm.DB {
	return c.db
}

// GetRedis 获取 Redis 连接
// 返回 Redis 客户端实例，供高级用户直接操作 Redis
func (c *Client) GetRedis() *redis.Client {
	return c.redis
}

// GetConfig 获取配置
// 返回配置对象，供高级用户访问配置信息
func (c *Client) GetConfig() *config.Config {
	return c.config
}

// GetUserDomainService 获取用户领域服务
// 返回用户领域服务实例，供高级用户直接调用领域服务
func (c *Client) GetUserDomainService() *userDomain.UserDomainService {
	return c.userDomainService
}

// GetVerificationService 获取验证码服务
// 返回验证码领域服务实例，供高级用户直接调用验证码服务
func (c *Client) GetVerificationService() *verificationDomain.VerificationService {
	return c.verificationService
}

// GetNotificationService 获取通知服务
// 返回通知领域服务实例，供高级用户直接调用通知服务
func (c *Client) GetNotificationService() *notificationDomain.NotificationService {
	return c.notificationService
}

// GetAuthService 获取认证服务
// 返回认证领域服务实例，供高级用户直接调用认证服务
func (c *Client) GetAuthService() *authDomain.AuthService {
	return c.authService
}

// Ping 测试连接
// 测试数据库和 Redis 连接是否正常
// 返回:
//
//	error: 连接失败的错误信息
func (c *Client) Ping(ctx context.Context) error {
	// 测试数据库连接
	sqlDB, err := c.db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	// 测试 Redis 连接
	if err := c.redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis 连接失败: %w", err)
	}

	return nil
}

// Close 关闭所有连接
// 关闭数据库和 Redis 连接，释放资源
// 返回:
//
//	error: 关闭过程中的错误
func (c *Client) Close() error {
	// 关闭数据库连接
	if c.db != nil {
		sqlDB, err := c.db.DB()
		if err != nil {
			return fmt.Errorf("获取数据库连接失败: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("关闭数据库连接失败: %w", err)
		}
	}

	// 关闭 Redis 连接
	if c.redis != nil {
		if err := c.redis.Close(); err != nil {
			return fmt.Errorf("关闭 Redis 连接失败: %w", err)
		}
	}

	return nil
}

// =============================================================================
// 用户相关操作
// =============================================================================

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Email    string // 邮箱（与手机号二选一）
	Phone    string // 手机号（与邮箱二选一）
	Password string // 密码（6-32位）
	Code     string // 验证码
	Username string // 用户名（可选，默认生成）
}

// Register 用户注册
// 参数:
//
//	ctx: 上下文
//	req: 注册请求
//
// 返回:
//
//	*UserInfo: 用户信息
//	error: 注册失败的错误
//
// 示例:
//
//	user, err := client.Register(ctx, &sdk.RegisterRequest{
//	    Email:    "test@example.com",
//	    Password: "password123",
//	    Code:     "123456",
//	})
func (c *Client) Register(ctx context.Context, req *RegisterRequest) (*UserInfo, error) {
	// 执行 BeforeCreate 钩子
	hookCtx := hook.NewContext(ctx, hook.BeforeCreate).WithRequest(req)
	if err := c.hookExecutor.Execute(hookCtx); err != nil {
		return nil, fmt.Errorf("BeforeCreate 钩子执行失败: %w", err)
	}

	// 调用应用服务注册用户
	userVo, err := c.userAppService.Register(ctx, &request.RegisterRequest{
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
		Code:     req.Code,
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}

	// 转换为 SDK 返回类型
	user := &UserInfo{
		ID:       userVo.ID,
		Username: userVo.Username,
		Nickname: "",
		Avatar:   "",
		Email:    userVo.Email,
		Phone:    "",
		Status:   userVo.Status,
	}

	// 执行 AfterCreate 钩子
	hookCtx = hook.NewContext(ctx, hook.AfterCreate).WithRequest(req)
	hookCtx.SetMetadata("user_id", user.ID)
	hookCtx.SetMetadata("user", user)
	if err := c.hookExecutor.Execute(hookCtx); err != nil {
		if c.enableLog {
			c.logger.Warn("AfterCreate 钩子执行失败", "error", err)
		}
		// AfterCreate 钩子失败不影响注册结果
	}

	return user, nil
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	LoginType string // 登录方式: username / email / phone
	Username  string // 用户名（LoginType=username 时必填）
	Email     string // 邮箱（LoginType=email 时必填）
	Phone     string // 手机号（LoginType=phone 时必填）
	Password  string // 密码（username/email 登录时必填）
	Code      string // 验证码（phone 登录时必填）
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string    // 访问令牌
	RefreshToken string    // 刷新令牌
	ExpiresIn    int       // 过期时间（秒）
	UserInfo     *UserInfo // 用户信息
}

// Login 用户登录
// 参数:
//
//	ctx: 上下文
//	req: 登录请求
//
// 返回:
//
//	*LoginResponse: 登录响应（包含 Token 和用户信息）
//	error: 登录失败的错误
//
// 示例:
//
//	// 用户名+密码登录
//	resp, err := client.Login(ctx, &sdk.LoginRequest{
//	    LoginType: "username",
//	    Username:  "john",
//	    Password:  "password123",
//	})
//	// 邮箱+密码登录
//	resp, err := client.Login(ctx, &sdk.LoginRequest{
//	    LoginType: "email",
//	    Email:     "test@example.com",
//	    Password:  "password123",
//	})
//	// 手机号+验证码登录
//	resp, err := client.Login(ctx, &sdk.LoginRequest{
//	    LoginType: "phone",
//	    Phone:     "13800138000",
//	    Code:      "123456",
//	})
func (c *Client) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 执行 BeforeLogin 钩子
	hookCtx := hook.NewContext(ctx, hook.BeforeLogin).WithRequest(req)
	if err := c.hookExecutor.Execute(hookCtx); err != nil {
		return nil, fmt.Errorf("BeforeLogin 钩子执行失败: %w", err)
	}

	// 调用应用服务登录
	loginVo, err := c.authAppService.Login(ctx, &request.LoginRequest{
		LoginType: req.LoginType,
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  req.Password,
		Code:      req.Code,
	})
	if err != nil {
		return nil, err
	}

	// 转换为 SDK 返回类型
	resp := &LoginResponse{
		AccessToken:  loginVo.AccessToken,
		RefreshToken: loginVo.RefreshToken,
		ExpiresIn:    loginVo.ExpiresIn,
		UserInfo: &UserInfo{
			ID:       loginVo.UserInfo.ID,
			Username: loginVo.UserInfo.Username,
			Nickname: "",
			Avatar:   "",
			Email:    loginVo.UserInfo.Email,
			Phone:    "",
			Status:   loginVo.UserInfo.Status,
		},
	}

	// 执行 AfterLogin 钩子
	hookCtx = hook.NewContext(ctx, hook.AfterLogin).WithRequest(req)
	hookCtx.SetMetadata("user_id", resp.UserInfo.ID)
	hookCtx.SetMetadata("user", resp.UserInfo)
	if err := c.hookExecutor.Execute(hookCtx); err != nil {
		if c.enableLog {
			c.logger.Warn("AfterLogin 钩子执行失败", "error", err)
		}
		// AfterLogin 钩子失败不影响登录结果
	}

	return resp, nil
}

// Logout 用户登出
// 参数:
//
//	ctx: 上下文
//	accessToken: 访问令牌
//	refreshToken: 刷新令牌
//
// 返回:
//
//	error: 登出失败的错误
//
// 示例:
//
//	err := client.Logout(ctx, accessToken, refreshToken)
func (c *Client) Logout(ctx context.Context, accessToken, refreshToken string) error {
	// 执行 BeforeLogout 钩子
	hookCtx := hook.NewContext(ctx, hook.BeforeLogout)
	hookCtx.SetMetadata("access_token", accessToken)
	hookCtx.SetMetadata("refresh_token", refreshToken)
	if err := c.hookExecutor.Execute(hookCtx); err != nil {
		return fmt.Errorf("BeforeLogout 钩子执行失败: %w", err)
	}

	// 调用应用服务登出
	if err := c.authAppService.Logout(ctx, accessToken, refreshToken); err != nil {
		return err
	}

	// 执行 AfterLogout 钩子
	hookCtx = hook.NewContext(ctx, hook.AfterLogout)
	if err := c.hookExecutor.Execute(hookCtx); err != nil {
		if c.enableLog {
			c.logger.Warn("AfterLogout 钩子执行失败", "error", err)
		}
		// AfterLogout 钩子失败不影响登出结果
	}

	return nil
}

// RefreshToken 刷新令牌
// 参数:
//
//	ctx: 上下文
//	refreshToken: 刷新令牌
//
// 返回:
//
//	accessToken: 新的访问令牌
//	newRefreshToken: 新的刷新令牌
//	expiresIn: 过期时间（秒）
//	error: 刷新失败的错误
//
// 示例:
//
//	accessToken, refreshToken, expiresIn, err := client.RefreshToken(ctx, oldRefreshToken)
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, expiresIn int, err error) {
	// 调用应用服务刷新令牌
	resp, err := c.authAppService.RefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", 0, err
	}

	return resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, nil
}

// UserInfo 用户信息
type UserInfo struct {
	ID       int    // 用户ID
	Username string // 用户名
	Nickname string // 昵称
	Avatar   string // 头像
	Email    string // 邮箱
	Phone    string // 手机号
	Status   int    // 状态（0:未激活 1:正常 2:禁用 3:锁定）
}

// GetUser 获取用户信息
// 参数:
//
//	ctx: 上下文
//	userID: 用户ID
//
// 返回:
//
//	*UserInfo: 用户信息
//	error: 获取失败的错误
//
// 示例:
//
//	user, err := client.GetUser(ctx, 1)
func (c *Client) GetUser(ctx context.Context, userID int) (*UserInfo, error) {
	userVo, err := c.userAppService.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:       userVo.ID,
		Username: userVo.Username,
		Nickname: "",
		Avatar:   "",
		Email:    userVo.Email,
		Phone:    "",
		Status:   userVo.Status,
	}, nil
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username *string // 用户名（可选）
	Status   *int    // 状态（可选）
}

// UpdateUser 更新用户信息
// 参数:
//
//	ctx: 上下文
//	userID: 用户ID
//	req: 更新请求
//
// 返回:
//
//	*UserInfo: 更新后的用户信息
//	error: 更新失败的错误
//
// 示例:
//
//	username := "new_username"
//	user, err := client.UpdateUser(ctx, 1, &sdk.UpdateUserRequest{
//	    Username: &username,
//	})
func (c *Client) UpdateUser(ctx context.Context, userID int, req *UpdateUserRequest) (*UserInfo, error) {
	// 执行 BeforeUpdate 钩子
	hookCtx := hook.NewContext(ctx, hook.BeforeUpdate).WithRequest(req)
	hookCtx.SetMetadata("user_id", userID)
	if err := c.hookExecutor.Execute(hookCtx); err != nil {
		return nil, fmt.Errorf("BeforeUpdate 钩子执行失败: %w", err)
	}

	// 调用应用服务更新用户
	userVo, err := c.userAppService.UpdateUser(ctx, userID, &request.UpdateUserRequest{
		Username: req.Username,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}

	// 转换为 SDK 返回类型
	user := &UserInfo{
		ID:       userVo.ID,
		Username: userVo.Username,
		Nickname: "",
		Avatar:   "",
		Email:    userVo.Email,
		Phone:    "",
		Status:   userVo.Status,
	}

	// 执行 AfterUpdate 钩子
	hookCtx = hook.NewContext(ctx, hook.AfterUpdate).WithRequest(req)
	hookCtx.SetMetadata("user_id", userID)
	hookCtx.SetMetadata("user", user)
	if err := c.hookExecutor.Execute(hookCtx); err != nil {
		if c.enableLog {
			c.logger.Warn("AfterUpdate 钩子执行失败", "error", err)
		}
		// AfterUpdate 钩子失败不影响更新结果
	}

	return user, nil
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email       string // 邮箱（与手机号二选一）
	Phone       string // 手机号（与邮箱二选一）
	Code        string // 验证码
	NewPassword string // 新密码
}

// ResetPassword 重置密码（忘记密码）
// 参数:
//
//	ctx: 上下文
//	req: 重置密码请求
//
// 返回:
//
//	error: 重置失败的错误
//
// 示例:
//
//	err := client.ResetPassword(ctx, &sdk.ResetPasswordRequest{
//	    Email:       "test@example.com",
//	    Code:        "123456",
//	    NewPassword: "newpassword123",
//	})
func (c *Client) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	// 执行 BeforeResetPassword 钩子
	hookCtx := hook.NewContext(ctx, hook.BeforeResetPassword).WithRequest(req)
	if err := c.hookExecutor.Execute(hookCtx); err != nil {
		return fmt.Errorf("BeforeResetPassword 钩子执行失败: %w", err)
	}

	// 调用应用服务重置密码
	if err := c.userAppService.ResetPassword(ctx, &request.ResetPasswordRequest{
		Email:       req.Email,
		Phone:       req.Phone,
		Code:        req.Code,
		NewPassword: req.NewPassword,
	}); err != nil {
		return err
	}

	// 执行 AfterResetPassword 钩子
	hookCtx = hook.NewContext(ctx, hook.AfterResetPassword).WithRequest(req)
	if err := c.hookExecutor.Execute(hookCtx); err != nil {
		if c.enableLog {
			c.logger.Warn("AfterResetPassword 钩子执行失败", "error", err)
		}
		// AfterResetPassword 钩子失败不影响重置结果
	}

	return nil
}

// =============================================================================
// 验证码相关操作
// =============================================================================

// SendVerificationCodeRequest 发送验证码请求
type SendVerificationCodeRequest struct {
	Target string // 目标（邮箱或手机号）
	Scene  string // 场景（register/login/reset_password 等）
}

// SendVerificationCode 发送验证码
// 参数:
//
//	ctx: 上下文
//	req: 发送验证码请求
//
// 返回:
//
//	error: 发送失败的错误
//
// 示例:
//
//	err := client.SendVerificationCode(ctx, &sdk.SendVerificationCodeRequest{
//	    Target: "test@example.com",
//	    Scene:  "register",
//	})
func (c *Client) SendVerificationCode(ctx context.Context, req *SendVerificationCodeRequest) error {
	// 1. 解析场景
	scene := verificationEnum.VerificationScene(req.Scene)
	if !scene.IsValid() {
		return fmt.Errorf("无效的验证码场景: %s", req.Scene)
	}

	// 2. 生成验证码
	code, err := c.verificationService.Generate(ctx, req.Target, scene)
	if err != nil {
		return fmt.Errorf("生成验证码失败: %w", err)
	}

	// 3. 判断目标类型并发送验证码
	targetType := c.detectTargetType(req.Target)
	switch targetType {
	case "email":
		// 发送邮件验证码
		if err := c.notificationService.SendEmailCode(ctx, req.Target, code.Code, scene); err != nil {
			return fmt.Errorf("发送邮件验证码失败: %w", err)
		}
	case "phone":
		// 发送短信验证码
		if err := c.notificationService.SendSMSCode(ctx, req.Target, code.Code, scene); err != nil {
			return fmt.Errorf("发送短信验证码失败: %w", err)
		}
	default:
		return fmt.Errorf("无法识别目标类型: %s", req.Target)
	}

	return nil
}

// VerifyCodeRequest 验证验证码请求
type VerifyCodeRequest struct {
	Target string // 目标（邮箱或手机号）
	Scene  string // 场景（register/login/reset_password 等）
	Code   string // 验证码
}

// VerifyCode 验证验证码
// 参数:
//
//	ctx: 上下文
//	req: 验证验证码请求
//
// 返回:
//
//	bool: 验证是否通过
//	error: 验证失败的错误
//
// 示例:
//
//	valid, err := client.VerifyCode(ctx, &sdk.VerifyCodeRequest{
//	    Target: "test@example.com",
//	    Scene:  "register",
//	    Code:   "123456",
//	})
func (c *Client) VerifyCode(ctx context.Context, req *VerifyCodeRequest) (bool, error) {
	// 1. 解析场景
	scene := verificationEnum.VerificationScene(req.Scene)
	if !scene.IsValid() {
		return false, fmt.Errorf("无效的验证码场景: %s", req.Scene)
	}

	// 2. 调用验证码服务验证
	if err := c.verificationService.Verify(ctx, req.Target, scene, req.Code); err != nil {
		return false, err
	}

	return true, nil
}

// =============================================================================
// 辅助方法
// =============================================================================

// detectTargetType 检测目标类型（邮箱或手机号）
func (c *Client) detectTargetType(target string) string {
	// 简单的邮箱正则判断
	if len(target) > 0 && (target[0] >= '0' && target[0] <= '9') {
		// 以数字开头，判断为手机号
		return "phone"
	}
	// 包含 @ 符号，判断为邮箱
	for _, ch := range target {
		if ch == '@' {
			return "email"
		}
	}
	// 默认判断为邮箱
	return "email"
}
