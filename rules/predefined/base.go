package predefined

import (
	"reflect"

	"github.com/2928807938/universal-service-user/rules/core"
)

// RequiredRule 必填验证规则
type RequiredRule struct{}

func NewRequiredRule() *RequiredRule {
	return &RequiredRule{}
}

func (r *RequiredRule) Name() string {
	return "required"
}

func (r *RequiredRule) Description() string {
	return "Field is required"
}

func (r *RequiredRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()

	if value == nil {
		return core.Failure(ctx.FieldName(), r.Name(), "field is required")
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		if val.Len() == 0 {
			return core.Failure(ctx.FieldName(), r.Name(), "field is required")
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if val.Len() == 0 {
			return core.Failure(ctx.FieldName(), r.Name(), "field is required")
		}
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			return core.Failure(ctx.FieldName(), r.Name(), "field is required")
		}
	}

	return core.Success()
}

// NotNullRule 非空验证规则
type NotNullRule struct{}

func NewNotNullRule() *NotNullRule {
	return &NotNullRule{}
}

func (r *NotNullRule) Name() string {
	return "not_null"
}

func (r *NotNullRule) Description() string {
	return "Field must not be null"
}

func (r *NotNullRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()

	if value == nil {
		return core.Failure(ctx.FieldName(), r.Name(), "field must not be null")
	}

	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		if val.IsNil() {
			return core.Failure(ctx.FieldName(), r.Name(), "field must not be null")
		}
	}

	return core.Success()
}

// NotEmptyRule 非空字符串验证规则
type NotEmptyRule struct{}

func NewNotEmptyRule() *NotEmptyRule {
	return &NotEmptyRule{}
}

func (r *NotEmptyRule) Name() string {
	return "not_empty"
}

func (r *NotEmptyRule) Description() string {
	return "Field must not be empty"
}

func (r *NotEmptyRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()

	if value == nil {
		return core.Failure(ctx.FieldName(), r.Name(), "field must not be empty")
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		if val.Len() == 0 {
			return core.Failure(ctx.FieldName(), r.Name(), "field must not be empty")
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if val.Len() == 0 {
			return core.Failure(ctx.FieldName(), r.Name(), "field must not be empty")
		}
	}

	return core.Success()
}
