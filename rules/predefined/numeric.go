package predefined

import (
	"fmt"

	"universal-service-user/rules/core"
)

// toFloat64 尝试将值转换为 float64
func toFloat64(value any) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

// MinRule 最小值验证规则
type MinRule struct {
	min float64
}

func NewMinRule(min float64) *MinRule {
	return &MinRule{min: min}
}

func (r *MinRule) Name() string {
	return "min"
}

func (r *MinRule) Description() string {
	return fmt.Sprintf("Field must be at least %v", r.min)
}

func (r *MinRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	num, ok := toFloat64(value)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a number")
	}

	if num < r.min {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must be at least %v, got %v", r.min, num))
	}

	return core.Success()
}

// MaxRule 最大值验证规则
type MaxRule struct {
	max float64
}

func NewMaxRule(max float64) *MaxRule {
	return &MaxRule{max: max}
}

func (r *MaxRule) Name() string {
	return "max"
}

func (r *MaxRule) Description() string {
	return fmt.Sprintf("Field must be at most %v", r.max)
}

func (r *MaxRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	num, ok := toFloat64(value)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a number")
	}

	if num > r.max {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must be at most %v, got %v", r.max, num))
	}

	return core.Success()
}

// RangeRule 数值范围验证规则
type RangeRule struct {
	min float64
	max float64
}

func NewRangeRule(min, max float64) *RangeRule {
	return &RangeRule{min: min, max: max}
}

func (r *RangeRule) Name() string {
	return "range"
}

func (r *RangeRule) Description() string {
	return fmt.Sprintf("Field must be between %v and %v", r.min, r.max)
}

func (r *RangeRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	num, ok := toFloat64(value)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a number")
	}

	if num < r.min || num > r.max {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must be between %v and %v, got %v", r.min, r.max, num))
	}

	return core.Success()
}

// PositiveRule 正数验证规则
type PositiveRule struct{}

func NewPositiveRule() *PositiveRule {
	return &PositiveRule{}
}

func (r *PositiveRule) Name() string {
	return "positive"
}

func (r *PositiveRule) Description() string {
	return "Field must be a positive number"
}

func (r *PositiveRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	num, ok := toFloat64(value)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a number")
	}

	if num <= 0 {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must be a positive number, got %v", num))
	}

	return core.Success()
}

// NegativeRule 负数验证规则
type NegativeRule struct{}

func NewNegativeRule() *NegativeRule {
	return &NegativeRule{}
}

func (r *NegativeRule) Name() string {
	return "negative"
}

func (r *NegativeRule) Description() string {
	return "Field must be a negative number"
}

func (r *NegativeRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	num, ok := toFloat64(value)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a number")
	}

	if num >= 0 {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must be a negative number, got %v", num))
	}

	return core.Success()
}

// NonNegativeRule 非负数验证规则
type NonNegativeRule struct{}

func NewNonNegativeRule() *NonNegativeRule {
	return &NonNegativeRule{}
}

func (r *NonNegativeRule) Name() string {
	return "non_negative"
}

func (r *NonNegativeRule) Description() string {
	return "Field must be a non-negative number"
}

func (r *NonNegativeRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	num, ok := toFloat64(value)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a number")
	}

	if num < 0 {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must be a non-negative number, got %v", num))
	}

	return core.Success()
}
