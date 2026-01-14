package core

// Rule 是所有规则必须实现的核心接口
type Rule interface {
	// Validate 执行验证逻辑
	Validate(ctx RuleContext) RuleResult

	// Name 返回规则名称
	Name() string

	// Description 返回规则描述
	Description() string
}

// RuleFunc 是函数式规则的类型定义
// 允许使用简单函数作为规则
type RuleFunc func(ctx RuleContext) RuleResult

// funcRule 是 RuleFunc 的包装器，实现 Rule 接口
type funcRule struct {
	name        string
	description string
	fn          RuleFunc
}

// NewFuncRule 创建一个函数式规则
func NewFuncRule(name, description string, fn RuleFunc) Rule {
	return &funcRule{
		name:        name,
		description: description,
		fn:          fn,
	}
}

func (r *funcRule) Validate(ctx RuleContext) RuleResult {
	return r.fn(ctx)
}

func (r *funcRule) Name() string {
	return r.name
}

func (r *funcRule) Description() string {
	return r.description
}
