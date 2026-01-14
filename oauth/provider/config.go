package provider

// Config OAuth 通用配置结构
type Config struct {
	Enabled     bool     // 是否启用
	AppID       string   // 应用 ID (或 Client ID)
	AppSecret   string   // 应用密钥 (或 Client Secret)
	RedirectURI string   // 回调地址
	Scopes      []string // 授权范围
}

// WechatConfig 微信 OAuth 配置
type WechatConfig struct {
	Config
}

// AlipayConfig 支付宝 OAuth 配置
type AlipayConfig struct {
	Config
	PrivateKey      string // 应用私钥
	AlipayPublicKey string // 支付宝公钥
}

// QQConfig QQ OAuth 配置
type QQConfig struct {
	Config
	AppKey string // App Key (QQ 使用 AppKey 而非 AppSecret)
}
