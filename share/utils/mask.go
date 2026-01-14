package utils

import (
	"strings"
	"unicode/utf8"
)

// MaskEmail 邮箱脱敏
// 示例:
//   - "user@example.com" -> "use***@example.com"
//   - "a@example.com" -> "a***@example.com"
//   - "" -> ""
func MaskEmail(email string) string {
	if email == "" {
		return ""
	}

	// 分割邮箱
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		// 邮箱格式不正确，返回原值
		return email
	}

	username := parts[0]
	domain := parts[1]

	if username == "" || domain == "" {
		return email
	}

	// 获取用户名字符数
	usernameLen := utf8.RuneCountInString(username)

	// 根据长度决定脱敏方式
	if usernameLen <= 3 {
		// 用户名长度 <= 3，保留第一个字符
		if usernameLen == 1 {
			return username[:1] + "***@" + domain
		}
		firstChar := string([]rune(username)[0])
		return firstChar + "***@" + domain
	}

	// 用户名长度 > 3，保留前3个字符
	runes := []rune(username)
	maskedUsername := string(runes[:3]) + "***"

	return maskedUsername + "@" + domain
}

// MaskPhone 手机号脱敏
// 示例:
//   - "13812345678" -> "138****5678"
//   - "" -> ""
func MaskPhone(phone string) string {
	if phone == "" {
		return ""
	}

	phoneLen := utf8.RuneCountInString(phone)

	// 标准中国大陆手机号长度为11位
	if phoneLen != 11 {
		// 非标准长度，不脱敏
		return phone
	}

	runes := []rune(phone)
	// 保留前3位和后4位，中间4位用****代替
	return string(runes[:3]) + "****" + string(runes[7:])
}

// MaskString 通用字符串脱敏
// prefixLen: 前缀保留长度
// suffixLen: 后缀保留长度
// maskChar: 脱敏字符（默认*）
// 示例:
//   - MaskString("12345678", 2, 2, "*") -> "12****78"
func MaskString(str string, prefixLen, suffixLen int, maskChar string) string {
	if str == "" {
		return ""
	}

	strLen := utf8.RuneCountInString(str)

	// 如果字符串长度小于等于前缀+后缀长度，不脱敏
	if strLen <= prefixLen+suffixLen {
		return str
	}

	runes := []rune(str)

	// 构建脱敏字符串
	maskLen := strLen - prefixLen - suffixLen
	mask := strings.Repeat(maskChar, maskLen)

	return string(runes[:prefixLen]) + mask + string(runes[strLen-suffixLen:])
}

// MaskIDCard 身份证号脱敏
// 示例:
//   - "110101199001011234" -> "110101********1234"
func MaskIDCard(idCard string) string {
	if idCard == "" {
		return ""
	}

	idLen := utf8.RuneCountInString(idCard)

	// 15位或18位身份证
	if idLen == 15 {
		return MaskString(idCard, 6, 3, "*")
	} else if idLen == 18 {
		return MaskString(idCard, 6, 4, "*")
	}

	// 其他长度不脱敏
	return idCard
}

// MaskBankCard 银行卡号脱敏
// 示例:
//   - "6222021234567890123" -> "622202***********0123"
func MaskBankCard(bankCard string) string {
	if bankCard == "" {
		return ""
	}

	cardLen := utf8.RuneCountInString(bankCard)

	if cardLen <= 10 {
		// 长度太短，保留前4后3
		return MaskString(bankCard, 4, 3, "*")
	}

	// 保留前6位和后4位
	return MaskString(bankCard, 6, 4, "*")
}
