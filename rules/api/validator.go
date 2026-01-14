package api

import (
	"reflect"
	"time"

	"universal-service-user/rules/core"
	"universal-service-user/rules/predefined"
)

// Validator 提供链式调用的验证 API
type Validator struct {
	fieldRules   map[string][]core.Rule // 按字段分组的规则
	currentField string                 // 当前正在配置的字段
	fieldOrder   []string               // 字段添加顺序
}

// New 创建一个新的验证器
func New() *Validator {
	return &Validator{
		fieldRules: make(map[string][]core.Rule),
		fieldOrder: make([]string, 0),
	}
}

// ForField 指定要验证的字段
func (v *Validator) ForField(fieldName string) *Validator {
	v.currentField = fieldName
	if _, exists := v.fieldRules[fieldName]; !exists {
		v.fieldRules[fieldName] = make([]core.Rule, 0)
		v.fieldOrder = append(v.fieldOrder, fieldName)
	}
	return v
}

// addRule 添加规则到当前字段
func (v *Validator) addRule(rule core.Rule) *Validator {
	if v.currentField != "" {
		v.fieldRules[v.currentField] = append(v.fieldRules[v.currentField], rule)
	}
	return v
}

// Rule 添加一个规则（用于添加自定义规则）
func (v *Validator) Rule(rule core.Rule) *Validator {
	return v.addRule(rule)
}

// ========== 基础规则 ==========

// Required 必填验证
func (v *Validator) Required() *Validator {
	return v.addRule(predefined.NewRequiredRule())
}

// NotNull 非空验证
func (v *Validator) NotNull() *Validator {
	return v.addRule(predefined.NewNotNullRule())
}

// NotEmpty 非空字符串验证
func (v *Validator) NotEmpty() *Validator {
	return v.addRule(predefined.NewNotEmptyRule())
}

// ========== 字符串规则 ==========

// Email 邮箱格式验证
func (v *Validator) Email() *Validator {
	return v.addRule(predefined.NewEmailRule())
}

// Length 字符串长度范围验证
func (v *Validator) Length(min, max int) *Validator {
	return v.addRule(predefined.NewLengthRule(min, max))
}

// MinLength 最小长度验证
func (v *Validator) MinLength(min int) *Validator {
	return v.addRule(predefined.NewMinLengthRule(min))
}

// MaxLength 最大长度验证
func (v *Validator) MaxLength(max int) *Validator {
	return v.addRule(predefined.NewMaxLengthRule(max))
}

// Pattern 正则表达式验证
func (v *Validator) Pattern(pattern string) *Validator {
	rule, err := predefined.NewPatternRule(pattern)
	if err != nil {
		return v.addRule(&invalidPatternRule{pattern: pattern, err: err})
	}
	return v.addRule(rule)
}

// PatternWithMessage 带自定义消息的正则表达式验证
func (v *Validator) PatternWithMessage(pattern, message string) *Validator {
	rule, err := predefined.NewPatternRule(pattern)
	if err != nil {
		return v.addRule(&invalidPatternRule{pattern: pattern, err: err})
	}
	return v.addRule(rule.WithMessage(message))
}

// URL URL格式验证
func (v *Validator) URL() *Validator {
	return v.addRule(predefined.NewURLRule())
}

// Contains 包含子串验证
func (v *Validator) Contains(substring string) *Validator {
	return v.addRule(predefined.NewContainsRule(substring))
}

// StartsWith 前缀验证
func (v *Validator) StartsWith(prefix string) *Validator {
	return v.addRule(predefined.NewStartsWithRule(prefix))
}

// EndsWith 后缀验证
func (v *Validator) EndsWith(suffix string) *Validator {
	return v.addRule(predefined.NewEndsWithRule(suffix))
}

// ========== 数值规则 ==========

// Min 最小值验证
func (v *Validator) Min(min float64) *Validator {
	return v.addRule(predefined.NewMinRule(min))
}

// Max 最大值验证
func (v *Validator) Max(max float64) *Validator {
	return v.addRule(predefined.NewMaxRule(max))
}

// Range 数值范围验证
func (v *Validator) Range(min, max float64) *Validator {
	return v.addRule(predefined.NewRangeRule(min, max))
}

// Positive 正数验证
func (v *Validator) Positive() *Validator {
	return v.addRule(predefined.NewPositiveRule())
}

