package service

import (
	"context"
	"fmt"

	"universal-service-user/api/user-api/converter"
	"universal-service-user/api/user-api/dto/request"
	"universal-service-user/api/user-api/dto/vo"
	authDomain "universal-service-user/auth/domain/service"
	"universal-service-user/auth/infrastructure/jwt"
	notificationDomain "universal-service-user/notification/domain/service"
	"universal-service-user/share/config"
	userErrors "universal-service-user/user/domain/errors"
	userDomain "universal-service-user/user/domain/service"
	verificationEnum "universal-service-user/verification/domain/enum"
	verificationDomain "universal-service-user/verification/domain/service"
)

// AuthAppService 认证应用服务
type AuthAppService struct {
	config              *config.Config
	userDomainService   *userDomain.UserDomainService
	verificationService *verificationDomain.VerificationService
	notificationService *notificationDomain.NotificationService
	authService         *authDomain.AuthService
	converter           *converter.UserConverter
}

// NewAuthAppService 创建认证应用服务
func NewAuthAppService(
	config *config.Config,
	userDomainService *userDomain.UserDomainService,
	verificationService *verificationDomain.VerificationService,
	notificationService *notificationDomain.NotificationService,
	authService *authDomain.AuthService,
) *AuthAppService {
	return &AuthAppService{
		config:              config,
		userDomainService:   userDomainService,
		verificationService: verificationService,
		notificationService: notificationService,
		authService:         authService,
		converter:           converter.NewUserConverter(),
	}
}

// Login 用户登录
func (s *AuthAppService) Login(ctx context.Context, req *request.LoginRequest) (*vo.LoginResponse, error) {
	var tokenPair *jwt.TokenPair
	var err error

	// 构建设备信息（可选）
	var deviceInfo *jwt.DeviceInfo
	// TODO: 从请求中提取设备信息

	// 根据登录类型调用对应的认证领域服务方法
	switch req.LoginType {
	case "username":
		// 用户名+密码登录（参数验证在领域服务中）
		tokenPair, err = s.authService.LoginByUsername(ctx, req.Username, req.Password, deviceInfo)

	case "email":
		// 邮箱+密码登录（参数验证在领域服务中）
		tokenPair, err = s.authService.LoginByEmail(ctx, req.Email, req.Password, deviceInfo)

	case "phone":
		// 手机号+验证码登录
		// 验证验证码
		if err := s.verificationService.VerifyAndDelete(ctx, req.Phone, verificationEnum.SceneLogin, req.Code); err != nil {
			return nil, err
		}
		tokenPair, err = s.authService.LoginByPhone(ctx, req.Phone, deviceInfo)

	default:
		return nil, userErrors.NewValidationError("不支持的登录方式: " + req.LoginType)
	}

	if err != nil {
		return nil, err
	}

	// 从 token 中解析用户信息
	claims, err := s.authService.ValidateToken(ctx, tokenPair.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("解析令牌失败: %w", err)
	}

	// 获取用户详细信息
	user, err := s.userDomainService.GetUser(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// 组装响应
	return &vo.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    s.config.JWT.Expire,
		UserInfo:     s.converter.ToVo(user),
	}, nil
}

// Logout 用户登出
func (s *AuthAppService) Logout(ctx context.Context, accessToken string, refreshToken string) error {
	// 调用认证领域服务登出（将 token 加入黑名单并删除会话）
	return s.authService.Logout(ctx, accessToken, refreshToken)
}

// RefreshToken 刷新令牌
func (s *AuthAppService) RefreshToken(ctx context.Context, refreshToken string) (*vo.RefreshTokenResponse, error) {
	// 调用认证领域服务刷新令牌
	tokenPair, err := s.authService.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return &vo.RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    s.config.JWT.Expire,
	}, nil
}
