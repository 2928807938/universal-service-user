package sms

import (
	"fmt"
	"sync"

	"github.com/2928807938/universal-service-user/notification/domain/provider"
)

// Registry 短信提供者注册中心
type Registry struct {
	providers map[string]provider.SMSProvider
	mu        sync.RWMutex
}

// globalRegistry 全局注册中心
var globalRegistry = &Registry{
	providers: make(map[string]provider.SMSProvider),
}

// Register 注册短信提供者
func Register(p provider.SMSProvider) error {
	return globalRegistry.Register(p)
}

// Get 获取短信提供者
func Get(name string) (provider.SMSProvider, error) {
	return globalRegistry.Get(name)
}

// List 列出所有已注册的提供者
func List() []string {
	return globalRegistry.List()
}

// Register 注册短信提供者
func (r *Registry) Register(p provider.SMSProvider) error {
	if p == nil {
		return fmt.Errorf("短信提供者不能为空")
	}

	name := p.GetName()
	if name == "" {
		return fmt.Errorf("短信提供者名称不能为空")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.providers[name]; exists {
		return fmt.Errorf("短信提供者已存在: %s", name)
	}

	r.providers[name] = p
	return nil
}

// Get 获取短信提供者
func (r *Registry) Get(name string) (provider.SMSProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.providers[name]
	if !exists {
		return nil, fmt.Errorf("短信提供者不存在: %s", name)
	}

	return p, nil
}

// List 列出所有已注册的提供者
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}

	return names
}
