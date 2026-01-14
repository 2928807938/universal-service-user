package sms

import (
	"fmt"

	"universal-service-user/notification/domain/provider"
)

// ProviderFactory 短信提供者工厂
type ProviderFactory struct {
	registry *Registry
}

// NewProviderFactory 创建短信提供者工厂
func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		registry: globalRegistry,
	}
}

// Create 创建短信提供者
// name: 提供者名称(如 "tencent", "aliyun")
// 返回对应的提供者实例,如果不存在则返回错误
func (f *ProviderFactory) Create(name string) (provider.SMSProvider, error) {
	return f.registry.Get(name)
}

// CreateWithValidation 创建并验证短信提供者
// name: 提供者名称
// 返回验证通过的提供者实例
func (f *ProviderFactory) CreateWithValidation(name string) (provider.SMSProvider, error) {
	// 创建提供者
	p, err := f.Create(name)
	if err != nil {
		return nil, err
	}

	// 验证配置
	if err := p.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("提供者配置验证失败: %w", err)
	}

	return p, nil
}

// ListAvailable 列出所有可用的提供者
func (f *ProviderFactory) ListAvailable() []string {
	return f.registry.List()
}
