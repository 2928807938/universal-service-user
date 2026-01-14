package entity

import (
	"time"

	basegorm "universal-service-user/share/repository/gorm"
)

// Session 会话实体 - 用于管理用户登录会话
type Session struct {
	ID               int    // 主键
	UserID           int    // 用户ID
	RefreshTokenHash string // Refresh Token 哈希值
	DeviceType       string // 设备类型: ios/android/web/desktop
	DeviceName       string // 设备名称
	IP               string // 登录IP
	LastActiveAt     time.Time // 最后活跃时间
	basegorm.AuditFields
}

// UpdateLastActiveAt 更新最后活跃时间
func (s *Session) UpdateLastActiveAt() {
	s.LastActiveAt = time.Now()
	s.Touch()
}

// IsExpired 判断会话是否过期（基于最后活跃时间）
func (s *Session) IsExpired(timeout time.Duration) bool {
	return time.Since(s.LastActiveAt) > timeout
}
