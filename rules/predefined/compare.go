package predefined

import (
	"fmt"
	"reflect"

	"github.com/2928807938/universal-service-user/rules/core"
)

// EqualsRule 等于指定值验证规则
type EqualsRule struct {
	expected any
}

func NewEqualsRule(expected any) *EqualsRule {
	return &EqualsRule{expected: expected}
}

func (r *EqualsRule) Name() string {
	return "equals"
}

func (r *EqualsRule) Description() string {
	return fmt.Sprintf("Field must equal: %v", r.expected)
}

func (r *EqualsRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil && r.expected == nil {
		return core.Success()
	}

	if !reflect.DeepEqual(value, r.expected) {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must equal %v", r.expected))
	}

	return core.Success()
}

// NotEqualsRule 不等于指定值验证规则
type NotEqualsRule struct {
	unexpected any
}

func NewNotEqualsRule(unexpected any) *NotEqualsRule {
	return &NotEqualsRule{unexpected: unexpected}
}

func (r *NotEqualsRule) Name() string {
	return "not_equals"
}

func (r *NotEqualsRule) Description() string {
	return fmt.Sprintf("Field must not equal: %v", r.unexpected)
}

func (r *NotEqualsRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()

	if reflect.DeepEqual(value, r.unexpected) {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must not equal %v", r.unexpected))
	}

	return core.Success()
}

// EqualsFieldRule 等于另一个字段的值验证规则（如密码确认）
type EqualsFieldRule struct {
	otherField string
}

func NewEqualsFieldRule(otherField string) *EqualsFieldRule {
	return &EqualsFieldRule{otherField: otherField}
}

func (r *EqualsFieldRule) Name() string {
	return "equals_field"
}

func (r *EqualsFieldRule) Description() string {
	return fmt.Sprintf("Field must equal field: %s", r.otherField)
}

func (r *EqualsFieldRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()

	otherValue, ok := ctx.GetFieldValue(r.otherField)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field %s not found", r.otherField))
	}

	if !reflect.DeepEqual(value, otherValue) {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must equal %s", r.otherField))
	}

	return core.Success()
}

// NotEqualsFieldRule 不等于另一个字段的值验证规则
type NotEqualsFieldRule struct {
	otherField string
}

func NewNotEqualsFieldRule(otherField string) *NotEqualsFieldRule {
	return &NotEqualsFieldRule{otherField: otherField}
}

func (r *NotEqualsFieldRule) Name() string {
	return "not_equals_field"
}

func (r *NotEqualsFieldRule) Description() string {
	return fmt.Sprintf("Field must not equal field: %s", r.otherField)
}

func (r *NotEqualsFieldRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()

	otherValue, ok := ctx.GetFieldValue(r.otherField)
	if !ok {
		// 如果另一个字段不存在，认为不相等
		return core.Success()
	}

	if reflect.DeepEqual(value, otherValue) {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must not equal %s", r.otherField))
	}

	return core.Success()
}

// GreaterThanFieldRule 大于另一个字段的值验证规则
type GreaterThanFieldRule struct {
	otherField string
}

func NewGreaterThanFieldRule(otherField string) *GreaterThanFieldRule {
	return &GreaterThanFieldRule{otherField: otherField}
}

func (r *GreaterThanFieldRule) Name() string {
	return "greater_than_field"
}

func (r *GreaterThanFieldRule) Description() string {
	return fmt.Sprintf("Field must be greater than field: %s", r.otherField)
}

func (r *GreaterThanFieldRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	otherValue, ok := ctx.GetFieldValue(r.otherField)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field %s not found", r.otherField))
	}

	num1, ok1 := toFloat64(value)
	num2, ok2 := toFloat64(otherValue)

	if !ok1 || !ok2 {
		return core.Failure(ctx.FieldName(), r.Name(), "both fields must be numbers")
	}

	if num1 <= num2 {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must be greater than %s", r.otherField))
	}

	return core.Success()
}

// LessThanFieldRule 小于另一个字段的值验证规则
type LessThanFieldRule struct {
	otherField string
}

func NewLessThanFieldRule(otherField string) *LessThanFieldRule {
	return &LessThanFieldRule{otherField: otherField}
}

