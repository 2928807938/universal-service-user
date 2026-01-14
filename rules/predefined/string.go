package predefined

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/2928807938/universal-service-user/rules/core"
)

// EmailRule 邮箱格式验证规则
type EmailRule struct {
	pattern *regexp.Regexp
}

func NewEmailRule() *EmailRule {
	// 简单但有效的邮箱正则
	pattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return &EmailRule{pattern: pattern}
}

func (r *EmailRule) Name() string {
	return "email"
}

func (r *EmailRule) Description() string {
	return "Field must be a valid email address"
}

func (r *EmailRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success() // nil 值由 required 规则处理
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success() // 空字符串由 required 规则处理
	}

	if !r.pattern.MatchString(str) {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a valid email address")
	}

	return core.Success()
}

// LengthRule 字符串长度范围验证规则
type LengthRule struct {
	min int
	max int
}

func NewLengthRule(min, max int) *LengthRule {
	return &LengthRule{min: min, max: max}
}

func (r *LengthRule) Name() string {
	return "length"
}

func (r *LengthRule) Description() string {
	return fmt.Sprintf("Field length must be between %d and %d", r.min, r.max)
}

func (r *LengthRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	length := utf8.RuneCountInString(str)
	if length < r.min || length > r.max {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field length must be between %d and %d, got %d", r.min, r.max, length))
	}

	return core.Success()
}

// MinLengthRule 最小长度验证规则
type MinLengthRule struct {
	min int
}

func NewMinLengthRule(min int) *MinLengthRule {
	return &MinLengthRule{min: min}
}

func (r *MinLengthRule) Name() string {
	return "min_length"
}

func (r *MinLengthRule) Description() string {
	return fmt.Sprintf("Field length must be at least %d", r.min)
}

func (r *MinLengthRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	length := utf8.RuneCountInString(str)
	if length < r.min {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field length must be at least %d, got %d", r.min, length))
	}

	return core.Success()
}

// MaxLengthRule 最大长度验证规则
type MaxLengthRule struct {
	max int
}

func NewMaxLengthRule(max int) *MaxLengthRule {
	return &MaxLengthRule{max: max}
}

func (r *MaxLengthRule) Name() string {
	return "max_length"
}

func (r *MaxLengthRule) Description() string {
	return fmt.Sprintf("Field length must be at most %d", r.max)
}

func (r *MaxLengthRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	length := utf8.RuneCountInString(str)
	if length > r.max {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field length must be at most %d, got %d", r.max, length))
	}

	return core.Success()
}

// PatternRule 正则表达式验证规则
type PatternRule struct {
	pattern *regexp.Regexp
	message string
}

func NewPatternRule(pattern string) (*PatternRule, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}
	return &PatternRule{
		pattern: re,
		message: fmt.Sprintf("field must match pattern: %s", pattern),
	}, nil
}

func MustPatternRule(pattern string) *PatternRule {
	r, err := NewPatternRule(pattern)
	if err != nil {
		panic(err)
	}
	return r
}

func (r *PatternRule) WithMessage(message string) *PatternRule {
	r.message = message
	return r
}

func (r *PatternRule) Name() string {
	return "pattern"
}

func (r *PatternRule) Description() string {
	return "Field must match the specified pattern"
}

func (r *PatternRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	if !r.pattern.MatchString(str) {
		return core.Failure(ctx.FieldName(), r.Name(), r.message)
	}

	return core.Success()
}

// URLRule URL格式验证规则
type URLRule struct{}

func NewURLRule() *URLRule {
	return &URLRule{}
}

func (r *URLRule) Name() string {
	return "url"
}

func (r *URLRule) Description() string {
	return "Field must be a valid URL"
}

func (r *URLRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	u, err := url.Parse(str)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a valid URL")
	}

	return core.Success()
}

// ContainsRule 包含子串验证规则
type ContainsRule struct {
	substring string
}

func NewContainsRule(substring string) *ContainsRule {
	return &ContainsRule{substring: substring}
}

func (r *ContainsRule) Name() string {
	return "contains"
}

