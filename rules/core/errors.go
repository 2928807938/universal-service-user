package core

import "fmt"

// RuleConfigError 表示规则配置错误
type RuleConfigError struct {
	RuleName string
	Message  string
}

func (e *RuleConfigError) Error() string {
	return fmt.Sprintf("rule config error [%s]: %s", e.RuleName, e.Message)
}

// NewRuleConfigError 创建一个规则配置错误
func NewRuleConfigError(ruleName, message string) *RuleConfigError {
	return &RuleConfigError{
		RuleName: ruleName,
		Message:  message,
	}
}

// RuleExecutionError 表示规则执行错误
type RuleExecutionError struct {
	RuleName string
	Field    string
	Cause    error
}

func (e *RuleExecutionError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("rule execution error [%s] on field [%s]: %s", e.RuleName, e.Field, e.Cause.Error())
	}
	return fmt.Sprintf("rule execution error [%s] on field [%s]", e.RuleName, e.Field)
}

func (e *RuleExecutionError) Unwrap() error {
	return e.Cause
}

// NewRuleExecutionError 创建一个规则执行错误
func NewRuleExecutionError(ruleName, field string, cause error) *RuleExecutionError {
	return &RuleExecutionError{
		RuleName: ruleName,
		Field:    field,
		Cause:    cause,
	}
}

// RuleNotFoundError 表示规则未找到错误
type RuleNotFoundError struct {
	RuleName string
}

func (e *RuleNotFoundError) Error() string {
	return fmt.Sprintf("rule not found: %s", e.RuleName)
}

// NewRuleNotFoundError 创建一个规则未找到错误
func NewRuleNotFoundError(ruleName string) *RuleNotFoundError {
	return &RuleNotFoundError{
		RuleName: ruleName,
	}
}

// DuplicateRuleError 表示规则重复注册错误
type DuplicateRuleError struct {
	RuleName string
}

func (e *DuplicateRuleError) Error() string {
	return fmt.Sprintf("rule already registered: %s", e.RuleName)
}

// NewDuplicateRuleError 创建一个规则重复注册错误
func NewDuplicateRuleError(ruleName string) *DuplicateRuleError {
	return &DuplicateRuleError{
		RuleName: ruleName,
	}
}
