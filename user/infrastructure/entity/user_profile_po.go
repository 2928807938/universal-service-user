package entity

import (
	"time"
)

// UserProfilePO 用户扩展信息持久化对象,与数据库表字段对应
type UserProfilePO struct {
	ID       int        `gorm:"primaryKey;autoIncrement"`
	UserID   int        `gorm:"type:int;uniqueIndex;not null"` // 用户ID（一对一关系）
	Bio      string     `gorm:"type:varchar(500)"`             // 个人简介/描述
	Gender   int        `gorm:"type:int;default:0"`            // 性别（0:未知 1:男 2:女）
	Birthday *time.Time `gorm:"type:date"`                     // 生日（可空）
	Location string     `gorm:"type:varchar(128)"`             // 所在地
	Company  string     `gorm:"type:varchar(128)"`             // 公司
	Position string     `gorm:"type:varchar(64)"`              // 职位

	// 审计字段 - 与数据库表字段对应
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Version   int       `gorm:"default:1" json:"version"`
}

// TableName 指定表名
func (UserProfilePO) TableName() string {
	return "user_profiles"
}

// GetID 获取实体主键
func (u *UserProfilePO) GetID() int {
	return u.ID
}
