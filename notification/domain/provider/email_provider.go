package provider

import "context"

// EmailMessage 邮件消息
type EmailMessage struct {
	To      []string          // 收件人
	Subject string            // 主题
	Body    string            // 邮件正文(HTML 或纯文本)
	IsHTML  bool              // 是否为 HTML 格式
	Cc      []string          // 抄送
	Bcc     []string          // 密送
	Params  map[string]string // 模板参数(用于模板邮件)
}

// EmailProvider 邮件提供者接口
type EmailProvider interface {
	// GetName 获取提供者唯一标识
	GetName() string

	// GetDisplayName 获取提供者显示名称
	GetDisplayName() string

	// Send 发送邮件
	// ctx: 上下文
	// message: 邮件消息
	// 返回 error 表示发送失败
	Send(ctx context.Context, message *EmailMessage) error

	// SendWithTemplate 使用模板发送邮件
	// ctx: 上下文
	// to: 收件人
	// templateID: 模板 ID
	// params: 模板参数
	// 返回 error 表示发送失败
	SendWithTemplate(ctx context.Context, to []string, templateID string, params map[string]string) error

	// ValidateConfig 验证配置是否完整有效
	// 返回 error 表示配置无效
	ValidateConfig() error
}
