package tencent

import (
	"context"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

	"universal-service-user/notification/domain/provider"
	smspkg "universal-service-user/notification/infrastructure/sms"
)

// TencentConfig 腾讯云短信配置
type TencentConfig struct {
	SecretID  string            // 腾讯云 SecretId
	SecretKey string            // 腾讯云 SecretKey
	AppID     string            // 短信应用 ID
	Sign      string            // 短信签名
	Templates map[string]string // 模板映射(key: 场景, value: 模板 ID)
}

// TencentProvider 腾讯云短信提供者
type TencentProvider struct {
	config *TencentConfig
}

// NewTencentProvider 创建腾讯云短信提供者
func NewTencentProvider(config *TencentConfig) (*TencentProvider, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	p := &TencentProvider{
		config: config,
	}

	// 注册到全局注册中心
	if err := smspkg.Register(p); err != nil {
		return nil, fmt.Errorf("注册腾讯云短信提供者失败: %w", err)
	}

	return p, nil
}

// GetName 获取提供者唯一标识
func (p *TencentProvider) GetName() string {
	return "tencent"
}

// GetDisplayName 获取提供者显示名称
func (p *TencentProvider) GetDisplayName() string {
	return "腾讯云短信"
}

// Send 发送短信
func (p *TencentProvider) Send(ctx context.Context, message *provider.SMSMessage) error {
	// 创建凭证
	credential := common.NewCredential(
		p.config.SecretID,
		p.config.SecretKey,
	)

	// 创建客户端配置
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"

	// 创建腾讯云短信客户端
	client, err := sms.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		return fmt.Errorf("创建腾讯云短信客户端失败: %w", err)
	}

	// 创建发送短信请求
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(p.config.AppID)
	request.SignName = common.StringPtr(p.config.Sign)
	request.TemplateId = common.StringPtr(message.TemplateID)
	request.PhoneNumberSet = common.StringPtrs([]string{"+86" + message.PhoneNumber})

	// 转换模板参数
	if len(message.Params) > 0 {
		params := make([]*string, 0, len(message.Params))
		for _, v := range message.Params {
			params = append(params, common.StringPtr(v))
		}
		request.TemplateParamSet = params
	}

	// 发送短信
	response, err := client.SendSms(request)
	if err != nil {
		return fmt.Errorf("发送短信失败: %w", err)
	}

	// 检查发送结果
	if len(response.Response.SendStatusSet) > 0 {
		status := response.Response.SendStatusSet[0]
		if *status.Code != "Ok" {
			return fmt.Errorf("短信发送失败: %s", *status.Message)
		}
	}

	return nil
}

// SendWithTemplate 使用模板发送短信
func (p *TencentProvider) SendWithTemplate(ctx context.Context, phoneNumber, templateID string, params map[string]string) error {
	message := &provider.SMSMessage{
		PhoneNumber: phoneNumber,
		TemplateID:  templateID,
		Params:      params,
		SignName:    p.config.Sign,
	}

	return p.Send(ctx, message)
}

// ValidateConfig 验证配置是否完整有效
func (p *TencentProvider) ValidateConfig() error {
	return validateConfig(p.config)
}

// validateConfig 验证配置
func validateConfig(config *TencentConfig) error {
	if config == nil {
		return fmt.Errorf("腾讯云短信配置不能为空")
	}
	if config.SecretID == "" {
		return fmt.Errorf("SecretID 不能为空")
	}
	if config.SecretKey == "" {
		return fmt.Errorf("SecretKey 不能为空")
	}
	if config.AppID == "" {
		return fmt.Errorf("AppID 不能为空")
	}
	if config.Sign == "" {
		return fmt.Errorf("Sign 不能为空")
	}

	return nil
}
