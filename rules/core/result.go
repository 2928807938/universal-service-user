package core

import "strings"

// RuleResult 表示单个规则的验证结果
type RuleResult struct {
	// Valid 表示验证是否通过
	Valid bool

	// Message 错误消息
	Message string

	// Field 验证的字段名
	Field string

	// RuleName 规则名称
	RuleName string
}

// Success 创建一个成功的验证结果
func Success() RuleResult {
	return RuleResult{Valid: true}
}

// Failure 创建一个失败的验证结果
func Failure(field, ruleName, message string) RuleResult {
	return RuleResult{
		Valid:    false,
		Message:  message,
		Field:    field,
		RuleName: ruleName,
	}
}

// ValidationError 表示验证错误
type ValidationError struct {
	Field    string `json:"field"`
	Rule     string `json:"rule"`
	Message  string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// ValidationResult 表示整体验证结果
type ValidationResult struct {
	valid  bool
	errors []*ValidationError
}

// NewValidationResult 创建一个新的验证结果
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		valid:  true,
		errors: make([]*ValidationError, 0),
	}
}

// AddError 添加一个验证错误
func (r *ValidationResult) AddError(field, rule, message string) {
	r.valid = false
	r.errors = append(r.errors, &ValidationError{
		Field:   field,
		Rule:    rule,
		Message: message,
	})
}

// AddRuleResult 添加规则结果
func (r *ValidationResult) AddRuleResult(result RuleResult) {
	if !result.Valid {
		r.AddError(result.Field, result.RuleName, result.Message)
	}
}

// IsValid 返回验证是否通过
func (r *ValidationResult) IsValid() bool {
	return r.valid
}

// Error 返回第一个错误（用于 fast-fail 模式）
func (r *ValidationResult) Error() error {
	if r.valid || len(r.errors) == 0 {
		return nil
	}
	return r.errors[0]
}

// Errors 返回所有错误
func (r *ValidationResult) Errors() []*ValidationError {
	return r.errors
}

// ErrorMap 返回按字段分组的错误映射
func (r *ValidationResult) ErrorMap() map[string][]string {
	result := make(map[string][]string)
	for _, err := range r.errors {
		result[err.Field] = append(result[err.Field], err.Message)
	}
	return result
}

// FirstErrorMap 返回每个字段的第一个错误
func (r *ValidationResult) FirstErrorMap() map[string]string {
	result := make(map[string]string)
	for _, err := range r.errors {
		if _, exists := result[err.Field]; !exists {
			result[err.Field] = err.Message
		}
	}
	return result
}

// ErrorString 返回所有错误的字符串表示
func (r *ValidationResult) ErrorString() string {
	if r.valid {
		return ""
	}
	var msgs []string
	for _, err := range r.errors {
		msgs = append(msgs, err.Field+": "+err.Message)
	}
	return strings.Join(msgs, "; ")
}

// Merge 合并另一个验证结果
func (r *ValidationResult) Merge(other *ValidationResult) {
	if other == nil {
		return
	}
	if !other.valid {
		r.valid = false
		r.errors = append(r.errors, other.errors...)
	}
}
