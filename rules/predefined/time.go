package predefined

import (
	"fmt"
	"time"

	"github.com/2928807938/universal-service-user/rules/core"
)

// toTime 尝试将值转换为 time.Time
func toTime(value any) (time.Time, bool) {
	switch v := value.(type) {
	case time.Time:
		return v, true
	case *time.Time:
		if v != nil {
			return *v, true
		}
		return time.Time{}, false
	case string:
		// 尝试多种常见格式解析
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05",
			"2006-01-02",
			"2006/01/02",
			"02-01-2006",
			"01/02/2006",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, v); err == nil {
				return t, true
			}
		}
		return time.Time{}, false
	case int64:
		// Unix 时间戳（秒）
		return time.Unix(v, 0), true
	default:
		return time.Time{}, false
	}
}

// BeforeRule 早于指定时间验证规则
type BeforeRule struct {
	before time.Time
}

func NewBeforeRule(before time.Time) *BeforeRule {
	return &BeforeRule{before: before}
}

// NewBeforeNowRule 创建一个早于当前时间的规则
func NewBeforeNowRule() *BeforeRule {
	return &BeforeRule{before: time.Now()}
}

func (r *BeforeRule) Name() string {
	return "before"
}

func (r *BeforeRule) Description() string {
	return fmt.Sprintf("Field must be before %s", r.before.Format(time.RFC3339))
}

func (r *BeforeRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	t, ok := toTime(value)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a valid time")
	}

	if !t.Before(r.before) {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must be before %s", r.before.Format(time.RFC3339)))
	}

	return core.Success()
}

// AfterRule 晚于指定时间验证规则
type AfterRule struct {
	after time.Time
}

func NewAfterRule(after time.Time) *AfterRule {
	return &AfterRule{after: after}
}

// NewAfterNowRule 创建一个晚于当前时间的规则
func NewAfterNowRule() *AfterRule {
	return &AfterRule{after: time.Now()}
}

func (r *AfterRule) Name() string {
	return "after"
}

func (r *AfterRule) Description() string {
	return fmt.Sprintf("Field must be after %s", r.after.Format(time.RFC3339))
}

func (r *AfterRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	t, ok := toTime(value)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a valid time")
	}

	if !t.After(r.after) {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must be after %s", r.after.Format(time.RFC3339)))
	}

	return core.Success()
}

// BetweenTimeRule 时间范围验证规则
type BetweenTimeRule struct {
	start time.Time
	end   time.Time
}

func NewBetweenTimeRule(start, end time.Time) *BetweenTimeRule {
	return &BetweenTimeRule{start: start, end: end}
}

func (r *BetweenTimeRule) Name() string {
	return "between_time"
}

func (r *BetweenTimeRule) Description() string {
	return fmt.Sprintf("Field must be between %s and %s",
		r.start.Format(time.RFC3339), r.end.Format(time.RFC3339))
}

func (r *BetweenTimeRule) Validate(ctx core.RuleContext) core.RuleResult {
	value := ctx.Value()
	if value == nil {
		return core.Success()
	}

	t, ok := toTime(value)
	if !ok {
		return core.Failure(ctx.FieldName(), r.Name(), "field must be a valid time")
	}

	if t.Before(r.start) || t.After(r.end) {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must be between %s and %s",
				r.start.Format(time.RFC3339), r.end.Format(time.RFC3339)))
	}

	return core.Success()
}

// DateFormatRule 日期格式验证规则
type DateFormatRule struct {
	format string
}

func NewDateFormatRule(format string) *DateFormatRule {
	return &DateFormatRule{format: format}
}

func (r *DateFormatRule) Name() string {
	return "date_format"
}

func (r *DateFormatRule) Description() string {
	return fmt.Sprintf("Field must match date format: %s", r.format)
}

func (r *DateFormatRule) Validate(ctx core.RuleContext) core.RuleResult {
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

	if _, err := time.Parse(r.format, str); err != nil {
		return core.Failure(ctx.FieldName(), r.Name(),
			fmt.Sprintf("field must match date format: %s", r.format))
	}

	return core.Success()
}
