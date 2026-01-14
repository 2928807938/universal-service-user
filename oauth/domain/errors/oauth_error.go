package errors

import (
	"fmt"

	"github.com/2928807938/universal-service-user/share/errors"
)

// ==================== OAuth 模块错误 ====================
// 错误码分段: 40xxx (40000-40999)

const (
	// OAuth 模块错误码
	OAuthUnsupportedProvider = 40001 // 不支持的登录平台
	OAuthAuthorizationFailed = 40002 // 第三方授权失败
	OAuthGetUserInfoFailed   = 40003 // 获取用户信息失败
	OAuthAccountAlreadyBound = 40004 // 账号已绑定其他用户
	OAuthBindingNotFound     = 40005 // OAuth 绑定不存在
	OAuthInvalidProvider     = 40006 // Provider 配置无效
	OAuthExchangeTokenFailed = 40007 // 换取 Token 失败
	OAuthRefreshTokenFailed  = 40008 // 刷新 Token 失败
	OAuthProviderNotEnabled  = 40009 // 平台未启用
)

// OAuthError OAuth 模块错误，继承自 AppError
type OAuthError struct {
	*errors.AppError
}

// NewOAuthError 创建 OAuth 错误
func NewOAuthError(code int, message string) *OAuthError {
	return &OAuthError{
		AppError: errors.New(code, message),
	}
}

// WrapOAuthError 包装原始错误
func WrapOAuthError(code int, message string, err error) *OAuthError {
	return &OAuthError{
		AppError: errors.Wrap(code, message, err),
	}
}

// Error 实现 error 接口
func (e *OAuthError) Error() string {
	if e.AppError.Err != nil {
		return fmt.Sprintf("[OAuth:%d] %s: %v", e.AppError.Code, e.AppError.Message, e.AppError.Err)
	}
	return fmt.Sprintf("[OAuth:%d] %s", e.AppError.Code, e.AppError.Message)
}

// ==================== OAuth 预定义错误（message 已定义） ====================

var (
	// ErrUnsupportedProvider 不支持的登录平台
	ErrUnsupportedProvider = &OAuthError{
		AppError: errors.New(OAuthUnsupportedProvider, "不支持的登录平台"),
	}

	// ErrAuthorizationFailed 第三方授权失败
	ErrAuthorizationFailed = &OAuthError{
		AppError: errors.New(OAuthAuthorizationFailed, "第三方授权失败"),
	}

	// ErrGetUserInfoFailed 获取用户信息失败
	ErrGetUserInfoFailed = &OAuthError{
		AppError: errors.New(OAuthGetUserInfoFailed, "获取用户信息失败"),
	}

	// ErrAccountAlreadyBound 账号已绑定其他用户
	ErrAccountAlreadyBound = &OAuthError{
		AppError: errors.New(OAuthAccountAlreadyBound, "账号已绑定其他用户"),
	}

	// ErrBindingNotFound OAuth 绑定不存在
	ErrBindingNotFound = &OAuthError{
		AppError: errors.New(OAuthBindingNotFound, "OAuth 绑定不存在"),
	}

	// ErrInvalidProvider Provider 配置无效
	ErrInvalidProvider = &OAuthError{
		AppError: errors.New(OAuthInvalidProvider, "Provider 配置无效"),
	}

	// ErrExchangeTokenFailed 换取 Token 失败
	ErrExchangeTokenFailed = &OAuthError{
		AppError: errors.New(OAuthExchangeTokenFailed, "换取 Token 失败"),
	}

	// ErrRefreshTokenFailed 刷新 Token 失败
	ErrRefreshTokenFailed = &OAuthError{
		AppError: errors.New(OAuthRefreshTokenFailed, "刷新 Token 失败"),
	}

	// ErrProviderNotEnabled 平台未启用
	ErrProviderNotEnabled = &OAuthError{
		AppError: errors.New(OAuthProviderNotEnabled, "该平台未启用"),
	}
)
