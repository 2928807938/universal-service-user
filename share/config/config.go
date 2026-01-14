package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server        ServerConfig        `yaml:"server"`
	Database      DatabaseConfig      `yaml:"database"`
	Redis         RedisConfig         `yaml:"redis"`
	JWT           JWTConfig           `yaml:"jwt"`
	Verification  VerificationConfig  `yaml:"verification"`
	Email         EmailConfig         `yaml:"email"`
	SMS           SMSConfig           `yaml:"sms"`
	OAuth         OAuthConfig         `yaml:"oauth"`
	LoginRateLimit LoginRateLimitConfig `yaml:"login_rate_limit"`
	Webhook       WebhookConfig       `yaml:"webhook"`
	Logging       LoggingConfig       `yaml:"logging"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string `yaml:"port"`
	Mode         string `yaml:"mode"`          // debug / release
	ReadTimeout  int    `yaml:"read_timeout"`  // 秒
	WriteTimeout int    `yaml:"write_timeout"` // 秒
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `yaml:"driver"`            // postgres / mysql / sqlite
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"dbname"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"` // 秒
	AutoMigrate     bool   `yaml:"auto_migrate"`      // 是否自动建表
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret        string `yaml:"secret"`
	Expire        int    `yaml:"expire"`         // Access Token 过期时间（秒）
	RefreshExpire int    `yaml:"refresh_expire"` // Refresh Token 过期时间（秒）
	Issuer        string `yaml:"issuer"`
}

// VerificationConfig 验证码配置
type VerificationConfig struct {
	CodeLength int `yaml:"code_length"` // 验证码长度
	Expire     int `yaml:"expire"`      // 过期时间（秒）
	RateLimit  int `yaml:"rate_limit"`  // 发送间隔（秒）
}

// EmailConfig 邮箱配置
type EmailConfig struct {
	Enabled   bool              `yaml:"enabled"`
	SMTP      SMTPConfig        `yaml:"smtp"`
	Templates map[string]string `yaml:"templates"`
	RateLimit RateLimitConfig   `yaml:"rate_limit"`
}

// SMTPConfig SMTP 配置
type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

// SMSConfig 短信配置
type SMSConfig struct {
	Enabled   bool              `yaml:"enabled"`
	Provider  string            `yaml:"provider"` // tencent / aliyun / custom
	Tencent   TencentSMSConfig  `yaml:"tencent"`
	Aliyun    AliyunSMSConfig   `yaml:"aliyun"`
	RateLimit RateLimitConfig   `yaml:"rate_limit"`
}

// TencentSMSConfig 腾讯云短信配置
type TencentSMSConfig struct {
	SecretID  string            `yaml:"secret_id"`
	SecretKey string            `yaml:"secret_key"`
	AppID     string            `yaml:"app_id"`
	Sign      string            `yaml:"sign"`
	Templates map[string]string `yaml:"templates"`
}

// AliyunSMSConfig 阿里云短信配置
type AliyunSMSConfig struct {
	AccessKeyID     string            `yaml:"access_key_id"`
	AccessKeySecret string            `yaml:"access_key_secret"`
	SignName        string            `yaml:"sign_name"`
	Templates       map[string]string `yaml:"templates"`
}

// RateLimitConfig 限频配置
type RateLimitConfig struct {
	Interval int `yaml:"interval"` // 秒
}

// OAuthConfig OAuth 配置
type OAuthConfig struct {
	Wechat OAuthProviderConfig `yaml:"wechat"`
	Alipay OAuthProviderConfig `yaml:"alipay"`
	QQ     OAuthProviderConfig `yaml:"qq"`
}

