package entity

import (
	"time"
)

// UserPO 用户持久化对象,与数据库表字段对应
type UserPO struct {
	ID           int    `gorm:"primaryKey;autoIncrement"`
	Username     string `gorm:"type:varchar(64);uniqueIndex"`           // 用户名（唯一，可空）
	Nickname     string `gorm:"type:varchar(64)"`                       // 昵称
	Avatar       string `gorm:"type:varchar(512)"`                      // 头像URL
	Email        string `gorm:"type:varchar(255);uniqueIndex"`          // 邮箱（唯一，可空）
	Phone        string `gorm:"type:varchar(20);uniqueIndex"`           // 手机号（唯一，可空）
	PasswordHash string `gorm:"type:varchar(255)"`                      // 密码哈希（可空，第三方登录用户可能无密码）
	Status       int    `gorm:"type:int;default:0;index"`               // 用户状态

	// 审计字段 - 与数据库表字段对应
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt time.Time `gorm:"index" json:"deleted_at,omitempty"`
	Version   int       `gorm:"default:1" json:"version"`
}

// TableName 指定表名
func (UserPO) TableName() string {
	return "users"
}

// GetID 获取实体主键
func (u *UserPO) GetID() int {
	return u.ID
}
