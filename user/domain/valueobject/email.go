package valueobject

import (
	"strings"
	"universal-service-user/rules"
	"universal-service-user/user/domain/errors"
)

// Email 邮箱值对象
type Email string

// NewEmail 创建邮箱值对象
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	// 使用规则引擎验证邮箱格式
	result := rules.ForField("email").Email().Validate(map[string]any{
		"email": email,
	})

	if !result.IsValid() {
		return "", errors.ErrUserInvalidEmail
	}

	return Email(email), nil
}

// String 转换为字符串
func (e Email) String() string {
	return string(e)
}

// Domain 获取邮箱域名
func (e Email) Domain() string {
	parts := strings.Split(string(e), "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// LocalPart 获取邮箱本地部分
func (e Email) LocalPart() string {
	parts := strings.Split(string(e), "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}
