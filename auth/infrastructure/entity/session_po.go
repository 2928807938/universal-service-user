package entity

import (
	"time"
)

// SessionPO 会话持久化对象,与数据库表字段对应
type SessionPO struct {
	ID               int       `gorm:"primaryKey;autoIncrement"`
	UserID           int       `gorm:"type:int;index;not null"`                                   // 用户ID
	RefreshTokenHash string    `gorm:"type:varchar(128);uniqueIndex;not null"`                    // Refresh Token 哈希值
	DeviceType       string    `gorm:"type:varchar(32)"`                                          // 设备类型
	DeviceName       string    `gorm:"type:varchar(128)"`                                         // 设备名称
	IP               string    `gorm:"type:varchar(64)"`                                          // 登录IP
	LastActiveAt     time.Time `gorm:"type:timestamp;index:idx_user_active,composite:1;not null"` // 最后活跃时间

	// 审计字段 - 与数据库表字段对应
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Version   int       `gorm:"default:1" json:"version"`
}

// TableName 指定表名
func (SessionPO) TableName() string {
	return "user_sessions"
}

// GetID 获取实体主键
func (s *SessionPO) GetID() int {
	return s.ID
}
