package config

import "gohub/pkg/config"

func init() {
	config.Add("captcha", func() map[string]any {
		return map[string]any{
			"height":            80,
			"width":             240,
			"length":            6,
			"maxSkew":           0.7,
			"dotCount":          80,
			"expire_time":       15,
			"debug_expire_time": 10080,
			"testing_key":       "captcha_skip_test",
		}
	})
}
