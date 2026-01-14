// Package rules 提供一个可扩展的规则引擎，支持内置验证规则和用户自定义规则。
//
// 基本用法：
//
//	// 单个字段验证（fast-fail 模式）
//	result := rules.ForField("email").
//		Required().
//		Email().
//		Validate(user)
//
//	if !result.IsValid() {
//		return result.Error()
//	}
//
//	// 多个字段验证（全量收集错误）
//	result := rules.ForField("username").Required().Length(3, 20).
//		ForField("email").Required().Email().
//		ForField("age").Required().Range(18, 120).
//		ValidateAll(user)
//
// 自定义规则：
//
//	// 注册自定义规则
//	rules.RegisterFunc("custom_id", "Custom ID validation",
//		func(ctx rules.RuleContext) rules.RuleResult {
//			id, ok := ctx.Value().(string)
//			if !ok || !strings.HasPrefix(id, "ID-") {
//				return rules.Failure(ctx.FieldName(), "custom_id", "ID must start with 'ID-'")
//			}
//			return rules.Success()
//		})
//
//	// 使用自定义规则
//	customRule, _ := rules.GetRule("custom_id")
//	result := rules.ForField("userId").
//		Rule(customRule).
//		Validate(data)
package rules

import (
	"github.com/2928807938/universal-service-user/rules/api"
	"github.com/2928807938/universal-service-user/rules/core"
	"github.com/2928807938/universal-service-user/rules/definition"
)

// 重新导出核心类型
type (
	// Rule 是所有规则必须实现的核心接口
	Rule = core.Rule

	// RuleFunc 是函数式规则的类型定义
	RuleFunc = core.RuleFunc

	// RuleContext 封装验证所需的输入数据
	RuleContext = core.RuleContext

	// RuleResult 表示单个规则的验证结果
	RuleResult = core.RuleResult

	// ValidationResult 表示整体验证结果
	ValidationResult = core.ValidationResult

	// ValidationError 表示验证错误
	ValidationError = core.ValidationError

	// Validator 提供链式调用的验证 API
	Validator = api.Validator

	// RuleRegistry 规则注册中心
	RuleRegistry = definition.RuleRegistry
)

// ========== Fluent API ==========

// ForField 创建一个新的验证器并指定第一个字段
func ForField(fieldName string) *Validator {
	return api.ForField(fieldName)
}

// New 创建一个新的验证器
func New() *Validator {
	return api.New()
}

// ========== 自定义规则 API ==========

// Register 注册一个自定义规则
func Register(rule Rule) error {
	return definition.Register(rule)
}

// RegisterFunc 使用函数注册一个自定义规则
func RegisterFunc(name, description string, fn RuleFunc) error {
	return definition.RegisterFunc(name, description, fn)
}

// GetRule 获取已注册的自定义规则
func GetRule(name string) (Rule, error) {
	return definition.Get(name)
}

// HasRule 检查规则是否已注册
func HasRule(name string) bool {
	return definition.Has(name)
}

// NewRegistry 创建一个新的规则注册中心
func NewRegistry() *RuleRegistry {
	return definition.NewRuleRegistry()
}

// ========== 核心辅助函数 ==========

// Success 创建一个成功的验证结果
func Success() RuleResult {
	return core.Success()
}

// Failure 创建一个失败的验证结果
func Failure(field, ruleName, message string) RuleResult {
	return core.Failure(field, ruleName, message)
}

// NewRuleContext 创建一个新的规则上下文
func NewRuleContext(fieldName string, value any, originalValue any) *core.DefaultRuleContext {
	return core.NewRuleContext(fieldName, value, originalValue)
}

// NewCustomRule 创建一个自定义规则
func NewCustomRule(name, description string, fn RuleFunc) Rule {
	return definition.NewCustomRule(name, description, fn)
}

// NewFuncRule 创建一个函数式规则
func NewFuncRule(name, description string, fn RuleFunc) Rule {
	return core.NewFuncRule(name, description, fn)
}
