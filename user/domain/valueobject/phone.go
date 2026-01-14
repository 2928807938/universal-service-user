package valueobject

import (
	"strings"
	"universal-service-user/rules"
	"universal-service-user/user/domain/errors"
)

// Phone 手机号值对象
type Phone string

// NewPhone 创建手机号值对象（默认为中国大陆手机号）
func NewPhone(phone string) (Phone, error) {
	phone = strings.TrimSpace(phone)

	// 使用规则引擎验证手机号格式
	result := rules.ForField("phone").Phone().Validate(map[string]any{
		"phone": phone,
	})

	if !result.IsValid() {
		return "", errors.ErrUserInvalidPhone
	}

	return Phone(phone), nil
}

// String 转换为字符串
func (p Phone) String() string {
	return string(p)
}

// IsValid 校验手机号是否有效
func (p Phone) IsValid() bool {
	// 使用规则引擎验证
	result := rules.ForField("phone").Phone().Validate(map[string]any{
		"phone": string(p),
	})
	return result.IsValid()
}

// Region 获取国家/地区代码
func (p Phone) Region() string {
	// 目前默认返回 CN（中国大陆）
	if p.IsValid() {
		return "CN"
	}
	return ""
}

// Mask 脱敏显示（138****8000）
func (p Phone) Mask() string {
	phone := string(p)
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}
