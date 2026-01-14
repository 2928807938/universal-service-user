package predefined

import (
	"fmt"
	"reflect"

	"github.com/2928807938/universal-service-user/rules/core"
)

// InRule 值在指定集合中验证规则
type InRule struct {
	values []any
}

func NewInRule(values ...any) *InRule {
	return &InRule{values: values}
}

func (r *InRule) Name() string {
	return "in"
}

func (r *InRule) Description() string {
	return fmt.Sprintf("Field must be one of: %v", r.values)
}

func (r *InRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	for _, v := range r.values {
		if reflect.DeepEqual(value, v) {
			return core.Success()
		}
	}

	return core.Failure(ctx.FieldName(), r.Name(),
		fmt.Sprintf("field must be one of: %v", r.values))
}

// NotInRule 值不在指定集合中验证规则
type NotInRule struct {
	values []any
}

func NewNotInRule(values ...any) *NotInRule {
	return &NotInRule{values: values}
}

func (r *NotInRule) Name() string {
	return "not_in"
}

func (r *NotInRule) Description() string {
	return fmt.Sprintf("Field must not be one of: %v", r.values)
}

func (r *NotInRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	for _, v := range r.values {
		if reflect.DeepEqual(value, v) {
			return core.Failure(ctx.FieldName(), r.Name(),
				fmt.Sprintf("field must not be one of: %v", r.values))
		}
	}

	return core.Success()
}

// ArrayLengthRule 数组/切片长度验证规则
type ArrayLengthRule struct {
	min int
	max int
}

func NewArrayLengthRule(min, max int) *ArrayLengthRule {
	return &ArrayLengthRule{min: min, max: max}
}

func (r *ArrayLengthRule) Name() string {
	return "array_length"
}

func (r *ArrayLengthRule) Description() string {
	return fmt.Sprintf("Array length must be between %d and %d", r.min, r.max)
}

func (r *ArrayLengthRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be an array or slice")
	}

	length := val.Len()
	if length < r.min || length > r.max {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("array length must be between %d and %d, got %d", r.min, r.max, length))
	}

	return core.Success()
}

// MinArrayLengthRule 数组最小长度验证规则
type MinArrayLengthRule struct {
	min int
}

func NewMinArrayLengthRule(min int) *MinArrayLengthRule {
	return &MinArrayLengthRule{min: min}
}

func (r *MinArrayLengthRule) Name() string {
	return "min_array_length"
}

func (r *MinArrayLengthRule) Description() string {
	return fmt.Sprintf("Array length must be at least %d", r.min)
}

func (r *MinArrayLengthRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be an array or slice")
	}

	length := val.Len()
	if length < r.min {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("array length must be at least %d, got %d", r.min, length))
	}

	return core.Success()
}

// MaxArrayLengthRule 数组最大长度验证规则
type MaxArrayLengthRule struct {
	max int
}

func NewMaxArrayLengthRule(max int) *MaxArrayLengthRule {
	return &MaxArrayLengthRule{max: max}
}

func (r *MaxArrayLengthRule) Name() string {
	return "max_array_length"
}

func (r *MaxArrayLengthRule) Description() string {
	return fmt.Sprintf("Array length must be at most %d", r.max)
}

func (r *MaxArrayLengthRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be an array or slice")
	}

	length := val.Len()
	if length > r.max {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("array length must be at most %d, got %d", r.max, length))
	}

	return core.Success()
}

// UniqueRule 数组元素唯一性验证规则
type UniqueRule struct{}

func NewUniqueRule() *UniqueRule {
	return &UniqueRule{}
}

func (r *UniqueRule) Name() string {
	return "unique"
}

func (r *UniqueRule) Description() string {
	return "Array elements must be unique"
}

func (r *UniqueRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be an array or slice")
	}

	seen := make(map[any]bool)
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i).Interface()
		if seen[elem] {
			return core.Failure(ctx.FieldName(), r.Name(),
				fmt.Sprintf("array elements must be unique, duplicate found: %v", elem))
		}
		seen[elem] = true
	}

	return core.Success()
}
