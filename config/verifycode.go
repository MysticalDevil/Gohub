package config

import "gohub/pkg/config"

func init() {
	config.Add("verifycode", func() map[string]any {
		return map[string]any{
			// Verify code length
			"code_length": config.Env("VERIFY_CODE_LENGTH", 6),

			// Expiration time, in minutes
			"expire_time": config.Env("VERIFY_CODE_EXPIRE", 15),

			// Expiration time, in minutes
			"debug_expire_time": 10080,
			// The local development environment verification code uses debug_code
			"debug_code": 123456,

			"debug_phone_prefix": "000",
			"debug_email_suffix": "@testing.com",
		}
	})
}
