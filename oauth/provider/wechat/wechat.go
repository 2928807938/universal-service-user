package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"universal-service-user/oauth/provider"
)

// Provider 微信 OAuth 提供者
type Provider struct {
	config      *provider.WechatConfig
	authURL     string
	tokenURL    string
	userInfoURL string
	httpClient  *http.Client
}

// New 创建微信 OAuth 提供者
func New(config *provider.WechatConfig) *Provider {
	return &Provider{
		config:      config,
		authURL:     "https://open.weixin.qq.com/connect/qrconnect",
		tokenURL:    "https://api.weixin.qq.com/sns/oauth2/access_token",
		userInfoURL: "https://api.weixin.qq.com/sns/userinfo",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetName 获取平台唯一标识
func (p *Provider) GetName() string {
	return "wechat"
}

// GetDisplayName 获取平台显示名称
func (p *Provider) GetDisplayName() string {
	return "微信"
}

// GetAuthURL 生成授权跳转 URL
func (p *Provider) GetAuthURL(ctx context.Context, state string) (string, error) {
	params := url.Values{}
	params.Add("appid", p.config.AppID)
	params.Add("redirect_uri", p.config.RedirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "snsapi_login")
	params.Add("state", state)

	return fmt.Sprintf("%s?%s#wechat_redirect", p.authURL, params.Encode()), nil
}

// ExchangeToken 用授权码换取 Access Token
func (p *Provider) ExchangeToken(ctx context.Context, code string) (*provider.TokenInfo, error) {
	params := url.Values{}
	params.Add("appid", p.config.AppID)
	params.Add("secret", p.config.AppSecret)
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")

	reqURL := fmt.Sprintf("%s?%s", p.tokenURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		OpenID       string `json:"openid"`
		Scope        string `json:"scope"`
		UnionID      string `json:"unionid"`
		ErrCode      int    `json:"errcode"`
		ErrMsg       string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("微信 API 错误 [%d]: %s", result.ErrCode, result.ErrMsg)
	}

	return &provider.TokenInfo{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(result.ExpiresIn) * time.Second),
		TokenType:    "Bearer",
		Scope:        result.Scope,
	}, nil
}

// GetUserInfo 获取第三方用户信息
func (p *Provider) GetUserInfo(ctx context.Context, token *provider.TokenInfo) (*provider.UserInfo, error) {
	// 实现省略,参考之前的代码
	return nil, fmt.Errorf("not implemented")
}

// RefreshToken 刷新访问令牌
func (p *Provider) RefreshToken(ctx context.Context, refreshToken string) (*provider.TokenInfo, error) {
	// 实现省略,参考之前的代码
	return nil, fmt.Errorf("not implemented")
}

// ValidateConfig 验证配置是否完整有效
func (p *Provider) ValidateConfig() error {
	if p.config == nil {
		return fmt.Errorf("微信配置不能为空")
	}

	if p.config.AppID == "" {
		return fmt.Errorf("微信 AppID 不能为空")
	}

	if p.config.AppSecret == "" {
		return fmt.Errorf("微信 AppSecret 不能为空")
	}

	if p.config.RedirectURI == "" {
		return fmt.Errorf("微信回调地址不能为空")
	}

	return nil
}
