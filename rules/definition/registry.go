package definition

import (
	"sync"

	"github.com/2928807938/universal-service-user/rules/core"
)

// RuleRegistry 是规则注册中心
// 管理用户自定义规则类型，支持规则的注册和查找
type RuleRegistry struct {
	mu    sync.RWMutex
	rules map[string]core.Rule
}

// NewRuleRegistry 创建一个新的规则注册中心
func NewRuleRegistry() *RuleRegistry {
	return &RuleRegistry{
		rules: make(map[string]core.Rule),
	}
}

// Register 注册一个规则
// 如果规则已存在，返回 DuplicateRuleError
func (r *RuleRegistry) Register(rule core.Rule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := rule.Name()
	if _, exists := r.rules[name]; exists {
		return core.NewDuplicateRuleError(name)
	}

	r.rules[name] = rule
	return nil
}

// RegisterWithName 使用指定名称注册一个规则
// 如果规则已存在，返回 DuplicateRuleError
func (r *RuleRegistry) RegisterWithName(name string, rule core.Rule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.rules[name]; exists {
		return core.NewDuplicateRuleError(name)
	}

	r.rules[name] = rule
	return nil
}

// RegisterFunc 使用函数注册一个规则
func (r *RuleRegistry) RegisterFunc(name, description string, fn core.RuleFunc) error {
	rule := NewCustomRule(name, description, fn)
	return r.RegisterWithName(name, rule)
}

// MustRegister 注册一个规则，如果失败则 panic
func (r *RuleRegistry) MustRegister(rule core.Rule) {
	if err := r.Register(rule); err != nil {
		panic(err)
	}
}

// MustRegisterFunc 使用函数注册一个规则，如果失败则 panic
func (r *RuleRegistry) MustRegisterFunc(name, description string, fn core.RuleFunc) {
	if err := r.RegisterFunc(name, description, fn); err != nil {
		panic(err)
	}
}

// Get 获取一个规则
// 如果规则不存在，返回 nil 和 RuleNotFoundError
func (r *RuleRegistry) Get(name string) (core.Rule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rule, exists := r.rules[name]
	if !exists {
		return nil, core.NewRuleNotFoundError(name)
	}

	return rule, nil
}

// Has 检查规则是否存在
func (r *RuleRegistry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.rules[name]
	return exists
}

// Unregister 注销一个规则
func (r *RuleRegistry) Unregister(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.rules[name]; exists {
		delete(r.rules, name)
		return true
	}
	return false
}

// List 列出所有已注册的规则名称
func (r *RuleRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.rules))
	for name := range r.rules {
		names = append(names, name)
	}
	return names
}

// Count 返回已注册规则的数量
func (r *RuleRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.rules)
}

// Clear 清除所有已注册的规则
func (r *RuleRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.rules = make(map[string]core.Rule)
}

// DefaultRegistry 是默认的全局规则注册中心
var DefaultRegistry = NewRuleRegistry()

// Register 在默认注册中心注册一个规则
func Register(rule core.Rule) error {
	return DefaultRegistry.Register(rule)
}

// RegisterFunc 在默认注册中心使用函数注册一个规则
func RegisterFunc(name, description string, fn core.RuleFunc) error {
	return DefaultRegistry.RegisterFunc(name, description, fn)
}

// Get 从默认注册中心获取一个规则
func Get(name string) (core.Rule, error) {
	return DefaultRegistry.Get(name)
}

// Has 检查默认注册中心中规则是否存在
func Has(name string) bool {
	return DefaultRegistry.Has(name)
}
