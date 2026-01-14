package service

import (
	"context"
	"time"

	autherrors "universal-service-user/auth/domain/errors"
	"universal-service-user/auth/domain/entity"
	"universal-service-user/auth/domain/repository"
	"universal-service-user/auth/infrastructure/jwt"
	userEntity "universal-service-user/user/domain/entity"
	userService "universal-service-user/user/domain/service"
)

// AuthService 认证领域服务
type AuthService struct {
	jwtProvider       *jwt.Provider
	sessionRepo       repository.SessionRepository
	userDomainService *userService.UserDomainService
}

// NewAuthService 创建认证服务
func NewAuthService(
	jwtProvider *jwt.Provider,
	sessionRepo repository.SessionRepository,
	userDomainService *userService.UserDomainService,
) *AuthService {
	return &AuthService{
		jwtProvider:       jwtProvider,
		sessionRepo:       sessionRepo,
		userDomainService: userDomainService,
	}
}

// LoginByUsername 通过用户名+密码登录
func (s *AuthService) LoginByUsername(ctx context.Context, username, password string, deviceInfo *jwt.DeviceInfo) (*jwt.TokenPair, error) {
	// 1. 验证用户名和密码
	user, err := s.userDomainService.LoginByUsername(ctx, username, password)
	if err != nil {
		return nil, err
	}

	// 2. 生成 Token 和创建会话
	return s.createTokenAndSession(ctx, user, deviceInfo)
}

// LoginByEmail 通过邮箱+密码登录
func (s *AuthService) LoginByEmail(ctx context.Context, email, password string, deviceInfo *jwt.DeviceInfo) (*jwt.TokenPair, error) {
	// 1. 验证邮箱和密码
	user, err := s.userDomainService.LoginByEmail(ctx, email, password)
	if err != nil {
		return nil, err
	}

	// 2. 生成 Token 和创建会话
	return s.createTokenAndSession(ctx, user, deviceInfo)
}

// LoginByPhone 通过手机号+验证码登录（验证码验证由外层处理）
func (s *AuthService) LoginByPhone(ctx context.Context, phone string, deviceInfo *jwt.DeviceInfo) (*jwt.TokenPair, error) {
	// 1. 验证手机号（验证码已在外层验证）
	user, err := s.userDomainService.LoginByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	// 2. 生成 Token 和创建会话
	return s.createTokenAndSession(ctx, user, deviceInfo)
}

// createTokenAndSession 创建 Token 和会话
func (s *AuthService) createTokenAndSession(ctx context.Context, user *userEntity.User, deviceInfo *jwt.DeviceInfo) (*jwt.TokenPair, error) {
	// 1. 生成 Token 对
	tokenPair, err := s.jwtProvider.GenerateTokenPair(ctx, user.ID, user.Username, deviceInfo)
	if err != nil {
		return nil, err
	}

	// 2. 创建会话记录
	session := &entity.Session{
		UserID:           user.ID,
		RefreshTokenHash: jwt.HashRefreshToken(tokenPair.RefreshToken),
		LastActiveAt:     time.Now(),
	}

	// 设置设备信息（如果提供）
	if deviceInfo != nil {
		session.DeviceType = deviceInfo.DeviceType
		session.DeviceName = deviceInfo.DeviceName
		session.IP = deviceInfo.IP
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		// 会话创建失败不影响登录，但需要记录日志
		// 可以在这里添加日志记录
	}

	return tokenPair, nil
}

// Logout 登出
func (s *AuthService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	// 1. 撤销 Access Token（加入黑名单）
	if err := s.jwtProvider.RevokeToken(ctx, accessToken); err != nil {
		// 撤销失败不影响登出流程
	}

	// 2. 撤销 Refresh Token（加入黑名单）
	if err := s.jwtProvider.RevokeToken(ctx, refreshToken); err != nil {
		// 撤销失败不影响登出流程
	}

	// 3. 删除会话记录
	tokenHash := jwt.HashRefreshToken(refreshToken)
	if err := s.sessionRepo.DeleteByRefreshTokenHash(ctx, tokenHash); err != nil {
		// 删除失败不影响登出流程
	}

	return nil
}

// LogoutAllDevices 登出所有设备
func (s *AuthService) LogoutAllDevices(ctx context.Context, userID int) error {
	// 1. 撤销用户的所有 Token
	if err := s.jwtProvider.RevokeAllUserTokens(ctx, userID); err != nil {
		return err
	}

	// 2. 删除用户的所有会话
	if err := s.sessionRepo.DeleteByUserID(ctx, userID); err != nil {
		return err
	}

	return nil
}

// RefreshToken 刷新 Token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*jwt.TokenPair, error) {
	// 1. 验证 Refresh Token
	claims, err := s.jwtProvider.ValidateToken(ctx, refreshToken)
	if err != nil {
		return nil, autherrors.WrapAuthError(autherrors.AuthRefreshTokenFailed, "Refresh Token 验证失败", err)
	}

	// 2. 检查会话是否存在
	tokenHash := jwt.HashRefreshToken(refreshToken)
	session, err := s.sessionRepo.GetByRefreshTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, autherrors.WrapAuthError(autherrors.AuthRefreshTokenFailed, "查询会话失败", err)
	}
	if session == nil {
		return nil, autherrors.ErrSessionNotFound
	}

	// 3. 生成新的 Token 对
	newTokenPair, err := s.jwtProvider.GenerateTokenPair(ctx, claims.UserID, claims.Username, nil)
	if err != nil {
		return nil, err
	}

	// 4. 更新会话记录
	session.RefreshTokenHash = jwt.HashRefreshToken(newTokenPair.RefreshToken)
	session.UpdateLastActiveAt()
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		// 更新失败不影响刷新流程
	}

	// 5. 撤销旧的 Refresh Token
	if err := s.jwtProvider.RevokeToken(ctx, refreshToken); err != nil {
		// 撤销失败不影响刷新流程
	}

	return newTokenPair, nil
}

// ValidateToken 验证 Token
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*jwt.Claims, error) {
	return s.jwtProvider.ValidateToken(ctx, tokenString)
}

// GetUserSessions 获取用户的所有会话
func (s *AuthService) GetUserSessions(ctx context.Context, userID int) ([]*entity.Session, error) {
	return s.sessionRepo.FindByUserID(ctx, userID)
}

// KickDevice 踢出指定设备（删除会话）
func (s *AuthService) KickDevice(ctx context.Context, sessionID int) error {
	// 1. 获取会话
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return autherrors.ErrSessionNotFound
	}

	// 2. 撤销该会话的 Token（通过时间戳方式，因为我们没有存储具体的 Token）
	// 这里可以选择撤销用户所有早于当前时间的 Token，或者只删除会话记录

	// 3. 删除会话
	return s.sessionRepo.Delete(ctx, sessionID)
}

// CleanExpiredSessions 清理过期会话
func (s *AuthService) CleanExpiredSessions(ctx context.Context, expiredDuration time.Duration) error {
	expiredBefore := time.Now().Add(-expiredDuration).Unix()
	return s.sessionRepo.DeleteExpiredSessions(ctx, expiredBefore)
}
