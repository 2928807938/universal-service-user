package provider

import "errors"

var (
	// ErrInvalidConfig 无效的配置
	ErrInvalidConfig = errors.New("oauth provider: 无效的配置")

	// ErrProviderNotFound Provider 不存在
	ErrProviderNotFound = errors.New("oauth provider: provider 不存在")

	// ErrProviderAlreadyRegistered Provider 已注册
	ErrProviderAlreadyRegistered = errors.New("oauth provider: provider 已注册")

	// ErrProviderDisabled Provider 未启用
	ErrProviderDisabled = errors.New("oauth provider: provider 未启用")
)
