package jwt

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	autherrors "github.com/2928807938/universal-service-user/auth/domain/errors"
)

// Provider JWT 提供者
type Provider struct {
	config    *Config
	blacklist Blacklist // Token 黑名单（可选，基于 Redis）
}

// NewProvider 创建 JWT 提供者
func NewProvider(config *Config, blacklist Blacklist) (*Provider, error) {
	if config == nil {
		config = DefaultConfig()
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Provider{
		config:    config,
		blacklist: blacklist,
	}, nil
}

// GenerateTokenPair 生成 Token 对（AccessToken + RefreshToken）
func (p *Provider) GenerateTokenPair(ctx context.Context, userID int, username string, deviceInfo *DeviceInfo) (*TokenPair, error) {
	now := time.Now()

	// 1. 生成 Access Token
	accessJTI := uuid.New().String()
	accessClaims := &Claims{
		UserID:    userID,
		Username:  username,
		TokenType: "access",
		JTI:       accessJTI,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(p.config.AccessTokenExpire).Unix(),
	}

	accessToken, err := p.generateToken(accessClaims)
	if err != nil {
		return nil, autherrors.WrapAuthError(autherrors.AuthTokenInvalid, "生成 Access Token 失败", err)
	}

	// 2. 生成 Refresh Token
	refreshJTI := uuid.New().String()
	refreshClaims := &Claims{
		UserID:    userID,
		Username:  username,
		TokenType: "refresh",
		JTI:       refreshJTI,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(p.config.RefreshTokenExpire).Unix(),
	}

	refreshToken, err := p.generateToken(refreshClaims)
	if err != nil {
		return nil, autherrors.WrapAuthError(autherrors.AuthTokenInvalid, "生成 Refresh Token 失败", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(p.config.AccessTokenExpire.Seconds()),
	}, nil
}

// generateToken 生成 Token
func (p *Provider) generateToken(claims *Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    claims.UserID,
		"username":   claims.Username,
		"token_type": claims.TokenType,
		"jti":        claims.JTI,
		"iat":        claims.IssuedAt,
		"exp":        claims.ExpiresAt,
		"iss":        p.config.Issuer,
	})

	signedToken, err := token.SignedString([]byte(p.config.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken 验证 Token（包含黑名单检查）
func (p *Provider) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	// 1. 解析 Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(p.config.Secret), nil
	})

	if err != nil {
		return nil, autherrors.WrapAuthError(autherrors.AuthTokenInvalid, "Token 解析失败", err)
	}

	if !token.Valid {
		return nil, autherrors.ErrTokenInvalid
	}

	// 2. 提取 Claims
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, autherrors.ErrTokenInvalid
	}

	claims := &Claims{
		UserID:    int(mapClaims["user_id"].(float64)),
		Username:  mapClaims["username"].(string),
		TokenType: mapClaims["token_type"].(string),
		JTI:       mapClaims["jti"].(string),
		IssuedAt:  int64(mapClaims["iat"].(float64)),
		ExpiresAt: int64(mapClaims["exp"].(float64)),
	}

	// 3. 检查是否过期
	if claims.IsExpired() {
		return nil, autherrors.ErrTokenExpired
	}

	// 4. 检查黑名单
	if p.blacklist != nil && p.config.EnableBlacklist {
		isBlacklisted, err := p.blacklist.IsBlacklisted(ctx, claims.JTI)
		if err != nil {
			// 黑名单检查失败不阻断验证，但记录错误
			// 可以在这里添加日志
		} else if isBlacklisted {
			return nil, autherrors.ErrTokenRevoked
		}
	}

	return claims, nil
}

// RefreshToken 刷新 Token（使用 Refresh Token 生成新的 Access Token）
func (p *Provider) RefreshToken(ctx context.Context, refreshTokenString string) (*TokenPair, error) {
	// 1. 验证 Refresh Token
	claims, err := p.ValidateToken(ctx, refreshTokenString)
	if err != nil {
		return nil, autherrors.WrapAuthError(autherrors.AuthRefreshTokenFailed, "Refresh Token 验证失败", err)
	}

	// 2. 检查 Token 类型
	if !claims.IsRefreshToken() {
		return nil, autherrors.NewAuthError(autherrors.AuthTokenInvalid, "无效的 Refresh Token")
	}

	// 3. 生成新的 Token 对
	return p.GenerateTokenPair(ctx, claims.UserID, claims.Username, nil)
}

// RevokeToken 撤销 Token（加入黑名单）
func (p *Provider) RevokeToken(ctx context.Context, tokenString string) error {
	if p.blacklist == nil || !p.config.EnableBlacklist {
		return nil // 未启用黑名单，直接返回
	}

	// 1. 解析 Token 获取 JTI
	claims, err := p.ValidateToken(ctx, tokenString)
	if err != nil {
		// Token 已过期或无效，无需加入黑名单
		return nil
	}

	// 2. 计算剩余有效期
	ttl := time.Until(time.Unix(claims.ExpiresAt, 0))
	if ttl <= 0 {
		return nil // 已过期，无需加入黑名单
	}

	// 3. 加入黑名单
	return p.blacklist.Add(ctx, claims.JTI, ttl)
}

// RevokeAllUserTokens 撤销用户的所有 Token（通过时间戳方式）
func (p *Provider) RevokeAllUserTokens(ctx context.Context, userID int) error {
	if p.blacklist == nil || !p.config.EnableBlacklist {
		return nil
	}

	// 记录用户黑名单时间戳
	return p.blacklist.AddUserBlacklist(ctx, userID, time.Now())
}

// IsUserTokenRevoked 检查用户 Token 是否被撤销
func (p *Provider) IsUserTokenRevoked(ctx context.Context, userID int, issuedAt time.Time) (bool, error) {
	if p.blacklist == nil || !p.config.EnableBlacklist {
		return false, nil
	}

	return p.blacklist.IsUserBlacklisted(ctx, userID, issuedAt)
}

// HashRefreshToken 对 Refresh Token 进行哈希（用于存储到数据库）
func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
