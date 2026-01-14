package oauth

import (
	"github.com/2928807938/universal-service-user/oauth/provider"
	"github.com/2928807938/universal-service-user/oauth/provider/alipay"
	"github.com/2928807938/universal-service-user/oauth/provider/qq"
	"github.com/2928807938/universal-service-user/oauth/provider/wechat"
)

// RegisterWechat 注册微信 Provider
func RegisterWechat(config *provider.WechatConfig) error {
	if config == nil {
		return provider.ErrInvalidConfig
	}

	p := wechat.New(config)
	return provider.GetManager().RegisterProvider(p)
}

// RegisterAlipay 注册支付宝 Provider
func RegisterAlipay(config *provider.AlipayConfig) error {
	if config == nil {
		return provider.ErrInvalidConfig
	}

	p := alipay.New(config)
	return provider.GetManager().RegisterProvider(p)
}

// RegisterQQ 注册 QQ Provider
func RegisterQQ(config *provider.QQConfig) error {
	if config == nil {
		return provider.ErrInvalidConfig
	}

	p := qq.New(config)
	return provider.GetManager().RegisterProvider(p)
}

// RegisterAll 一次性注册所有启用的 Providers
func RegisterAll(configs map[string]interface{}) error {
	// 注册微信
	if wechatCfg, ok := configs["wechat"].(*provider.WechatConfig); ok && wechatCfg != nil && wechatCfg.Enabled {
		if err := RegisterWechat(wechatCfg); err != nil {
			return err
		}
	}

	// 注册支付宝
	if alipayCfg, ok := configs["alipay"].(*provider.AlipayConfig); ok && alipayCfg != nil && alipayCfg.Enabled {
		if err := RegisterAlipay(alipayCfg); err != nil {
			return err
		}
	}

	// 注册 QQ
	if qqCfg, ok := configs["qq"].(*provider.QQConfig); ok && qqCfg != nil && qqCfg.Enabled {
		if err := RegisterQQ(qqCfg); err != nil {
			return err
		}
	}

	return nil
}