func (r *ContainsRule) Description() string {
	return fmt.Sprintf("Field must contain: %s", r.substring)
}

func (r *ContainsRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	if !strings.Contains(str, r.substring) {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must contain: %s", r.substring))
	}

	return core.Success()
}

// StartsWithRule 前缀验证规则
type StartsWithRule struct {
	prefix string
}

func NewStartsWithRule(prefix string) *StartsWithRule {
	return &StartsWithRule{prefix: prefix}
}

func (r *StartsWithRule) Name() string {
	return "starts_with"
}

func (r *StartsWithRule) Description() string {
	return fmt.Sprintf("Field must start with: %s", r.prefix)
}

func (r *StartsWithRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	if !strings.HasPrefix(str, r.prefix) {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must start with: %s", r.prefix))
	}

	return core.Success()
}

// EndsWithRule 后缀验证规则
type EndsWithRule struct {
	suffix string
}

func NewEndsWithRule(suffix string) *EndsWithRule {
	return &EndsWithRule{suffix: suffix}
}

func (r *EndsWithRule) Name() string {
	return "ends_with"
}

func (r *EndsWithRule) Description() string {
	return fmt.Sprintf("Field must end with: %s", r.suffix)
}

func (r *EndsWithRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	if !strings.HasSuffix(str, r.suffix) {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must end with: %s", r.suffix))
	}

	return core.Success()
}

// AlphaRule 只包含字母验证规则
type AlphaRule struct {
	pattern *regexp.Regexp
}

func NewAlphaRule() *AlphaRule {
	return &AlphaRule{pattern: regexp.MustCompile(`^[a-zA-Z]+$`)}
}

func (r *AlphaRule) Name() string {
	return "alpha"
}

func (r *AlphaRule) Description() string {
	return "Field must contain only letters"
}

func (r *AlphaRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	if !r.pattern.MatchString(str) {
		return core.Failure(ctx.FieldName(), r.Name(), "field must contain only letters")
	}

	return core.Success()
}

// AlphanumericRule 只包含字母和数字验证规则
type AlphanumericRule struct {
	pattern *regexp.Regexp
}

func NewAlphanumericRule() *AlphanumericRule {
	return &AlphanumericRule{pattern: regexp.MustCompile(`^[a-zA-Z0-9]+$`)}
}

func (r *AlphanumericRule) Name() string {
	return "alphanumeric"
}

func (r *AlphanumericRule) Description() string {
	return "Field must contain only letters and numbers"
}

func (r *AlphanumericRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	if !r.pattern.MatchString(str) {
		return core.Failure(ctx.FieldName(), r.Name(), "field must contain only letters and numbers")
	}

	return core.Success()
}

// NumericRule 只包含数字字符验证规则
type NumericRule struct {
	pattern *regexp.Regexp
}

func NewNumericRule() *NumericRule {
	return &NumericRule{pattern: regexp.MustCompile(`^[0-9]+$`)}
}

func (r *NumericRule) Name() string {
	return "numeric"
}

func (r *NumericRule) Description() string {
	return "Field must contain only numeric characters"
}

func (r *NumericRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	if !r.pattern.MatchString(str) {
		return core.Failure(ctx.FieldName(), r.Name(), "field must contain only numeric characters")
	}

	return core.Success()
}

// UUIDRule UUID格式验证规则
type UUIDRule struct {
	pattern *regexp.Regexp
}

func NewUUIDRule() *UUIDRule {
	// 支持 UUID v1-v5 格式
	pattern := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	return &UUIDRule{pattern: pattern}
}

func (r *UUIDRule) Name() string {
	return "uuid"
}

func (r *UUIDRule) Description() string {
	return "Field must be a valid UUID"
}

func (r *UUIDRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	if !r.pattern.MatchString(str) {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a valid UUID")
	}

	return core.Success()
}

// IPRule IP地址验证规则
type IPRule struct {
	version int // 0: both, 4: IPv4 only, 6: IPv6 only
}

func NewIPRule() *IPRule {
	return &IPRule{version: 0}
}

func NewIPv4Rule() *IPRule {
	return &IPRule{version: 4}
}

