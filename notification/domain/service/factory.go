package service

import (
	"fmt"

	"github.com/2928807938/universal-service-user/notification/domain/provider"
	"github.com/2928807938/universal-service-user/notification/infrastructure/email"
	"github.com/2928807938/universal-service-user/notification/infrastructure/sms/aliyun"
	"github.com/2928807938/universal-service-user/notification/infrastructure/sms/tencent"
)

// NewNotificationServiceFromConfig 从配置创建通知服务
func NewNotificationServiceFromConfig(config *NotificationConfig) (*NotificationService, error) {
	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	var emailProvider provider.EmailProvider
	var smsProvider provider.SMSProvider

	// 初始化邮箱提供者
	if config.Email != nil && config.Email.Enabled {
		var err error
		emailProvider, err = createEmailProvider(config.Email)
		if err != nil {
			return nil, fmt.Errorf("创建邮箱提供者失败: %w", err)
		}
	}

	// 初始化短信提供者
	if config.SMS != nil && config.SMS.Enabled {
		var err error
		smsProvider, err = createSMSProvider(config.SMS)
		if err != nil {
			return nil, fmt.Errorf("创建短信提供者失败: %w", err)
		}
	}

	return NewNotificationService(emailProvider, smsProvider), nil
}

// createEmailProvider 创建邮箱提供者
func createEmailProvider(config *EmailConfig) (provider.EmailProvider, error) {
	switch config.Provider {
	case "smtp", "":
		return email.NewSMTPProvider(config.SMTP)
	default:
		return nil, fmt.Errorf("不支持的邮箱提供者: %s", config.Provider)
	}
}

// createSMSProvider 创建短信提供者
func createSMSProvider(config *SMSConfig) (provider.SMSProvider, error) {
	switch config.Provider {
	case "tencent":
		return tencent.NewTencentProvider(config.Tencent)
	case "aliyun":
		return aliyun.NewAliyunProvider(config.Aliyun)
	default:
		return nil, fmt.Errorf("不支持的短信提供者: %s", config.Provider)
	}
}
