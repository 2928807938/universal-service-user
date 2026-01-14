package service

import (
	"context"

	"github.com/2928807938/universal-service-user/oauth/domain/entity"
	oautherrors "github.com/2928807938/universal-service-user/oauth/domain/errors"
	"github.com/2928807938/universal-service-user/oauth/domain/repository"
	"github.com/2928807938/universal-service-user/oauth/provider"
)

// OAuthService OAuth 领域服务
type OAuthService struct {
	repo    repository.OAuthBindingRepository
	manager *provider.Manager
}

// NewOAuthService 创建 OAuth 领域服务
func NewOAuthService(repo repository.OAuthBindingRepository) *OAuthService {
	return &OAuthService{
		repo:    repo,
		manager: provider.GetManager(),
	}
}

// GetAuthURL 获取授权 URL
// ctx: 上下文
// providerName: 平台标识
// state: 状态码
// 返回授权 URL 和错误信息
func (s *OAuthService) GetAuthURL(ctx context.Context, providerName, state string) (string, error) {
	p, err := s.getProvider(providerName)
	if err != nil {
		return "", err
	}

	return p.GetAuthURL(ctx, state)
}

// ExchangeToken 用授权码换取 Token
// ctx: 上下文
// providerName: 平台标识
// code: 授权码
// 返回令牌信息和错误信息
func (s *OAuthService) ExchangeToken(ctx context.Context, providerName, code string) (*provider.TokenInfo, error) {
	p, err := s.getProvider(providerName)
	if err != nil {
		return nil, err
	}

	token, err := p.ExchangeToken(ctx, code)
	if err != nil {
		return nil, oautherrors.ErrExchangeTokenFailed
	}

	return token, nil
}

// GetUserInfo 获取第三方用户信息
// ctx: 上下文
// providerName: 平台标识
// token: 访问令牌
// 返回用户信息和错误信息
func (s *OAuthService) GetUserInfo(ctx context.Context, providerName string, token *provider.TokenInfo) (*provider.UserInfo, error) {
	p, err := s.getProvider(providerName)
	if err != nil {
		return nil, err
	}

	userInfo, err := p.GetUserInfo(ctx, token)
	if err != nil {
		return nil, oautherrors.ErrGetUserInfoFailed
	}

	return userInfo, nil
}

// BindAccount 绑定第三方账号
// ctx: 上下文
// userID: 用户 ID
// userInfo: 第三方用户信息
// token: 令牌信息
// 返回绑定实体和错误信息
func (s *OAuthService) BindAccount(ctx context.Context, userID int, userInfo *provider.UserInfo, token *provider.TokenInfo) (*entity.OAuthBinding, error) {
	// 检查是否已存在绑定
	exists, err := s.repo.ExistsByProviderAndOpenID(ctx, userInfo.Provider, userInfo.OpenID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, oautherrors.ErrAccountAlreadyBound
	}

	// 创建绑定
	binding := &entity.OAuthBinding{
		UserID:       userID,
		Provider:     userInfo.Provider,
		OpenID:       userInfo.OpenID,
		UnionID:      userInfo.UnionID,
		Nickname:     userInfo.Nickname,
		Avatar:       userInfo.Avatar,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.ExpiresAt,
	}

	if err := s.repo.Create(ctx, binding); err != nil {
		return nil, err
	}

	return binding, nil
}

// UnbindAccount 解绑第三方账号
// ctx: 上下文
// userID: 用户 ID
// provider: 平台标识
// 返回错误信息
func (s *OAuthService) UnbindAccount(ctx context.Context, userID int, provider string) error {
	return s.repo.DeleteByUserIDAndProvider(ctx, userID, provider)
}

