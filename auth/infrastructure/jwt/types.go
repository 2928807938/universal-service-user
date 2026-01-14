package jwt

import "time"

// TokenPair Token 对
type TokenPair struct {
	AccessToken  string `json:"access_token"`  // 访问令牌
	RefreshToken string `json:"refresh_token"` // 刷新令牌
	ExpiresIn    int64  `json:"expires_in"`    // 访问令牌过期时间（秒）
}

// Claims JWT Claims
type Claims struct {
	UserID    int    `json:"user_id"`    // 用户ID
	Username  string `json:"username"`   // 用户名
	TokenType string `json:"token_type"` // token 类型: access/refresh
	JTI       string `json:"jti"`        // JWT ID，用于黑名单
	IssuedAt  int64  `json:"iat"`        // 签发时间
	ExpiresAt int64  `json:"exp"`        // 过期时间
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	DeviceID   string `json:"device_id"`   // 设备唯一标识
	DeviceType string `json:"device_type"` // 设备类型: ios/android/web/desktop
	DeviceName string `json:"device_name"` // 设备名称
	IP         string `json:"ip"`          // 登录IP
	Browser    string `json:"browser"`     // 浏览器
	OS         string `json:"os"`          // 操作系统
}

// IsExpired 判断 Claims 是否过期
func (c *Claims) IsExpired() bool {
	return time.Now().Unix() > c.ExpiresAt
}

// IsAccessToken 判断是否为访问令牌
func (c *Claims) IsAccessToken() bool {
	return c.TokenType == "access"
}

// IsRefreshToken 判断是否为刷新令牌
func (c *Claims) IsRefreshToken() bool {
	return c.TokenType == "refresh"
}
