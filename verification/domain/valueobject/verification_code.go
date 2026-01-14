package valueobject

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/2928807938/universal-service-user/verification/domain/enum"
	verrors "github.com/2928807938/universal-service-user/verification/domain/errors"
)

// VerificationCode 验证码值对象
type VerificationCode struct {
	Code      string                 // 验证码
	Target    string                 // 目标（邮箱或手机号）
	Scene     enum.VerificationScene // 使用场景
	ExpiresAt time.Time              // 过期时间
	Verified  bool                   // 是否已验证
	CreatedAt time.Time              // 创建时间
}

// NewVerificationCode 创建新的验证码
func NewVerificationCode(target string, scene enum.VerificationScene, codeLength int, expireDuration time.Duration) (*VerificationCode, error) {
	// 验证场景是否有效
	if !scene.IsValid() {
		return nil, verrors.ErrSceneInvalid
	}

	// 验证目标是否为空
	if target == "" {
		return nil, verrors.ErrTargetInvalid
	}

	// 生成验证码
	code, err := generateCode(codeLength)
	if err != nil {
		return nil, verrors.WrapVerificationError(verrors.VerificationGenerateFailed, "生成验证码失败", err)
	}

	now := time.Now()
	return &VerificationCode{
		Code:      code,
		Target:    target,
		Scene:     scene,
		ExpiresAt: now.Add(expireDuration),
		Verified:  false,
		CreatedAt: now,
	}, nil
}

// IsExpired 验证码是否已过期
func (v *VerificationCode) IsExpired() bool {
	return time.Now().After(v.ExpiresAt)
}

// IsVerified 验证码是否已使用
func (v *VerificationCode) IsVerified() bool {
	return v.Verified
}

// MarkAsVerified 标记验证码为已使用
func (v *VerificationCode) MarkAsVerified() {
	v.Verified = true
}

// Verify 验证验证码是否正确
func (v *VerificationCode) Verify(code string) error {
	// 检查是否已过期
	if v.IsExpired() {
		return verrors.ErrCodeExpired
	}

	// 检查是否已使用
	if v.IsVerified() {
		return verrors.ErrCodeUsed
	}

	// 检查验证码是否匹配
	if v.Code != code {
		return verrors.ErrCodeIncorrect
	}

	return nil
}

// RemainingTime 获取剩余有效时间
func (v *VerificationCode) RemainingTime() time.Duration {
	if v.IsExpired() {
		return 0
	}
	return time.Until(v.ExpiresAt)
}

// generateCode 生成指定长度的数字验证码
func generateCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("验证码长度必须大于0")
	}

	// 生成随机数字验证码
	code := ""
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += num.String()
	}

	return code, nil
}
