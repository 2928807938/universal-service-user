package entity

import (
	basegorm "github.com/2928807938/universal-service-user/share/repository/gorm"
	"github.com/2928807938/universal-service-user/user/domain/enum"
	"github.com/2928807938/universal-service-user/user/domain/valueobject"
)

// User 用户实体 - 聚合根
type User struct {
	ID           int
	Username     string            // 用户名（唯一，可空）
	Nickname     string            // 昵称（用户显示名称）
	Avatar       string            // 头像URL
	Email        valueobject.Email // 邮箱（唯一，可空）
	Phone        valueobject.Phone // 手机号（唯一，可空）
	PasswordHash string            // 密码哈希（可空，第三方登录用户可能无密码）
	Status       enum.UserStatus   // 用户状态
	basegorm.AuditFields
}

// Activate 激活用户
func (u *User) Activate() {
	u.Status = enum.UserStatusActive
	u.Touch()
}

// Disable 禁用用户
func (u *User) Disable() {
	u.Status = enum.UserStatusDisabled
	u.Touch()
}

// ChangePassword 修改密码
func (u *User) ChangePassword(newPasswordHash string) {
	u.PasswordHash = newPasswordHash
	u.Touch()
}

// UpdateEmail 更新邮箱
func (u *User) UpdateEmail(email valueobject.Email) {
	u.Email = email
	u.Touch()
}

// UpdatePhone 更新手机号
func (u *User) UpdatePhone(phone valueobject.Phone) {
	u.Phone = phone
	u.Touch()
}

// UpdateNickname 更新昵称
func (u *User) UpdateNickname(nickname string) {
	u.Nickname = nickname
	u.Touch()
}

// UpdateAvatar 更新头像
func (u *User) UpdateAvatar(avatar string) {
	u.Avatar = avatar
	u.Touch()
}

// Lock 锁定用户
func (u *User) Lock() {
	u.Status = enum.UserStatusLocked
	u.Touch()
}

// Unlock 解锁用户
func (u *User) Unlock() {
	u.Status = enum.UserStatusActive
	u.Touch()
}

// IsActive 判断用户是否激活
func (u *User) IsActive() bool {
	return u.Status == enum.UserStatusActive
}

// IsLocked 判断用户是否被锁定
func (u *User) IsLocked() bool {
	return u.Status == enum.UserStatusLocked
}

// IsDisabled 判断用户是否被禁用
func (u *User) IsDisabled() bool {
	return u.Status == enum.UserStatusDisabled
}

// CanLogin 判断用户是否可以登录
func (u *User) CanLogin() bool {
	return u.Status == enum.UserStatusActive
}
