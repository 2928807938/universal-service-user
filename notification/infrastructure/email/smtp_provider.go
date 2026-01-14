package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"universal-service-user/notification/domain/provider"
)

// SMTPConfig SMTP 配置
type SMTPConfig struct {
	Host      string            // SMTP 服务器地址
	Port      int               // SMTP 端口
	Username  string            // 用户名
	Password  string            // 密码
	From      string            // 发件人(格式: "系统通知 <noreply@example.com>")
	Templates map[string]string // 邮件模板(key: 模板ID, value: 模板内容)
}

// SMTPProvider SMTP 邮件提供者
type SMTPProvider struct {
	config *SMTPConfig
	auth   smtp.Auth
}

// NewSMTPProvider 创建 SMTP 邮件提供者
func NewSMTPProvider(config *SMTPConfig) (*SMTPProvider, error) {
	if err := validateSMTPConfig(config); err != nil {
		return nil, err
	}

	// 创建 SMTP 认证
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	return &SMTPProvider{
		config: config,
		auth:   auth,
	}, nil
}

// GetName 获取提供者唯一标识
func (p *SMTPProvider) GetName() string {
	return "smtp"
}

// GetDisplayName 获取提供者显示名称
func (p *SMTPProvider) GetDisplayName() string {
	return "SMTP 邮件"
}

// Send 发送邮件
func (p *SMTPProvider) Send(ctx context.Context, message *provider.EmailMessage) error {
	// 构建邮件内容
	msg := p.buildMessage(message)

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", p.config.Host, p.config.Port)
	recipients := append(message.To, message.Cc...)
	recipients = append(recipients, message.Bcc...)

	if err := smtp.SendMail(addr, p.auth, p.extractEmail(p.config.From), recipients, []byte(msg)); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	return nil
}

// SendWithTemplate 使用模板发送邮件
func (p *SMTPProvider) SendWithTemplate(ctx context.Context, to []string, templateID string, params map[string]string) error {
	// 获取模板内容
	templateContent, exists := p.config.Templates[templateID]
	if !exists {
		return fmt.Errorf("模板不存在: %s", templateID)
	}

	// 解析模板
	tmpl, err := template.New(templateID).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("解析模板失败: %w", err)
	}

	// 渲染模板
	var body bytes.Buffer
	if err := tmpl.Execute(&body, params); err != nil {
		return fmt.Errorf("渲染模板失败: %w", err)
	}

	// 发送邮件
	message := &provider.EmailMessage{
		To:      to,
		Subject: getSubjectFromParams(params, templateID),
		Body:    body.String(),
		IsHTML:  false,
	}

	return p.Send(ctx, message)
}

// ValidateConfig 验证配置是否完整有效
func (p *SMTPProvider) ValidateConfig() error {
	return validateSMTPConfig(p.config)
}

// buildMessage 构建邮件内容
func (p *SMTPProvider) buildMessage(message *provider.EmailMessage) string {
	var msg strings.Builder

	// From
	msg.WriteString(fmt.Sprintf("From: %s\r\n", p.config.From))

	// To
	if len(message.To) > 0 {
		msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(message.To, ", ")))
	}

	// Cc
	if len(message.Cc) > 0 {
		msg.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(message.Cc, ", ")))
	}

	// Subject
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", message.Subject))

	// MIME
	msg.WriteString("MIME-Version: 1.0\r\n")
	if message.IsHTML {
		msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}

	// Body
	msg.WriteString("\r\n")
	msg.WriteString(message.Body)

	return msg.String()
}

// extractEmail 从 "Name <email>" 格式中提取邮箱地址
func (p *SMTPProvider) extractEmail(from string) string {
	start := strings.Index(from, "<")
	end := strings.Index(from, ">")

	if start != -1 && end != -1 && start < end {
		return from[start+1 : end]
	}

	return from
}

// validateSMTPConfig 验证 SMTP 配置
func validateSMTPConfig(config *SMTPConfig) error {
	if config == nil {
		return fmt.Errorf("SMTP 配置不能为空")
	}
	if config.Host == "" {
		return fmt.Errorf("SMTP Host 不能为空")
	}
	if config.Port <= 0 {
		return fmt.Errorf("SMTP Port 必须大于 0")
	}
	if config.Username == "" {
		return fmt.Errorf("SMTP Username 不能为空")
	}
	if config.Password == "" {
		return fmt.Errorf("SMTP Password 不能为空")
	}
	if config.From == "" {
		return fmt.Errorf("SMTP From 不能为空")
	}

	return nil
}

// getSubjectFromParams 从参数中获取主题,如果没有则使用默认值
func getSubjectFromParams(params map[string]string, defaultSubject string) string {
	if subject, exists := params["subject"]; exists {
		return subject
	}
	return defaultSubject
}