// GetBindingsByUserID 获取用户的所有绑定
// ctx: 上下文
// userID: 用户 ID
// 返回绑定列表和错误信息
func (s *OAuthService) GetBindingsByUserID(ctx context.Context, userID int) ([]*entity.OAuthBinding, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// GetBindingByProviderAndOpenID 根据平台和 OpenID 获取绑定
// ctx: 上下文
// provider: 平台标识
// openID: 平台用户唯一标识
// 返回绑定实体和错误信息
func (s *OAuthService) GetBindingByProviderAndOpenID(ctx context.Context, provider, openID string) (*entity.OAuthBinding, error) {
	return s.repo.FindByProviderAndOpenID(ctx, provider, openID)
}

// GetBindingByProviderAndUnionID 根据平台和 UnionID 获取绑定
// ctx: 上下文
// provider: 平台标识
// unionID: 平台联合标识
// 返回绑定实体和错误信息
func (s *OAuthService) GetBindingByProviderAndUnionID(ctx context.Context, provider, unionID string) (*entity.OAuthBinding, error) {
	if unionID == "" {
		return nil, oautherrors.ErrBindingNotFound
	}
	return s.repo.FindByProviderAndUnionID(ctx, provider, unionID)
}

// RefreshToken 刷新访问令牌
// ctx: 上下文
// binding: OAuth 绑定实体
// 返回新的令牌信息和错误信息
func (s *OAuthService) RefreshToken(ctx context.Context, binding *entity.OAuthBinding) (*provider.TokenInfo, error) {
	p, err := s.getProvider(binding.Provider)
	if err != nil {
		return nil, err
	}

	token, err := p.RefreshToken(ctx, binding.RefreshToken)
	if err != nil {
		return nil, oautherrors.ErrRefreshTokenFailed
	}

	// 更新绑定的 Token 信息
	binding.UpdateToken(token.AccessToken, token.RefreshToken, token.ExpiresAt)
	if err := s.repo.Update(ctx, binding); err != nil {
		return nil, err
	}

	return token, nil
}

// UpdateBinding 更新绑定信息
// ctx: 上下文
// binding: OAuth 绑定实体
// userInfo: 新的用户信息
// token: 新的令牌信息
// 返回错误信息
func (s *OAuthService) UpdateBinding(ctx context.Context, binding *entity.OAuthBinding, userInfo *provider.UserInfo, token *provider.TokenInfo) error {
	// 更新用户信息
	binding.UpdateUserInfo(userInfo.Nickname, userInfo.Avatar)

	// 更新 UnionID (如果有)
	if userInfo.UnionID != "" && userInfo.UnionID != binding.UnionID {
		binding.UpdateUnionID(userInfo.UnionID)
	}

	// 更新 Token 信息
	binding.UpdateToken(token.AccessToken, token.RefreshToken, token.ExpiresAt)

	return s.repo.Update(ctx, binding)
}

// IsAccountBound 判断第三方账号是否已绑定
// ctx: 上下文
// provider: 平台标识
// openID: 平台用户唯一标识
// 返回是否已绑定和错误信息
func (s *OAuthService) IsAccountBound(ctx context.Context, provider, openID string) (bool, error) {
	return s.repo.ExistsByProviderAndOpenID(ctx, provider, openID)
}

// getProvider 获取 OAuth 提供者
func (s *OAuthService) getProvider(providerName string) (provider.Provider, error) {
	p, err := s.manager.GetProvider(providerName)
	if err != nil {
		return nil, oautherrors.ErrUnsupportedProvider
	}

	if err := p.ValidateConfig(); err != nil {
		return nil, oautherrors.ErrInvalidProvider
	}

	return p, nil
}

// GetOrCreateBinding 获取或创建绑定
// 如果绑定已存在则更新信息,否则创建新绑定
func (s *OAuthService) GetOrCreateBinding(ctx context.Context, userID int, userInfo *provider.UserInfo, token *provider.TokenInfo) (*entity.OAuthBinding, error) {
	// 先尝试查找已有绑定
	binding, err := s.repo.FindByProviderAndOpenID(ctx, userInfo.Provider, userInfo.OpenID)
	if err == nil {
		// 绑定已存在,更新信息
		if err := s.UpdateBinding(ctx, binding, userInfo, token); err != nil {
			return nil, err
		}
		return binding, nil
	}

	// 如果有 UnionID,尝试通过 UnionID 查找
	if userInfo.UnionID != "" {
		binding, err = s.repo.FindByProviderAndUnionID(ctx, userInfo.Provider, userInfo.UnionID)
		if err == nil {
			// 找到了,更新信息
			if err := s.UpdateBinding(ctx, binding, userInfo, token); err != nil {
				return nil, err
			}
			return binding, nil
		}
	}

	// 绑定不存在,创建新绑定
	return s.BindAccount(ctx, userID, userInfo, token)
}

// ValidateTokenAndRefreshIfNeeded 验证 Token 并在需要时刷新
// ctx: 上下文
// binding: OAuth 绑定实体
// 返回有效的令牌信息和错误信息
func (s *OAuthService) ValidateTokenAndRefreshIfNeeded(ctx context.Context, binding *entity.OAuthBinding) (*provider.TokenInfo, error) {
	// 检查 Token 是否过期
	if !binding.IsTokenExpired() {
		// Token 未过期,直接返回
		return &provider.TokenInfo{
			AccessToken:  binding.AccessToken,
			RefreshToken: binding.RefreshToken,
			ExpiresAt:    binding.ExpiresAt,
			TokenType:    "Bearer",
		}, nil
	}

	// Token 已过期,尝试刷新
	return s.RefreshToken(ctx, binding)
}
