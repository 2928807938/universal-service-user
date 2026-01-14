package provider

import "context"

// SMSMessage 短信消息
type SMSMessage struct {
	PhoneNumber string            // 手机号
	Content     string            // 短信内容(用于纯文本短信)
	TemplateID  string            // 模板 ID
	Params      map[string]string // 模板参数
	SignName    string            // 签名(某些厂商需要)
}

// SMSProvider 短信提供者接口
type SMSProvider interface {
	// GetName 获取提供者唯一标识(如 "tencent", "aliyun")
	GetName() string

	// GetDisplayName 获取提供者显示名称(如 "腾讯云短信", "阿里云短信")
	GetDisplayName() string

	// Send 发送短信
	// ctx: 上下文
	// message: 短信消息
	// 返回 error 表示发送失败
	Send(ctx context.Context, message *SMSMessage) error

	// SendWithTemplate 使用模板发送短信
	// ctx: 上下文
	// phoneNumber: 手机号
	// templateID: 模板 ID
	// params: 模板参数
	// 返回 error 表示发送失败
	SendWithTemplate(ctx context.Context, phoneNumber, templateID string, params map[string]string) error

	// ValidateConfig 验证配置是否完整有效
	// 返回 error 表示配置无效
	ValidateConfig() error
}
