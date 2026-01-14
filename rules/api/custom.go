package api

import (
	"universal-service-user/rules/core"
	"universal-service-user/rules/definition"
)

// CustomRuleBuilder 自定义规则构建器
// 提供链式 API 来创建自定义规则
type CustomRuleBuilder struct {
	name        string
	description string
	validateFn  core.RuleFunc
}

// NewCustomRule 创建一个自定义规则构建器
func NewCustomRule(name string) *CustomRuleBuilder {
	return &CustomRuleBuilder{
		name: name,
	}
}

// Description 设置规则描述
func (b *CustomRuleBuilder) Description(desc string) *CustomRuleBuilder {
	b.description = desc
	return b
}

// Validate 设置验证函数
func (b *CustomRuleBuilder) Validate(fn core.RuleFunc) *CustomRuleBuilder {
	b.validateFn = fn
	return b
}

// Build 构建规则
func (b *CustomRuleBuilder) Build() core.Rule {
	return definition.NewCustomRule(b.name, b.description, b.validateFn)
}

// Register 构建并注册到默认注册中心
func (b *CustomRuleBuilder) Register() error {
	rule := b.Build()
	return definition.Register(rule)
}

// MustRegister 构建并注册，失败则 panic
func (b *CustomRuleBuilder) MustRegister() {
	if err := b.Register(); err != nil {
		panic(err)
	}
}

// ========== 规则注册 API ==========

// Register 注册一个规则到默认注册中心
func Register(rule core.Rule) error {
	return definition.Register(rule)
}

// RegisterFunc 使用函数注册规则到默认注册中心
func RegisterFunc(name, description string, fn core.RuleFunc) error {
	return definition.RegisterFunc(name, description, fn)
}

// GetRule 从默认注册中心获取规则
func GetRule(name string) (core.Rule, error) {
	return definition.Get(name)
}

// HasRule 检查规则是否存在
func HasRule(name string) bool {
	return definition.Has(name)
}

// ========== 辅助函数 ==========

// Success 创建成功的验证结果
func Success() core.RuleResult {
	return core.Success()
}

// Failure 创建失败的验证结果
func Failure(field, ruleName, message string) core.RuleResult {
	return core.Failure(field, ruleName, message)
}
