package service

import (
	"context"
	"fmt"

	"universal-service-user/notification/domain/provider"
	"universal-service-user/verification/domain/enum"
)

// NotificationService 通知领域服务
type NotificationService struct {
	emailProvider provider.EmailProvider
	smsProvider   provider.SMSProvider
}

// NewNotificationService 创建通知服务
func NewNotificationService(emailProvider provider.EmailProvider, smsProvider provider.SMSProvider) *NotificationService {
	return &NotificationService{
		emailProvider: emailProvider,
		smsProvider:   smsProvider,
	}
}

// SendEmailCode 发送邮箱验证码
// ctx: 上下文
// email: 收件人邮箱
// code: 验证码
// scene: 验证码场景
func (s *NotificationService) SendEmailCode(ctx context.Context, email, code string, scene enum.VerificationScene) error {
	if s.emailProvider == nil {
		return fmt.Errorf("邮箱提供者未配置")
	}

	// 构建邮件主题和内容
	subject, body := s.buildEmailContent(code, scene)

	// 发送邮件
	message := &provider.EmailMessage{
		To:      []string{email},
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	}

	if err := s.emailProvider.Send(ctx, message); err != nil {
		return fmt.Errorf("发送邮箱验证码失败: %w", err)
	}

	return nil
}

// SendSMSCode 发送短信验证码
// ctx: 上下文
// phoneNumber: 手机号
// code: 验证码
// scene: 验证码场景
func (s *NotificationService) SendSMSCode(ctx context.Context, phoneNumber, code string, scene enum.VerificationScene) error {
	if s.smsProvider == nil {
		return fmt.Errorf("短信提供者未配置")
	}

	// 获取模板 ID
	templateID := s.getTemplateID(scene)

	// 构建模板参数
	params := map[string]string{
		"code": code,
	}

	// 发送短信
	if err := s.smsProvider.SendWithTemplate(ctx, phoneNumber, templateID, params); err != nil {
		return fmt.Errorf("发送短信验证码失败: %w", err)
	}

	return nil
}

// SendEmail 发送通用邮件
// ctx: 上下文
// to: 收件人列表
// subject: 主题
// body: 邮件正文
// isHTML: 是否为 HTML 格式
func (s *NotificationService) SendEmail(ctx context.Context, to []string, subject, body string, isHTML bool) error {
	if s.emailProvider == nil {
		return fmt.Errorf("邮箱提供者未配置")
	}

	message := &provider.EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  isHTML,
	}

	if err := s.emailProvider.Send(ctx, message); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	return nil
}

// SendSMS 发送通用短信
// ctx: 上下文
// phoneNumber: 手机号
// content: 短信内容
func (s *NotificationService) SendSMS(ctx context.Context, phoneNumber, content string) error {
	if s.smsProvider == nil {
		return fmt.Errorf("短信提供者未配置")
	}

	message := &provider.SMSMessage{
		PhoneNumber: phoneNumber,
		Content:     content,
	}

	if err := s.smsProvider.Send(ctx, message); err != nil {
		return fmt.Errorf("发送短信失败: %w", err)
	}

	return nil
}

// buildEmailContent 构建邮件内容
func (s *NotificationService) buildEmailContent(code string, scene enum.VerificationScene) (subject, body string) {
	sceneDesc := scene.Description()

	switch scene {
	case enum.SceneRegister:
		subject = "注册验证码"
		body = fmt.Sprintf("您的注册验证码是：%s，有效期5分钟。如非本人操作，请忽略此邮件。", code)

	case enum.SceneLogin:
		subject = "登录验证码"
		body = fmt.Sprintf("您的登录验证码是：%s，有效期5分钟。如非本人操作，请忽略此邮件。", code)

	case enum.SceneResetPassword:
		subject = "重置密码验证码"
		body = fmt.Sprintf("您的重置密码验证码是：%s，有效期5分钟。如非本人操作，请立即修改密码。", code)

	case enum.SceneChangeEmailOld:
		subject = "更换邮箱验证"
		body = fmt.Sprintf("您正在更换邮箱，验证码是：%s，有效期5分钟。如非本人操作，请立即修改密码。", code)

	case enum.SceneChangeEmailNew:
		subject = "绑定新邮箱验证"
		body = fmt.Sprintf("您正在绑定新邮箱，验证码是：%s，有效期5分钟。如非本人操作，请忽略此邮件。", code)

	case enum.SceneChangePhoneOld:
		subject = "更换手机号验证"
		body = fmt.Sprintf("您正在更换手机号，验证码是：%s，有效期5分钟。如非本人操作，请立即修改密码。", code)

	case enum.SceneChangePhoneNew:
		subject = "绑定新手机号验证"
		body = fmt.Sprintf("您正在绑定新手机号，验证码是：%s，有效期5分钟。如非本人操作，请忽略此短信。", code)

	default:
		subject = fmt.Sprintf("%s验证码", sceneDesc)
		body = fmt.Sprintf("您的验证码是：%s，有效期5分钟。", code)
	}

	return subject, body
}

// getTemplateID 根据场景获取模板 ID
func (s *NotificationService) getTemplateID(scene enum.VerificationScene) string {
	// 这里返回场景的字符串表示作为模板 ID
	// 实际使用时，应该从配置中读取对应的模板 ID
	return scene.String()
}
