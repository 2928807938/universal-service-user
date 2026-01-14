package jwt

import "time"

// Config JWT 配置
type Config struct {
	Secret             string        // JWT 密钥
	AccessTokenExpire  time.Duration // Access Token 过期时间（默认 2 小时）
	RefreshTokenExpire time.Duration // Refresh Token 过期时间（默认 7 天）
	Issuer             string        // 签发者
	EnableBlacklist    bool          // 是否启用黑名单（需要 Redis）
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Secret:             "default-jwt-secret-please-change-in-production",
		AccessTokenExpire:  2 * time.Hour,
		RefreshTokenExpire: 7 * 24 * time.Hour,
		Issuer:             "universal-user-service",
		EnableBlacklist:    true,
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Secret == "" {
		return ErrInvalidConfig("JWT secret cannot be empty")
	}
	if c.AccessTokenExpire <= 0 {
		return ErrInvalidConfig("Access token expire must be greater than 0")
	}
	if c.RefreshTokenExpire <= 0 {
		return ErrInvalidConfig("Refresh token expire must be greater than 0")
	}
	if c.Issuer == "" {
		c.Issuer = "universal-user-service"
	}
	return nil
}

// ErrInvalidConfig 配置错误
func ErrInvalidConfig(message string) error {
	return &ConfigError{Message: message}
}

// ConfigError 配置错误
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return "JWT配置错误: " + e.Message
}