// Negative 负数验证
func (v *Validator) Negative() *Validator {
	return v.addRule(predefined.NewNegativeRule())
}

// NonNegative 非负数验证
func (v *Validator) NonNegative() *Validator {
	return v.addRule(predefined.NewNonNegativeRule())
}

// ========== 集合规则 ==========

// In 值在指定集合中验证
func (v *Validator) In(values ...any) *Validator {
	return v.addRule(predefined.NewInRule(values...))
}

// NotIn 值不在指定集合中验证
func (v *Validator) NotIn(values ...any) *Validator {
	return v.addRule(predefined.NewNotInRule(values...))
}

// ArrayLength 数组长度范围验证
func (v *Validator) ArrayLength(min, max int) *Validator {
	return v.addRule(predefined.NewArrayLengthRule(min, max))
}

// MinArrayLength 数组最小长度验证
func (v *Validator) MinArrayLength(min int) *Validator {
	return v.addRule(predefined.NewMinArrayLengthRule(min))
}

// MaxArrayLength 数组最大长度验证
func (v *Validator) MaxArrayLength(max int) *Validator {
	return v.addRule(predefined.NewMaxArrayLengthRule(max))
}

// Unique 数组元素唯一性验证
func (v *Validator) Unique() *Validator {
	return v.addRule(predefined.NewUniqueRule())
}

// ========== 时间规则 ==========

// Before 早于指定时间验证
func (v *Validator) Before(t time.Time) *Validator {
	return v.addRule(predefined.NewBeforeRule(t))
}

// BeforeNow 早于当前时间验证
func (v *Validator) BeforeNow() *Validator {
	return v.addRule(predefined.NewBeforeNowRule())
}

// After 晚于指定时间验证
func (v *Validator) After(t time.Time) *Validator {
	return v.addRule(predefined.NewAfterRule(t))
}

// AfterNow 晚于当前时间验证
func (v *Validator) AfterNow() *Validator {
	return v.addRule(predefined.NewAfterNowRule())
}

// BetweenTime 时间范围验证
func (v *Validator) BetweenTime(start, end time.Time) *Validator {
	return v.addRule(predefined.NewBetweenTimeRule(start, end))
}

// DateFormat 日期格式验证
func (v *Validator) DateFormat(format string) *Validator {
	return v.addRule(predefined.NewDateFormatRule(format))
}

// ========== 更多字符串规则 ==========

// Alpha 只包含字母验证
func (v *Validator) Alpha() *Validator {
	return v.addRule(predefined.NewAlphaRule())
}

// Alphanumeric 只包含字母和数字验证
func (v *Validator) Alphanumeric() *Validator {
	return v.addRule(predefined.NewAlphanumericRule())
}

// Numeric 只包含数字字符验证
func (v *Validator) Numeric() *Validator {
	return v.addRule(predefined.NewNumericRule())
}

// UUID UUID格式验证
func (v *Validator) UUID() *Validator {
	return v.addRule(predefined.NewUUIDRule())
}

// IP IP地址验证
func (v *Validator) IP() *Validator {
	return v.addRule(predefined.NewIPRule())
}

// IPv4 IPv4地址验证
func (v *Validator) IPv4() *Validator {
	return v.addRule(predefined.NewIPv4Rule())
}

// IPv6 IPv6地址验证
func (v *Validator) IPv6() *Validator {
	return v.addRule(predefined.NewIPv6Rule())
}

// Phone 手机号验证（中国大陆）
func (v *Validator) Phone() *Validator {
	return v.addRule(predefined.NewPhoneRule())
}

// JSON JSON格式验证
func (v *Validator) JSON() *Validator {
	return v.addRule(predefined.NewJSONRule())
}

// Lowercase 小写字母验证
func (v *Validator) Lowercase() *Validator {
	return v.addRule(predefined.NewLowercaseRule())
}

// Uppercase 大写字母验证
func (v *Validator) Uppercase() *Validator {
	return v.addRule(predefined.NewUppercaseRule())
}

// ========== 比较规则 ==========

// Equals 等于指定值验证
func (v *Validator) Equals(expected any) *Validator {
	return v.addRule(predefined.NewEqualsRule(expected))
}

// NotEquals 不等于指定值验证
func (v *Validator) NotEquals(unexpected any) *Validator {
	return v.addRule(predefined.NewNotEqualsRule(unexpected))
}

