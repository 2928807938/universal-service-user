package enum

// UserStatus 用户状态枚举
type UserStatus int

const (
	UserStatusInactive UserStatus = 0 // 未激活
	UserStatusActive   UserStatus = 1 // 已激活
	UserStatusDisabled UserStatus = 2 // 已禁用
	UserStatusLocked   UserStatus = 3 // 已锁定（多次登录失败）
)

// String 返回状态描述
func (s UserStatus) String() string {
	switch s {
	case UserStatusInactive:
		return "inactive"
	case UserStatusActive:
		return "active"
	case UserStatusDisabled:
		return "disabled"
	case UserStatusLocked:
		return "locked"
	default:
		return "unknown"
	}
}

// IsValid 验证状态是否有效
func (s UserStatus) IsValid() bool {
	return s >= UserStatusInactive && s <= UserStatusLocked
}

// IsActive 是否激活状态
func (s UserStatus) IsActive() bool {
	return s == UserStatusActive
}

// IsInactive 是否未激活状态
func (s UserStatus) IsInactive() bool {
	return s == UserStatusInactive
}

// IsDisabled 是否禁用状态
func (s UserStatus) IsDisabled() bool {
	return s == UserStatusDisabled
}

// IsLocked 是否锁定状态
func (s UserStatus) IsLocked() bool {
	return s == UserStatusLocked
}
