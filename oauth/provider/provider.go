package provider

import "context"

// Provider OAuth 提供者接口
// 定义统一的 OAuth 2.0 授权流程接口
type Provider interface {
	// GetName 获取平台唯一标识
	// 返回值示例: "wechat", "alipay", "qq"
	GetName() string

	// GetDisplayName 获取平台显示名称
	// 返回值示例: "微信", "支付宝", "QQ"
	GetDisplayName() string

	// GetAuthURL 生成授权跳转 URL
	// ctx: 上下文
	// state: 状态码(用于防止 CSRF 攻击)
	// 返回授权 URL 和错误信息
	GetAuthURL(ctx context.Context, state string) (string, error)

	// ExchangeToken 用授权码换取 Access Token
	// ctx: 上下文
	// code: 授权码
	// 返回令牌信息和错误信息
	ExchangeToken(ctx context.Context, code string) (*TokenInfo, error)

	// GetUserInfo 获取第三方用户信息
	// ctx: 上下文
	// token: 访问令牌
	// 返回用户信息和错误信息
	GetUserInfo(ctx context.Context, token *TokenInfo) (*UserInfo, error)

	// RefreshToken 刷新访问令牌(可选)
	// ctx: 上下文
	// refreshToken: 刷新令牌
	// 返回新的令牌信息和错误信息
	RefreshToken(ctx context.Context, refreshToken string) (*TokenInfo, error)

	// ValidateConfig 验证配置是否完整有效
	// 返回 error 表示配置无效
	ValidateConfig() error
}
