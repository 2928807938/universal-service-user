package service

import (
	"fmt"

	"github.com/2928807938/universal-service-user/notification/infrastructure/email"
	"github.com/2928807938/universal-service-user/notification/infrastructure/sms/aliyun"
	"github.com/2928807938/universal-service-user/notification/infrastructure/sms/tencent"
)

// NotificationConfig 通知配置
type NotificationConfig struct {
	Email *EmailConfig // 邮箱配置
	SMS   *SMSConfig   // 短信配置
}

// EmailConfig 邮箱配置
type EmailConfig struct {
	Enabled  bool              // 是否启用
	Provider string            // 提供者类型(当前仅支持 "smtp")
	SMTP     *email.SMTPConfig // SMTP 配置
}

// SMSConfig 短信配置
type SMSConfig struct {
	Enabled  bool                   // 是否启用
	Provider string                 // 提供者类型("tencent" 或 "aliyun")
	Tencent  *tencent.TencentConfig // 腾讯云配置
	Aliyun   *aliyun.AliyunConfig   // 阿里云配置
}

// Validate 验证配置
func (c *NotificationConfig) Validate() error {
	if c == nil {
		return fmt.Errorf("通知配置不能为空")
	}

	// 验证邮箱配置
	if c.Email != nil && c.Email.Enabled {
		if err := c.Email.Validate(); err != nil {
			return fmt.Errorf("邮箱配置验证失败: %w", err)
		}
	}

	// 验证短信配置
	if c.SMS != nil && c.SMS.Enabled {
		if err := c.SMS.Validate(); err != nil {
			return fmt.Errorf("短信配置验证失败: %w", err)
		}
	}

	return nil
}

// Validate 验证邮箱配置
func (c *EmailConfig) Validate() error {
	if c == nil {
		return fmt.Errorf("邮箱配置不能为空")
	}

	if !c.Enabled {
		return nil
	}

	// 目前仅支持 SMTP
	if c.Provider != "smtp" && c.Provider != "" {
		return fmt.Errorf("不支持的邮箱提供者: %s", c.Provider)
	}

	// 验证 SMTP 配置
	if c.SMTP == nil {
		return fmt.Errorf("SMTP 配置不能为空")
	}

	if c.SMTP.Host == "" {
		return fmt.Errorf("SMTP Host 不能为空")
	}
	if c.SMTP.Port <= 0 {
		return fmt.Errorf("SMTP Port 必须大于 0")
	}
	if c.SMTP.Username == "" {
		return fmt.Errorf("SMTP Username 不能为空")
	}
	if c.SMTP.Password == "" {
		return fmt.Errorf("SMTP Password 不能为空")
	}
	if c.SMTP.From == "" {
		return fmt.Errorf("SMTP From 不能为空")
	}

	return nil
}

// Validate 验证短信配置
func (c *SMSConfig) Validate() error {
	if c == nil {
		return fmt.Errorf("短信配置不能为空")
	}

	if !c.Enabled {
		return nil
	}

	// 检查提供者类型
	if c.Provider != "tencent" && c.Provider != "aliyun" {
		return fmt.Errorf("不支持的短信提供者: %s (仅支持 tencent 或 aliyun)", c.Provider)
	}

	// 根据提供者类型验证对应的配置
	switch c.Provider {
	case "tencent":
		if c.Tencent == nil {
			return fmt.Errorf("腾讯云短信配置不能为空")
		}
		if c.Tencent.SecretID == "" {
			return fmt.Errorf("腾讯云 SecretID 不能为空")
		}
		if c.Tencent.SecretKey == "" {
			return fmt.Errorf("腾讯云 SecretKey 不能为空")
		}
		if c.Tencent.AppID == "" {
			return fmt.Errorf("腾讯云 AppID 不能为空")
		}
		if c.Tencent.Sign == "" {
			return fmt.Errorf("腾讯云 Sign 不能为空")
		}

	case "aliyun":
		if c.Aliyun == nil {
			return fmt.Errorf("阿里云短信配置不能为空")
		}
		if c.Aliyun.AccessKeyID == "" {
			return fmt.Errorf("阿里云 AccessKeyID 不能为空")
		}
		if c.Aliyun.AccessKeySecret == "" {
			return fmt.Errorf("阿里云 AccessKeySecret 不能为空")
		}
		if c.Aliyun.SignName == "" {
			return fmt.Errorf("阿里云 SignName 不能为空")
		}
	}

	return nil
}

// DefaultConfig 返回默认配置
func DefaultConfig() *NotificationConfig {
	return &NotificationConfig{
		Email: &EmailConfig{
			Enabled:  false,
			Provider: "smtp",
		},
		SMS: &SMSConfig{
			Enabled:  false,
			Provider: "tencent",
		},
	}
}
