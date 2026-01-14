package aliyun

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"

	"github.com/2928807938/universal-service-user/notification/domain/provider"
	"github.com/2928807938/universal-service-user/notification/infrastructure/sms"
)

// AliyunConfig 阿里云短信配置
type AliyunConfig struct {
	AccessKeyID     string            // 阿里云 AccessKeyId
	AccessKeySecret string            // 阿里云 AccessKeySecret
	SignName        string            // 短信签名
	Templates       map[string]string // 模板映射(key: 场景, value: 模板 ID)
}

// AliyunProvider 阿里云短信提供者
type AliyunProvider struct {
	config *AliyunConfig
}

// NewAliyunProvider 创建阿里云短信提供者
func NewAliyunProvider(config *AliyunConfig) (*AliyunProvider, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	p := &AliyunProvider{
		config: config,
	}

	// 注册到全局注册中心
	if err := sms.Register(p); err != nil {
		return nil, fmt.Errorf("注册阿里云短信提供者失败: %w", err)
	}

	return p, nil
}

// GetName 获取提供者唯一标识
func (p *AliyunProvider) GetName() string {
	return "aliyun"
}

// GetDisplayName 获取提供者显示名称
func (p *AliyunProvider) GetDisplayName() string {
	return "阿里云短信"
}

// Send 发送短信
func (p *AliyunProvider) Send(ctx context.Context, message *provider.SMSMessage) error {
	// 创建阿里云短信客户端
	client, err := dysmsapi.NewClientWithAccessKey(
		"cn-hangzhou", // 地域 ID,默认使用杭州
		p.config.AccessKeyID,
		p.config.AccessKeySecret,
	)
	if err != nil {
		return fmt.Errorf("创建阿里云短信客户端失败: %w", err)
	}

	// 创建发送短信请求
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = message.PhoneNumber
	request.SignName = p.config.SignName
	request.TemplateCode = message.TemplateID

	// 转换参数为 JSON
	if len(message.Params) > 0 {
		paramsJSON, err := json.Marshal(message.Params)
		if err != nil {
			return fmt.Errorf("序列化参数失败: %w", err)
		}
		request.TemplateParam = string(paramsJSON)
	}

	// 发送短信
	response, err := client.SendSms(request)
	if err != nil {
		return fmt.Errorf("发送短信失败: %w", err)
	}

	// 检查发送结果
	if response.Code != "OK" {
		return fmt.Errorf("短信发送失败: %s - %s", response.Code, response.Message)
	}

	return nil
}

// SendWithTemplate 使用模板发送短信
func (p *AliyunProvider) SendWithTemplate(ctx context.Context, phoneNumber, templateID string, params map[string]string) error {
	message := &provider.SMSMessage{
		PhoneNumber: phoneNumber,
		TemplateID:  templateID,
		Params:      params,
		SignName:    p.config.SignName,
	}

	return p.Send(ctx, message)
}

// ValidateConfig 验证配置是否完整有效
func (p *AliyunProvider) ValidateConfig() error {
	return validateConfig(p.config)
}

// validateConfig 验证配置
func validateConfig(config *AliyunConfig) error {
	if config == nil {
		return fmt.Errorf("阿里云短信配置不能为空")
	}
	if config.AccessKeyID == "" {
		return fmt.Errorf("AccessKeyID 不能为空")
	}
	if config.AccessKeySecret == "" {
		return fmt.Errorf("AccessKeySecret 不能为空")
	}
	if config.SignName == "" {
		return fmt.Errorf("SignName 不能为空")
	}

	return nil
}
