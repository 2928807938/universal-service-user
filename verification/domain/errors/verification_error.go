package errors

import (
	"fmt"

	"universal-service-user/share/errors"
)

// ==================== Verification 模块错误 ====================
// 错误码分段: 20xxx (20000-20999)
// 对应设计文档 10.5 验证码相关错误码 (20xxx)

const (
	// Verification 模块错误码
	VerificationCodeIncorrect     = 20001 // 验证码错误
	VerificationCodeExpired       = 20002 // 验证码已过期
	VerificationCodeUsed          = 20003 // 验证码已使用
	VerificationCodeTooFrequent   = 20004 // 验证码发送过于频繁
	VerificationCodeNotFound      = 20005 // 验证码不存在
	VerificationSceneInvalid      = 20006 // 验证码场景无效
	VerificationTargetInvalid     = 20007 // 验证目标无效
	VerificationGenerateFailed    = 20008 // 验证码生成失败
	VerificationStoreFailed       = 20009 // 验证码存储失败
)

// VerificationError Verification 模块错误，继承自 AppError
type VerificationError struct {
	*errors.AppError
}

// NewVerificationError 创建 Verification 错误
func NewVerificationError(code int, message string) *VerificationError {
	return &VerificationError{
		AppError: errors.New(code, message),
	}
}

// WrapVerificationError 包装原始错误
func WrapVerificationError(code int, message string, err error) *VerificationError {
	return &VerificationError{
		AppError: errors.Wrap(code, message, err),
	}
}

// Error 实现 error 接口
func (e *VerificationError) Error() string {
	if e.AppError.Err != nil {
		return fmt.Sprintf("[Verification:%d] %s: %v", e.AppError.Code, e.AppError.Message, e.AppError.Err)
	}
	return fmt.Sprintf("[Verification:%d] %s", e.AppError.Code, e.AppError.Message)
}

// ==================== Verification 预定义错误（message 已定义） ====================

var (
	// ErrCodeIncorrect 验证码错误
	ErrCodeIncorrect = &VerificationError{
		AppError: errors.New(VerificationCodeIncorrect, "验证码错误"),
	}

	// ErrCodeExpired 验证码已过期
	ErrCodeExpired = &VerificationError{
		AppError: errors.New(VerificationCodeExpired, "验证码已过期"),
	}

	// ErrCodeUsed 验证码已使用
	ErrCodeUsed = &VerificationError{
		AppError: errors.New(VerificationCodeUsed, "验证码已使用"),
	}

	// ErrCodeTooFrequent 验证码发送过于频繁
	ErrCodeTooFrequent = &VerificationError{
		AppError: errors.New(VerificationCodeTooFrequent, "验证码发送过于频繁，请稍后再试"),
	}

	// ErrCodeNotFound 验证码不存在
	ErrCodeNotFound = &VerificationError{
		AppError: errors.New(VerificationCodeNotFound, "验证码不存在或已过期"),
	}

	// ErrSceneInvalid 验证码场景无效
	ErrSceneInvalid = &VerificationError{
		AppError: errors.New(VerificationSceneInvalid, "验证码场景无效"),
	}

	// ErrTargetInvalid 验证目标无效
	ErrTargetInvalid = &VerificationError{
		AppError: errors.New(VerificationTargetInvalid, "验证目标无效（邮箱或手机号格式错误）"),
	}

	// ErrGenerateFailed 验证码生成失败
	ErrGenerateFailed = &VerificationError{
		AppError: errors.New(VerificationGenerateFailed, "验证码生成失败"),
	}

	// ErrStoreFailed 验证码存储失败
	ErrStoreFailed = &VerificationError{
		AppError: errors.New(VerificationStoreFailed, "验证码存储失败"),
	}
)