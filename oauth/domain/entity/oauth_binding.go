package entity

import (
	basegorm "github.com/2928807938/universal-service-user/share/repository/gorm"
	"time"
)

// OAuthBinding OAuth 绑定实体 - 聚合根
// 表示用户与第三方平台账号的绑定关系
type OAuthBinding struct {
	ID                   int       // 主键
	UserID               int       // 用户 ID (外键)
	Provider             string    // 平台标识 (wechat/alipay/qq)
	OpenID               string    // 平台用户唯一标识
	UnionID              string    // 平台联合标识 (可空, 如微信)
	Nickname             string    // 平台昵称
	Avatar               string    // 平台头像
	AccessToken          string    // 访问令牌 (加密存储)
	RefreshToken         string    // 刷新令牌 (加密存储)
	ExpiresAt            time.Time // 令牌过期时间
	basegorm.AuditFields           // 审计字段
}

// IsTokenExpired 判断 Token 是否过期
func (b *OAuthBinding) IsTokenExpired() bool {
	return time.Now().After(b.ExpiresAt)
}

// UpdateToken 更新 Token 信息
func (b *OAuthBinding) UpdateToken(accessToken, refreshToken string, expiresAt time.Time) {
	b.AccessToken = accessToken
	b.RefreshToken = refreshToken
	b.ExpiresAt = expiresAt
	b.Touch()
}

// UpdateUserInfo 更新用户信息
func (b *OAuthBinding) UpdateUserInfo(nickname, avatar string) {
	b.Nickname = nickname
	b.Avatar = avatar
	b.Touch()
}

// UpdateUnionID 更新 UnionID
func (b *OAuthBinding) UpdateUnionID(unionID string) {
	b.UnionID = unionID
	b.Touch()
}