// OAuthProviderConfig OAuth 平台配置
type OAuthProviderConfig struct {
	Enabled         bool   `yaml:"enabled"`
	AppID           string `yaml:"app_id"`
	AppSecret       string `yaml:"app_secret"`
	ClientID        string `yaml:"client_id"` // 兼容不同平台的命名
	ClientSecret    string `yaml:"client_secret"`
	PrivateKey      string `yaml:"private_key"`       // 支付宝私钥
	AlipayPublicKey string `yaml:"alipay_public_key"` // 支付宝公钥
	RedirectURI     string `yaml:"redirect_uri"`
}

// LoginRateLimitConfig 登录限频配置
type LoginRateLimitConfig struct {
	Enabled bool                   `yaml:"enabled"`
	IP      LoginLimitRule         `yaml:"ip"`
	Account LoginLimitRule         `yaml:"account"`
	Device  LoginLimitRule         `yaml:"device"`
}

// LoginLimitRule 登录限制规则
type LoginLimitRule struct {
	MaxAttempts   int `yaml:"max_attempts"`   // 最大尝试次数
	MaxFailures   int `yaml:"max_failures"`   // 最大失败次数
	Window        int `yaml:"window"`         // 时间窗口（秒）
	BlockDuration int `yaml:"block_duration"` // 封禁时长（秒）
	LockDuration  int `yaml:"lock_duration"`  // 锁定时长（秒）
}

// WebhookConfig Webhook 配置
type WebhookConfig struct {
	Enabled bool     `yaml:"enabled"`
	URL     string   `yaml:"url"`
	Secret  string   `yaml:"secret"`
	Events  []string `yaml:"events"`
	Timeout int      `yaml:"timeout"` // 秒
	Retry   int      `yaml:"retry"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level             string `yaml:"level"`               // debug / info / warn / error
	Format            string `yaml:"format"`              // json / text
	EnableSensibleMask bool  `yaml:"enable_sensible_mask"` // 敏感信息脱敏
}

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 如果未指定配置文件路径,尝试从默认位置加载
	if configPath == "" {
		configPath = findConfigFile()
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 从环境变量覆盖敏感配置
	overrideFromEnv(&config)

	// 验证必需的配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &config, nil
}

// findConfigFile 查找配置文件
func findConfigFile() string {
	possiblePaths := []string{
		"config.yaml",
		"config.yml",
		"configs/config.yaml",
		"configs/config.yml",
		"./config.yaml",
		"./config.yml",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// 返回默认路径
	return "config.yaml"
}

// overrideFromEnv 从环境变量覆盖配置
func overrideFromEnv(config *Config) {
	// 数据库配置
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &config.Database.Port)
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		config.Database.DBName = dbname
	}

	// Redis 配置
	if host := os.Getenv("REDIS_HOST"); host != "" {
		config.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &config.Redis.Port)
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		config.Redis.Password = password
	}

	// JWT Secret
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		config.JWT.Secret = secret
	}

	// 服务端口
	if port := os.Getenv("PORT"); port != "" {
		config.Server.Port = port
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 验证数据库配置
	if c.Database.Driver == "" {
		return fmt.Errorf("数据库驱动不能为空")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("数据库主机不能为空")
	}

	// 验证 Redis 配置
	if c.Redis.Host == "" {
		return fmt.Errorf("Redis 主机不能为空")
	}

	// 验证 JWT 配置
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT Secret 不能为空")
	}

	return nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			c.Host, c.User, c.Password, c.DBName, c.Port)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.User, c.Password, c.Host, c.Port, c.DBName)
	case "sqlite":
		return filepath.Join(".", c.DBName+".db")
	default:
		return ""
	}
}

// GetRedisAddr 获取 Redis 地址
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetAccessTokenDuration 获取 Access Token 过期时长
func (c *JWTConfig) GetAccessTokenDuration() time.Duration {
	return time.Duration(c.Expire) * time.Second
}

// GetRefreshTokenDuration 获取 Refresh Token 过期时长
func (c *JWTConfig) GetRefreshTokenDuration() time.Duration {
	return time.Duration(c.RefreshExpire) * time.Second
}
