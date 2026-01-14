package provider

import (
	"fmt"
	"sync"
)

// Manager OAuth Provider 管理器
type Manager struct {
	providers map[string]Provider
	mu        sync.RWMutex
}

// NewManager 创建 OAuth Provider 管理器
func NewManager() *Manager {
	return &Manager{
		providers: make(map[string]Provider),
	}
}

// globalManager 全局管理器实例
var globalManager = NewManager()

// GetManager 获取全局管理器
func GetManager() *Manager {
	return globalManager
}

// RegisterProvider 注册 Provider
func (m *Manager) RegisterProvider(p Provider) error {
	if p == nil {
		return fmt.Errorf("provider 不能为空")
	}

	name := p.GetName()
	if name == "" {
		return fmt.Errorf("provider 名称不能为空")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查是否已注册
	if _, exists := m.providers[name]; exists {
		return fmt.Errorf("provider 已注册: %s", name)
	}

	// 验证配置
	if err := p.ValidateConfig(); err != nil {
		return fmt.Errorf("provider 配置无效: %w", err)
	}

	m.providers[name] = p
	return nil
}

// GetProvider 获取 Provider
func (m *Manager) GetProvider(name string) (Provider, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, exists := m.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider 不存在: %s", name)
	}

	return p, nil
}

// HasProvider 判断 Provider 是否存在
func (m *Manager) HasProvider(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.providers[name]
	return exists
}

// ListProviders 列出所有已注册的 Providers
func (m *Manager) ListProviders() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.providers))
	for name := range m.providers {
		names = append(names, name)
	}
	return names
}

// UnregisterProvider 注销 Provider
func (m *Manager) UnregisterProvider(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.providers[name]; !exists {
		return fmt.Errorf("provider 不存在: %s", name)
	}

	delete(m.providers, name)
	return nil
}

// Clear 清空所有 Providers
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.providers = make(map[string]Provider)
}

// GetEnabledProviders 获取所有已启用的 Providers
func (m *Manager) GetEnabledProviders() []Provider {
	m.mu.RLock()
	defer m.mu.RUnlock()

	enabled := make([]Provider, 0, len(m.providers))
	for _, p := range m.providers {
		enabled = append(enabled, p)
	}
	return enabled
}
