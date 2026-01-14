package errors

import (
	"fmt"
	"time"

	"universal-service-user/share/errors"
)

// ==================== User 模块错误 ====================
// 错误码分段: 11xxx (11000-11999)

const (
	// User 模块错误码
	UserNotFound              = 11001 // 用户不存在
	UserAlreadyExists         = 11002 // 用户已存在
	UserInvalidPassword       = 11003 // 密码错误
	UserDisabled              = 11004 // 用户已被禁用
	UserExpired               = 11005 // 用户已过期
	UserInvalidEmail          = 11006 // 邮箱格式错误
	UserEmailAlreadyExists    = 11007 // 邮箱已被使用
	UserUsernameExists        = 11008 // 用户名已被使用
	UserValidationFailed      = 11009 // 参数验证失败
	UserPasswordStrengthWeak  = 11010 // 密码强度不够
	UserPasswordMismatch      = 11011 // 密码不匹配
	UserInvalidPhone          = 11012 // 手机号格式错误
	UserPhoneAlreadyExists    = 11013 // 手机号已被使用
	UserLocked                = 11014 // 用户已被锁定
	UserStatusAbnormal        = 11015 // 用户状态异常
	UserInactive              = 11016 // 用户未激活
	UserAccountTempLocked     = 11017 // 账号临时锁定（登录失败次数过多）
	UserLoginAttemptExceeded  = 11018 // 登录尝试次数过多
)

// UserError User 模块错误，继承自 AppError
type UserError struct {
	*errors.AppError
}

// NewUserError 创建 User 错误
func NewUserError(code int, message string) *UserError {
	return &UserError{
		AppError: errors.New(code, message),
	}
}

// WrapUserError 包装原始错误
func WrapUserError(code int, message string, err error) *UserError {
	return &UserError{
		AppError: errors.Wrap(code, message, err),
	}
}

// Error 实现 error 接口
func (e *UserError) Error() string {
	if e.AppError.Err != nil {
		return fmt.Sprintf("[User:%d] %s: %v", e.AppError.Code, e.AppError.Message, e.AppError.Err)
	}
	return fmt.Sprintf("[User:%d] %s", e.AppError.Code, e.AppError.Message)
}

// NewValidationError 创建参数验证错误
func NewValidationError(message string) *UserError {
	return NewUserError(UserValidationFailed, message)
}

// ==================== User 预定义错误（message 已定义） ====================

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = &UserError{
		AppError: errors.New(UserNotFound, "用户不存在"),
	}

	// ErrUserAlreadyExists 用户已存在
	ErrUserAlreadyExists = &UserError{
		AppError: errors.New(UserAlreadyExists, "用户已存在"),
	}

	// ErrUserInvalidPassword 密码错误
	ErrUserInvalidPassword = &UserError{
		AppError: errors.New(UserInvalidPassword, "密码错误"),
	}

	// ErrUserDisabled 用户已被禁用
	ErrUserDisabled = &UserError{
		AppError: errors.New(UserDisabled, "用户已被禁用"),
	}

	// ErrUserExpired 用户已过期
	ErrUserExpired = &UserError{
		AppError: errors.New(UserExpired, "用户已过期"),
	}

	// ErrUserInvalidEmail 邮箱格式不正确
	ErrUserInvalidEmail = &UserError{
		AppError: errors.New(UserInvalidEmail, "邮箱格式不正确"),
	}

	// ErrUserEmailAlreadyExists 邮箱已被使用
	ErrUserEmailAlreadyExists = &UserError{
		AppError: errors.New(UserEmailAlreadyExists, "邮箱已被使用"),
	}

	// ErrUserUsernameExists 用户名已被使用
	ErrUserUsernameExists = &UserError{
		AppError: errors.New(UserUsernameExists, "用户名已被使用"),
	}

	// ErrPasswordStrengthWeak 密码强度不够
	ErrPasswordStrengthWeak = &UserError{
		AppError: errors.New(UserPasswordStrengthWeak, "密码强度不够"),
	}

	// ErrPasswordMismatch 密码不匹配
	ErrPasswordMismatch = &UserError{
		AppError: errors.New(UserPasswordMismatch, "密码不匹配"),
	}

	// ErrUserInvalidPhone 手机号格式不正确
	ErrUserInvalidPhone = &UserError{
		AppError: errors.New(UserInvalidPhone, "手机号格式不正确"),
	}

	// ErrUserPhoneAlreadyExists 手机号已被使用
	ErrUserPhoneAlreadyExists = &UserError{
		AppError: errors.New(UserPhoneAlreadyExists, "手机号已被使用"),
	}

	// ErrUserLocked 用户已被锁定
	ErrUserLocked = &UserError{
		AppError: errors.New(UserLocked, "账号已被锁定，请联系管理员"),
	}

	// ErrUserStatusAbnormal 用户状态异常
	ErrUserStatusAbnormal = &UserError{
		AppError: errors.New(UserStatusAbnormal, "账号状态异常，无法登录"),
	}

	// ErrUserInactive 用户未激活
	ErrUserInactive = &UserError{
		AppError: errors.New(UserInactive, "用户未激活"),
	}
)

// NewLoginAttemptError 创建登录尝试次数过多错误
// remainingAttempts: 剩余尝试次数
func NewLoginAttemptError(remainingAttempts int) *UserError {
	return &UserError{
		AppError: errors.New(UserLoginAttemptExceeded, fmt.Sprintf("密码错误，还剩 %d 次尝试机会", remainingAttempts)),
	}
}

// NewAccountTempLockedError 创建账号临时锁定错误
// lockDuration: 锁定时长
func NewAccountTempLockedError(lockDuration time.Duration) *UserError {
	minutes := int(lockDuration.Minutes())
	message := fmt.Sprintf("登录失败次数过多，账号已锁定 %d 分钟", minutes)
	if minutes < 1 {
		message = fmt.Sprintf("登录失败次数过多，账号已锁定 %d 秒", int(lockDuration.Seconds()))
	}
	return &UserError{
		AppError: errors.New(UserAccountTempLocked, message),
	}
}
