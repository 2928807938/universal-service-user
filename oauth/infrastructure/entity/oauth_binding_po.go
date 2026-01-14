package entity

import (
	"time"
)

// OAuthBindingPO OAuth 绑定持久化对象,与数据库表字段对应
type OAuthBindingPO struct {
	ID           int    `gorm:"primaryKey;autoIncrement"`
	UserID       int    `gorm:"type:int;index:idx_user_provider,priority:1;not null"` // 用户 ID
	Provider     string `gorm:"type:varchar(32);index:idx_user_provider,priority:2;not null"` // 平台标识
	OpenID       string `gorm:"type:varchar(128);index;not null"` // 平台用户唯一标识
	UnionID      string `gorm:"type:varchar(128);index"` // 平台联合标识(可空)
	Nickname     string `gorm:"type:varchar(64)"` // 平台昵称
	Avatar       string `gorm:"type:varchar(512)"` // 平台头像
	AccessToken  string `gorm:"type:varchar(512)"` // 访问令牌(加密存储)
	RefreshToken string `gorm:"type:varchar(512)"` // 刷新令牌(加密存储)
	ExpiresAt    time.Time `gorm:"type:timestamp"` // 令牌过期时间

	// 审计字段 - 与数据库表字段对应
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt time.Time `gorm:"index" json:"deleted_at,omitempty"`
	Version   int       `gorm:"default:1" json:"version"`
}

// TableName 指定表名
func (OAuthBindingPO) TableName() string {
	return "user_oauth_bindings"
}

// GetID 获取实体主键
func (o *OAuthBindingPO) GetID() int {
	return o.ID
}
