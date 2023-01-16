package config

import "gohub/pkg/config"

func init() {
	config.Add("redis", func() map[string]any {
		return map[string]any{
			"host":     config.Env("REDIS_HOST", "127.0.0.1"),
			"port":     config.Env("REDIS_PORT", "6379"),
			"password": config.Env("REDIS_PASSWORD", ""),
			// Use 1 for business class storage (picture verification code, SMS verification code, session)
			"database": config.Env("REDIS_MAIN_DB", 1),

			// Use 0 for the cache package, and clearing the cache should not affect the business
			"database_cache": config.Env("REDIS_CACHE_DB", 0),
		}
	})
}
