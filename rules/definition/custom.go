package definition

import (
	"github.com/2928807938/universal-service-user/rules/core"
)

// CustomRule 是通用的自定义规则实现
// 接受验证函数作为参数，支持闭包捕获上下文
type CustomRule struct {
	name        string
	description string
	validateFn  core.RuleFunc
}

// NewCustomRule 创建一个自定义规则
func NewCustomRule(name, description string, validateFn core.RuleFunc) *CustomRule {
	return &CustomRule{
		name:        name,
		description: description,
		validateFn:  validateFn,
	}
}

func (r *CustomRule) Name() string {
	return r.name
}

func (r *CustomRule) Description() string {
	return r.description
}

func (r *CustomRule) Validate(ctx core.RuleContext) core.RuleResult {
	if r.validateFn == nil {
		return core.Success()
	}
	return r.validateFn(ctx)
}

// CustomRuleBuilder 提供链式构建自定义规则的能力
type CustomRuleBuilder struct {
	name        string
	description string
	validateFn  core.RuleFunc
}

// NewCustomRuleBuilder 创建一个自定义规则构建器
func NewCustomRuleBuilder(name string) *CustomRuleBuilder {
	return &CustomRuleBuilder{
		name: name,
	}
}

// Description 设置规则描述
func (b *CustomRuleBuilder) Description(desc string) *CustomRuleBuilder {
	b.description = desc
	return b
}

// ValidateWith 设置验证函数
func (b *CustomRuleBuilder) ValidateWith(fn core.RuleFunc) *CustomRuleBuilder {
	b.validateFn = fn
	return b
}

// Build 构建自定义规则
func (b *CustomRuleBuilder) Build() *CustomRule {
	return NewCustomRule(b.name, b.description, b.validateFn)
}