func (r *LessThanFieldRule) Name() string {
	return "less_than_field"
}

func (r *LessThanFieldRule) Description() string {
	return fmt.Sprintf("Field must be less than field: %s", r.otherField)
}

func (r *LessThanFieldRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	otherValue, ok := ctx.GetFieldValue(r.otherField)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field %s not found", r.otherField))
	}

	num1, ok1 := toFloat64(value)
	num2, ok2 := toFloat64(otherValue)

	if !ok1 || !ok2 {
		return core.Failure(ctx.FieldName(), r.Name(), "both fields must be numbers")
	}

	if num1 >= num2 {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must be less than %s", r.otherField))
	}

	return core.Success()
}

// RequiredIfRule 条件必填验证规则
type RequiredIfRule struct {
	otherField    string
	expectedValue any
}

func NewRequiredIfRule(otherField string, expectedValue any) *RequiredIfRule {
	return &RequiredIfRule{
		otherField:    otherField,
		expectedValue: expectedValue,
	}
}

func (r *RequiredIfRule) Name() string {
	return "required_if"
}

func (r *RequiredIfRule) Description() string {
	return fmt.Sprintf("Field is required if %s equals %v", r.otherField, r.expectedValue)
}

func (r *RequiredIfRule) Validate(ctx core.RuleContext) core.RuleResult {
	otherValue, ok := ctx.GetFieldValue(r.otherField)
	if !ok {
		// 如果条件字段不存在，不需要验证
		return core.Success()
	}

	// 检查条件是否满足
	if !reflect.DeepEqual(otherValue, r.expectedValue) {
		// 条件不满足，不需要验证
		return core.Success()
	}

	// 条件满足，验证当前字段必填
	value := ctx.Value()
	if value == nil {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field is required when %s is %v", r.otherField, r.expectedValue))
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		if val.Len() == 0 {
			return core.Failure(ctx.FieldName(), r.Name(),
				fmt.Sprintf("field is required when %s is %v", r.otherField, r.expectedValue))
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if val.Len() == 0 {
			return core.Failure(ctx.FieldName(), r.Name(),
				fmt.Sprintf("field is required when %s is %v", r.otherField, r.expectedValue))
		}
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			return core.Failure(ctx.FieldName(), r.Name(),
				fmt.Sprintf("field is required when %s is %v", r.otherField, r.expectedValue))
		}
	}

	return core.Success()
}

// RequiredUnlessRule 条件非必填验证规则
type RequiredUnlessRule struct {
	otherField    string
	expectedValue any
}

func NewRequiredUnlessRule(otherField string, expectedValue any) *RequiredUnlessRule {
	return &RequiredUnlessRule{
		otherField:    otherField,
		expectedValue: expectedValue,
	}
}

func (r *RequiredUnlessRule) Name() string {
	return "required_unless"
}

func (r *RequiredUnlessRule) Description() string {
	return fmt.Sprintf("Field is required unless %s equals %v", r.otherField, r.expectedValue)
}

func (r *RequiredUnlessRule) Validate(ctx core.RuleContext) core.RuleResult {
	otherValue, ok := ctx.GetFieldValue(r.otherField)
	if !ok {
		// 如果条件字段不存在，需要验证
	} else if reflect.DeepEqual(otherValue, r.expectedValue) {
		// 条件满足（等于期望值），不需要验证
		return core.Success()
	}

	// 条件不满足，验证当前字段必填
	value := ctx.Value()
	if value == nil {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field is required unless %s is %v", r.otherField, r.expectedValue))
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		if val.Len() == 0 {
			return core.Failure(ctx.FieldName(), r.Name(),
				fmt.Sprintf("field is required unless %s is %v", r.otherField, r.expectedValue))
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if val.Len() == 0 {
			return core.Failure(ctx.FieldName(), r.Name(),
				fmt.Sprintf("field is required unless %s is %v", r.otherField, r.expectedValue))
		}
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			return core.Failure(ctx.FieldName(), r.Name(),
				fmt.Sprintf("field is required unless %s is %v", r.otherField, r.expectedValue))
		}
	}

	return core.Success()
}