// EqualsField 等于另一个字段的值验证（如密码确认）
func (v *Validator) EqualsField(otherField string) *Validator {
	return v.addRule(predefined.NewEqualsFieldRule(otherField))
}

// NotEqualsField 不等于另一个字段的值验证
func (v *Validator) NotEqualsField(otherField string) *Validator {
	return v.addRule(predefined.NewNotEqualsFieldRule(otherField))
}

// GreaterThanField 大于另一个字段的值验证
func (v *Validator) GreaterThanField(otherField string) *Validator {
	return v.addRule(predefined.NewGreaterThanFieldRule(otherField))
}

// LessThanField 小于另一个字段的值验证
func (v *Validator) LessThanField(otherField string) *Validator {
	return v.addRule(predefined.NewLessThanFieldRule(otherField))
}

// RequiredIf 条件必填验证
func (v *Validator) RequiredIf(otherField string, expectedValue any) *Validator {
	return v.addRule(predefined.NewRequiredIfRule(otherField, expectedValue))
}

// RequiredUnless 条件非必填验证
func (v *Validator) RequiredUnless(otherField string, expectedValue any) *Validator {
	return v.addRule(predefined.NewRequiredUnlessRule(otherField, expectedValue))
}

// ========== 执行验证 ==========

// Validate 执行验证（fast-fail 模式）
// 遇到第一个错误立即返回
func (v *Validator) Validate(data any) *core.ValidationResult {
	result := core.NewValidationResult()

	for _, fieldName := range v.fieldOrder {
		rules := v.fieldRules[fieldName]
		fieldValue := extractFieldValue(data, fieldName)
		ctx := core.NewRuleContext(fieldName, fieldValue, data)

		for _, rule := range rules {
			ruleResult := rule.Validate(ctx)
			if !ruleResult.Valid {
				result.AddError(fieldName, rule.Name(), ruleResult.Message)
				return result // fast-fail
			}
		}
	}

	return result
}

// ValidateAll 执行验证（全量收集模式）
// 收集所有字段的所有错误后返回
func (v *Validator) ValidateAll(data any) *core.ValidationResult {
	result := core.NewValidationResult()

	for _, fieldName := range v.fieldOrder {
		rules := v.fieldRules[fieldName]
		fieldValue := extractFieldValue(data, fieldName)
		ctx := core.NewRuleContext(fieldName, fieldValue, data)

		for _, rule := range rules {
			ruleResult := rule.Validate(ctx)
			if !ruleResult.Valid {
				result.AddError(fieldName, rule.Name(), ruleResult.Message)
			}
		}
	}

	return result
}

// Reset 重置验证器状态
func (v *Validator) Reset() *Validator {
	v.fieldRules = make(map[string][]core.Rule)
	v.fieldOrder = make([]string, 0)
	v.currentField = ""
	return v
}

// ========== 全局便捷函数 ==========

// ForField 创建一个新的验证器并指定第一个字段
func ForField(fieldName string) *Validator {
	return New().ForField(fieldName)
}

// ========== 内部辅助函数 ==========

// extractFieldValue 从数据中提取字段值
func extractFieldValue(data any, fieldName string) any {
	if data == nil {
		return nil
	}

	// 如果 data 是 map
	if m, ok := data.(map[string]any); ok {
		return m[fieldName]
	}

	// 如果 data 是结构体或结构体指针
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		// 尝试通过标签查找
		return findFieldByTag(val, fieldName)
	}

	return field.Interface()
}

// findFieldByTag 通过 json 或 form 标签查找字段
func findFieldByTag(val reflect.Value, tagName string) any {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// 检查 json 标签
		if tag := field.Tag.Get("json"); tag == tagName {
			return val.Field(i).Interface()
		}

		// 检查 form 标签
		if tag := field.Tag.Get("form"); tag == tagName {
			return val.Field(i).Interface()
		}
	}
	return nil
}

// invalidPatternRule 无效正则表达式的占位规则
type invalidPatternRule struct {
	pattern string
	err     error
}

func (r *invalidPatternRule) Name() string {
	return "invalid_pattern"
}

func (r *invalidPatternRule) Description() string {
	return "Invalid pattern configuration"
}

func (r *invalidPatternRule) Validate(ctx core.RuleContext) core.RuleResult {
	return core.Failure(ctx.FieldName(), r.Name(), "invalid pattern: "+r.err.Error())
}