func NewIPv6Rule() *IPRule {
	return &IPRule{version: 6}
}

func (r *IPRule) Name() string {
	switch r.version {
	case 4:
		return "ipv4"
	case 6:
		return "ipv6"
	default:
		return "ip"
	}
}

func (r *IPRule) Description() string {
	switch r.version {
	case 4:
		return "Field must be a valid IPv4 address"
	case 6:
		return "Field must be a valid IPv6 address"
	default:
		return "Field must be a valid IP address"
	}
}

func (r *IPRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	// IPv4 正则
	ipv4Pattern := `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	// IPv6 简化正则
	ipv6Pattern := `^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$|^::$|^([0-9a-fA-F]{1,4}:){1,7}:$|^:(:([0-9a-fA-F]{1,4})){1,7}$`

	isIPv4 := regexp.MustCompile(ipv4Pattern).MatchString(str)
	isIPv6 := regexp.MustCompile(ipv6Pattern).MatchString(str)

	switch r.version {
	case 4:
		if !isIPv4 {
			return core.Failure(ctx.FieldName(), r.Name(), "field must be a valid IPv4 address")
		}
	case 6:
		if !isIPv6 {
			return core.Failure(ctx.FieldName(), r.Name(), "field must be a valid IPv6 address")
		}
	default:
		if !isIPv4 && !isIPv6 {
			return core.Failure(ctx.FieldName(), r.Name(), "field must be a valid IP address")
		}
	}

	return core.Success()
}

// PhoneRule 手机号验证规则
type PhoneRule struct {
	pattern *regexp.Regexp
	region  string
}

// NewPhoneRule 创建中国大陆手机号验证规则
func NewPhoneRule() *PhoneRule {
	// 中国大陆手机号
	pattern := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return &PhoneRule{pattern: pattern, region: "CN"}
}

// NewPhoneRuleWithPattern 创建自定义手机号验证规则
func NewPhoneRuleWithPattern(pattern string, region string) (*PhoneRule, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid phone pattern: %w", err)
	}
	return &PhoneRule{pattern: re, region: region}, nil
}

func (r *PhoneRule) Name() string {
	return "phone"
}

func (r *PhoneRule) Description() string {
	return fmt.Sprintf("Field must be a valid phone number (%s)", r.region)
}

func (r *PhoneRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	if !r.pattern.MatchString(str) {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a valid phone number")
	}

	return core.Success()
}

// JSONRule JSON格式验证规则
type JSONRule struct{}

func NewJSONRule() *JSONRule {
	return &JSONRule{}
}

func (r *JSONRule) Name() string {
	return "json"
}

func (r *JSONRule) Description() string {
	return "Field must be a valid JSON string"
}

func (r *JSONRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	// 简单的 JSON 验证：检查首尾字符
	str = strings.TrimSpace(str)
	if len(str) < 2 {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a valid JSON string")
	}

	first, last := str[0], str[len(str)-1]
	if !((first == '{' && last == '}') || (first == '[' && last == ']') ||
		(first == '"' && last == '"')) {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a valid JSON string")
	}

	return core.Success()
}

// LowercaseRule 小写字母验证规则
type LowercaseRule struct{}

func NewLowercaseRule() *LowercaseRule {
	return &LowercaseRule{}
}

func (r *LowercaseRule) Name() string {
	return "lowercase"
}

func (r *LowercaseRule) Description() string {
	return "Field must be lowercase"
}

func (r *LowercaseRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	if str != strings.ToLower(str) {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be lowercase")
	}

	return core.Success()
}

// UppercaseRule 大写字母验证规则
type UppercaseRule struct{}

func NewUppercaseRule() *UppercaseRule {
	return &UppercaseRule{}
}

func (r *UppercaseRule) Name() string {
	return "uppercase"
}

func (r *UppercaseRule) Description() string {
	return "Field must be uppercase"
}

func (r *UppercaseRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	str, ok := value.(string)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a string")
	}

	if str == "" {
		return core.Success()
	}

	if str != strings.ToUpper(str) {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be uppercase")
	}

	return core.Success()
}
