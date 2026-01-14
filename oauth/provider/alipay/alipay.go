package alipay

import (
	"context"
	"fmt"

	"universal-service-user/oauth/provider"
)

// Provider 支付宝 OAuth 提供者
type Provider struct {
	config *provider.AlipayConfig
}

// New 创建支付宝 OAuth 提供者
func New(config *provider.AlipayConfig) *Provider {
	return &Provider{config: config}
}

// GetName 获取平台唯一标识
func (p *Provider) GetName() string { return "alipay" }

// GetDisplayName 获取平台显示名称
func (p *Provider) GetDisplayName() string { return "支付宝" }

// GetAuthURL 生成授权跳转 URL
func (p *Provider) GetAuthURL(ctx context.Context, state string) (string, error) {
	return fmt.Sprintf("https://openauth.alipay.com/oauth2/publicAppAuthorize.htm?app_id=%s&redirect_uri=%s&scope=auth_user&state=%s",
		p.config.AppID, p.config.RedirectURI, state), nil
}

// ExchangeToken 用授权码换取 Access Token
func (p *Provider) ExchangeToken(ctx context.Context, code string) (*provider.TokenInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetUserInfo 获取第三方用户信息
func (p *Provider) GetUserInfo(ctx context.Context, token *provider.TokenInfo) (*provider.UserInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

// RefreshToken 刷新访问令牌
func (p *Provider) RefreshToken(ctx context.Context, refreshToken string) (*provider.TokenInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

// ValidateConfig 验证配置是否完整有效
func (p *Provider) ValidateConfig() error {
	if p.config == nil || p.config.AppID == "" || p.config.PrivateKey == "" {
		return fmt.Errorf("支付宝配置不完整")
	}
	return nil
}
