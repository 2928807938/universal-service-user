package provider

import "time"

// UserInfo OAuth 用户信息统一格式
type UserInfo struct {
	Provider string         // 平台标识(wechat/alipay/qq)
	OpenID   string         // 平台用户唯一标识
	UnionID  string         // 联合标识(可选,如微信)
	Nickname string         // 昵称
	Avatar   string         // 头像 URL
	Gender   int            // 性别(0:未知 1:男 2:女)
	Raw      map[string]any // 原始数据(用于扩展)
}

// TokenInfo OAuth 令牌信息
type TokenInfo struct {
	AccessToken  string    // 访问令牌
	RefreshToken string    // 刷新令牌(可选)
	ExpiresAt    time.Time // 令牌过期时间
	TokenType    string    // 令牌类型(通常是 "Bearer")
	Scope        string    // 授权范围
}

// AuthState OAuth 授权状态(用于防止 CSRF 攻击)
type AuthState struct {
	State       string    // 状态码
	RedirectURI string    // 回调地址
	CreatedAt   time.Time // 创建时间
}
