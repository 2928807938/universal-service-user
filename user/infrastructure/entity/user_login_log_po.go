package entity

import (
	"time"
)

// UserLoginLogPO 用户登录日志持久化对象,与数据库表字段对应
type UserLoginLogPO struct {
	ID         int        `gorm:"primaryKey;autoIncrement"`
	UserID     int        `gorm:"type:int;index;not null"`                                     // 用户ID
	LoginType  string     `gorm:"type:varchar(32);not null"`                                   // 登录方式
	LoginAt    time.Time  `gorm:"type:timestamp;index:idx_user_login_at,composite:1;not null"` // 登录时间
	LogoutAt   *time.Time `gorm:"type:timestamp"`                                              // 登出时间（可空）
	IP         string     `gorm:"type:varchar(64);index"`                                      // 登录IP
	DeviceType string     `gorm:"type:varchar(32)"`                                            // 设备类型
	DeviceName string     `gorm:"type:varchar(128)"`                                           // 设备名称
	Browser    string     `gorm:"type:varchar(64)"`                                            // 浏览器
	OS         string     `gorm:"type:varchar(32)"`                                            // 操作系统
	Status     int        `gorm:"type:int;not null"`                                           // 状态（1:成功 2:失败）
	FailReason string     `gorm:"type:varchar(128)"`                                           // 失败原因

	// 审计字段 - 与数据库表字段对应
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Version   int       `gorm:"default:1" json:"version"`
}

// TableName 指定表名
func (UserLoginLogPO) TableName() string {
	return "user_login_logs"
}

// GetID 获取实体主键
func (u *UserLoginLogPO) GetID() int {
	return u.ID
}

// LoginStatus 登录状态枚举
const (
	LoginStatusSuccess = 1 // 登录成功
	LoginStatusFailed  = 2 // 登录失败
)

// LoginType 登录方式枚举
const (
	LoginTypeUsername = "username" // 用户名+密码
	LoginTypeEmail    = "email"    // 邮箱+密码
	LoginTypePhone    = "phone"    // 手机号+验证码
	LoginTypeWechat   = "wechat"   // 微信登录
	LoginTypeAlipay   = "alipay"   // 支付宝登录
	LoginTypeQQ       = "qq"       // QQ登录
)
