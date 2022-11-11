// Package config Site configuration information
package config

import "gohub/pkg/config"

func init() {
	config.Add("app", func() map[string]any {
		return map[string]any{
			// App name
			"name": config.Env("APP_NAME", "Gohub"),
			// The current environment, used to distinguish multiple environments, generally local, stage, production, test
			"env": config.Env("APP_ENV", "production"),
			// Whether to enter debug mode
			"debug": config.Env("APP_DEBUG", false),
			// Application service port
			"port": config.Env("APP_PORT", "3000"),
			// Encrypted session, JWT encryption
			"key": config.Env("APP_KEY", "8f9603f2be2480ecc6ac3a47fc4888e15be260e290b865582c613fc7b065c02a"),
			// To generate links
			"url": config.Env("APP_URL", "http://localhost:3000"),
			// Set time zone
			"timezone": config.Env("TIMEZONE", "Asia/Shanghai"),
		}
	})
}
