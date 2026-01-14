package sdk

import (
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/2928807938/universal-service-user/notification/infrastructure/email"
	"github.com/2928807938/universal-service-user/share/config"
)

// Option SDK 配置选项函数
type Option func(*Client)

// WithDatabase 设置数据库连接
// db: GORM 数据库实例
// 示例:
//
//	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//	sdk.New(sdk.WithDatabase(db))
func WithDatabase(db *gorm.DB) Option {
	return func(c *Client) {
		c.db = db
	}
}

// WithRedis 设置 Redis 连接
// rdb: Redis 客户端实例
// 示例:
//
//	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
//	sdk.New(sdk.WithRedis(rdb))
func WithRedis(rdb *redis.Client) Option {
	return func(c *Client) {
		c.redis = rdb
	}
}

// WithConfig 使用完整的配置对象
// cfg: 配置对象
// 示例:
//
//	cfg, _ := config.LoadConfig("config.yaml")
//	sdk.New(sdk.WithConfig(cfg))
func WithConfig(cfg *config.Config) Option {
	return func(c *Client) {
		c.config = cfg
	}
}

// WithAutoMigrate 设置是否自动建表
// autoMigrate: true 表示自动建表，false 表示跳过
// 默认值: true
// 示例:
//
//	sdk.New(sdk.WithAutoMigrate(true))
func WithAutoMigrate(autoMigrate bool) Option {
	return func(c *Client) {
		c.autoMigrate = autoMigrate
	}
}

// WithJWT 设置 JWT 配置
// secret: JWT 密钥
// accessExpire: Access Token 过期时间
// refreshExpire: Refresh Token 过期时间
// issuer: 签发者
// 示例:
//
//	sdk.New(sdk.WithJWT("your-secret", 2*time.Hour, 7*24*time.Hour, "user-service"))
func WithJWT(secret string, accessExpire, refreshExpire time.Duration, issuer string) Option {
	return func(c *Client) {
		if c.config == nil {
			c.config = &config.Config{}
		}
		c.config.JWT = config.JWTConfig{
			Secret:        secret,
			Expire:        int(accessExpire.Seconds()),
			RefreshExpire: int(refreshExpire.Seconds()),
			Issuer:        issuer,
		}
	}
}

// WithVerification 设置验证码配置
// codeLength: 验证码长度
// expire: 验证码过期时间
// rateLimit: 发送间隔
// 示例:
//
//	sdk.New(sdk.WithVerification(6, 5*time.Minute, 60*time.Second))
func WithVerification(codeLength int, expire, rateLimit time.Duration) Option {
	return func(c *Client) {
		if c.config == nil {
			c.config = &config.Config{}
		}
		c.config.Verification = config.VerificationConfig{
			CodeLength: codeLength,
			Expire:     int(expire.Seconds()),
			RateLimit:  int(rateLimit.Seconds()),
		}
	}
}

// WithEmailProvider 设置邮箱服务提供者
// cfg: SMTP 配置
// 示例:
//
//	smtpCfg := &email.SMTPConfig{
//	    Host:     "smtp.example.com",
//	    Port:     587,
//	    Username: "noreply@example.com",
//	    Password: "password",
//	    From:     "系统通知 <noreply@example.com>",
//	    Templates: map[string]string{
//	        "register": "验证码：{code}",
//	    },
//	}
//	sdk.New(sdk.WithEmailProvider(smtpCfg))
func WithEmailProvider(cfg *email.SMTPConfig) Option {
	return func(c *Client) {
		c.emailConfig = cfg
	}
}

// WithSMSProvider 设置短信服务提供者
// provider: 短信服务商名称 (tencent/aliyun)
// cfg: 短信服务配置
// 注意: 使用此选项前需要确保已经实现对应的短信服务提供者
// 示例:
//
//	sdk.New(sdk.WithSMSProvider("tencent", smsConfig))
func WithSMSProvider(provider string, cfg interface{}) Option {
	return func(c *Client) {
		c.smsProvider = provider
		c.smsConfig = cfg
	}
}

// WithOAuthProvider 设置 OAuth 提供者配置
// provider: OAuth 平台名称 (wechat/alipay/qq)
// cfg: OAuth 配置
// 注意: 使用此选项前需要确保已经实现对应的 OAuth 提供者
// 示例:
//
//	oauthCfg := &config.OAuthProviderConfig{
//	    Enabled:     true,
//	    AppID:       "your-app-id",
//	    AppSecret:   "your-app-secret",
//	    RedirectURI: "https://your-domain.com/callback",
//	}
//	sdk.New(sdk.WithOAuthProvider("wechat", oauthCfg))
func WithOAuthProvider(provider string, cfg *config.OAuthProviderConfig) Option {
	return func(c *Client) {
		if c.config == nil {
			c.config = &config.Config{}
		}
		switch provider {
		case "wechat":
			c.config.OAuth.Wechat = *cfg
		case "alipay":
			c.config.OAuth.Alipay = *cfg
		case "qq":
			c.config.OAuth.QQ = *cfg
		}
	}
}

// WithLoginRateLimit 设置登录限频配置
// enabled: 是否启用限频
// ipMaxAttempts: 同一 IP 最大尝试次数
// ipWindow: IP 限制时间窗口
// ipBlockDuration: IP 封禁时长
// accountMaxFailures: 同一账号最大失败次数
// accountWindow: 账号限制时间窗口
// accountLockDuration: 账号锁定时长
// 示例:
//
//	sdk.New(sdk.WithLoginRateLimit(
//	    true,
//	    5, 60*time.Second, 10*time.Minute,
//	    5, 10*time.Minute, 30*time.Minute,
//	))
func WithLoginRateLimit(
	enabled bool,
	ipMaxAttempts int, ipWindow, ipBlockDuration time.Duration,
	accountMaxFailures int, accountWindow, accountLockDuration time.Duration,
) Option {
	return func(c *Client) {
		if c.config == nil {
			c.config = &config.Config{}
		}
		c.config.LoginRateLimit = config.LoginRateLimitConfig{
			Enabled: enabled,
			IP: config.LoginLimitRule{
				MaxAttempts:   ipMaxAttempts,
				Window:        int(ipWindow.Seconds()),
				BlockDuration: int(ipBlockDuration.Seconds()),
			},
			Account: config.LoginLimitRule{
				MaxFailures:  accountMaxFailures,
				Window:       int(accountWindow.Seconds()),
				LockDuration: int(accountLockDuration.Seconds()),
			},
		}
	}
}

// WithLogger 设置日志记录器
// logger: 自定义日志实现
// 示例:
//
//	sdk.New(sdk.WithLogger(myLogger))
func WithLogger(logger Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithEnableLog 设置是否启用日志
// enable: true 表示启用日志，false 表示禁用
// 默认值: false
// 示例:
//
//	sdk.New(sdk.WithEnableLog(true))
func WithEnableLog(enable bool) Option {
	return func(c *Client) {
		c.enableLog = enable
	}
}

// Logger 日志接口
// 使用者可以实现此接口来自定义日志行为
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}
