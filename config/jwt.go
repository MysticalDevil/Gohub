package config

import "gohub/pkg/config"

func init() {
	config.Add("jwt", func() map[string]any {
		return map[string]any{
			// use config.GetString("app.key")
			// "signing_key":

			// Expiration time, in minutes, generally no more than two hours
			"expire_time": config.Env("JWT_EXPIRE_TIME", 120),
			// Allowed refresh time, the unit is minutes, 86400 is two months, counted from the token signature time
			"max_refresh_time": config.Env("JWT_MAX_REFRESH_TIME", 86400),
			// The expiration time in debug mode is convenient for local debugging and development
			"debug_expire_time": 86400,
		}
	})
}
