package errors

import (
	"fmt"

	"universal-service-user/share/errors"
)

// ==================== Auth 模块错误 ====================
// 错误码分段: 30xxx (30000-30999)

const (
	// Auth 模块错误码 - 基础认证错误 (30001-30099)
	AuthNotLoggedIn        = 30001 // 未登录
	AuthTokenExpired       = 30002 // Token 已过期
	AuthTokenInvalid       = 30003 // Token 无效
	AuthTokenRevoked       = 30004 // Token 已被撤销
	AuthRefreshTokenFailed = 30005 // 刷新 Token 失败
	AuthTokenGenerateFailed = 30006 // 生成 Token 失败

	// Auth 模块错误码 - 登录防刷 (30501-30599)
	AuthIPRateLimited      = 30501 // IP 请求过于频繁
	AuthAccountLocked      = 30502 // 登录失败次数过多，账号已锁定
	AuthDeviceRateLimited  = 30503 // 设备请求过于频繁
	AuthCaptchaRequired    = 30504 // 需要验证码
	AuthCaptchaInvalid     = 30505 // 验证码错误

	// Auth 模块错误码 - 会话管理 (30601-30699)
	AuthSessionNotFound    = 30601 // 会话不存在
	AuthSessionExpired     = 30602 // 会话已过期
	AuthInvalidDeviceInfo  = 30603 // 设备信息无效
)

// AuthError Auth 模块错误，继承自 AppError
type AuthError struct {
	*errors.AppError
}

// NewAuthError 创建 Auth 错误
func NewAuthError(code int, message string) *AuthError {
	return &AuthError{
		AppError: errors.New(code, message),
	}
}

// WrapAuthError 包装原始错误
func WrapAuthError(code int, message string, err error) *AuthError {
	return &AuthError{
		AppError: errors.Wrap(code, message, err),
	}
}

// Error 实现 error 接口
func (e *AuthError) Error() string {
	if e.AppError.Err != nil {
		return fmt.Sprintf("[Auth:%d] %s: %v", e.AppError.Code, e.AppError.Message, e.AppError.Err)
	}
	return fmt.Sprintf("[Auth:%d] %s", e.AppError.Code, e.AppError.Message)
}

// ==================== Auth 预定义错误（message 已定义） ====================

var (
	// 基础认证错误
	ErrNotLoggedIn = &AuthError{
		AppError: errors.New(AuthNotLoggedIn, "未登录"),
	}

	ErrTokenExpired = &AuthError{
		AppError: errors.New(AuthTokenExpired, "Token 已过期"),
	}

	ErrTokenInvalid = &AuthError{
		AppError: errors.New(AuthTokenInvalid, "Token 无效"),
	}

	ErrTokenRevoked = &AuthError{
		AppError: errors.New(AuthTokenRevoked, "Token 已被撤销"),
	}

	ErrRefreshTokenFailed = &AuthError{
		AppError: errors.New(AuthRefreshTokenFailed, "刷新 Token 失败"),
	}

	ErrTokenGenerateFailed = &AuthError{
		AppError: errors.New(AuthTokenGenerateFailed, "生成 Token 失败"),
	}

	// 登录防刷错误
	ErrIPRateLimited = &AuthError{
		AppError: errors.New(AuthIPRateLimited, "IP 请求过于频繁，请稍后重试"),
	}

	ErrAccountLocked = &AuthError{
		AppError: errors.New(AuthAccountLocked, "登录失败次数过多，账号已被锁定"),
	}

	ErrDeviceRateLimited = &AuthError{
		AppError: errors.New(AuthDeviceRateLimited, "设备请求过于频繁，请稍后重试"),
	}

	ErrCaptchaRequired = &AuthError{
		AppError: errors.New(AuthCaptchaRequired, "需要验证码"),
	}

	ErrCaptchaInvalid = &AuthError{
		AppError: errors.New(AuthCaptchaInvalid, "验证码错误"),
	}

	// 会话管理错误
	ErrSessionNotFound = &AuthError{
		AppError: errors.New(AuthSessionNotFound, "会话不存在"),
	}

	ErrSessionExpired = &AuthError{
		AppError: errors.New(AuthSessionExpired, "会话已过期"),
	}

	ErrInvalidDeviceInfo = &AuthError{
		AppError: errors.New(AuthInvalidDeviceInfo, "设备信息无效"),
	}
)
