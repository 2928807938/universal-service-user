package hook

import (
	"sync"
)

// Registry 钩子注册中心
// 用于注册和管理各种类型的钩子函数
// 线程安全
type Registry struct {
	// hooks 存储所有注册的钩子
	// key: HookType, value: 该类型的所有钩子函数列表
	hooks map[HookType][]HookFunc

	// mu 读写锁，保证并发安全
	mu sync.RWMutex
}

// NewRegistry 创建新的钩子注册中心
func NewRegistry() *Registry {
	return &Registry{
		hooks: make(map[HookType][]HookFunc),
	}
}

// Register 注册钩子
// hookType: 钩子类型
// hookFunc: 钩子函数
// 说明：同一个钩子类型可以注册多个钩子函数，执行时按注册顺序依次执行
func (r *Registry) Register(hookType HookType, hookFunc HookFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.hooks[hookType] = append(r.hooks[hookType], hookFunc)
}

// Get 获取指定类型的所有钩子函数
// hookType: 钩子类型
// 返回该类型的所有钩子函数列表，如果没有注册则返回空切片
func (r *Registry) Get(hookType HookType) []HookFunc {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 返回副本，避免外部修改
	hooks, ok := r.hooks[hookType]
	if !ok {
		return []HookFunc{}
	}

	result := make([]HookFunc, len(hooks))
	copy(result, hooks)
	return result
}

// Has 检查是否注册了指定类型的钩子
// hookType: 钩子类型
// 返回 true 表示已注册该类型的钩子，false 表示未注册
func (r *Registry) Has(hookType HookType) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	hooks, ok := r.hooks[hookType]
	return ok && len(hooks) > 0
}

// Count 获取指定类型已注册的钩子数量
// hookType: 钩子类型
// 返回该类型已注册的钩子数量
func (r *Registry) Count(hookType HookType) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	hooks, ok := r.hooks[hookType]
	if !ok {
		return 0
	}
	return len(hooks)
}

// Clear 清空指定类型的所有钩子
// hookType: 钩子类型
func (r *Registry) Clear(hookType HookType) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.hooks, hookType)
}

// ClearAll 清空所有钩子
func (r *Registry) ClearAll() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.hooks = make(map[HookType][]HookFunc)
}

// GetAllTypes 获取所有已注册钩子的类型
// 返回所有已注册钩子的类型列表
func (r *Registry) GetAllTypes() []HookType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]HookType, 0, len(r.hooks))
	for hookType := range r.hooks {
		types = append(types, hookType)
	}
	return types
}

// 全局钩子注册中心
var defaultRegistry = NewRegistry()

// Register 在全局注册中心注册钩子
func Register(hookType HookType, hookFunc HookFunc) {
	defaultRegistry.Register(hookType, hookFunc)
}

// Get 从全局注册中心获取钩子
func Get(hookType HookType) []HookFunc {
	return defaultRegistry.Get(hookType)
}

// Has 检查全局注册中心是否有指定类型的钩子
func Has(hookType HookType) bool {
	return defaultRegistry.Has(hookType)
}

// Count 获取全局注册中心中指定类型的钩子数量
func Count(hookType HookType) int {
	return defaultRegistry.Count(hookType)
}

// Clear 清空全局注册中心中指定类型的钩子
func Clear(hookType HookType) {
	defaultRegistry.Clear(hookType)
}

// ClearAll 清空全局注册中心中的所有钩子
func ClearAll() {
	defaultRegistry.ClearAll()
}

// GetDefaultRegistry 获取全局钩子注册中心
// 用于在需要独立注册中心的场景下获取全局实例
func GetDefaultRegistry() *Registry {
	return defaultRegistry
}
