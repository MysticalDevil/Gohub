package config

import "gohub/pkg/config"

func init() {
	config.Add("sms", func() map[string]any {
		return map[string]any{
			"aliyun": map[string]any{
				"access_key_id":     config.Env("SMS_ALIYUN_ACCESS_ID"),
				"access_key_secret": config.Env("SMS_ALIYUN_ACCESS_SECRET"),
				"sign_name":         config.Env("SMS_ALIYUN_SIGN_NAME"),
				"template_code":     config.Env("SMS_ALIYUN_TEMPLATE_CODE", "SMS_154950909"),
			},
		}
	})
}
