package config

import "gohub/pkg/config"

func init() {
	config.Add("database", func() map[string]any {
		return map[string]any{
			// Default database
			"connection":           config.Env("DB_CONNECTION", "postgresql"),
			"max_idle_connections": config.Env("DB_MAX_IDLE_CONNECTIONS", 100),
			"max_open_connections": config.Env("DB_MAX_OPEN_CONNECTIONS", 25),
			"max_life_seconds":     config.Env("DB_MAX_LIFE_SECONDS", 5*60),
			"mysql": map[string]any{
				"host":     config.Env("DB_HOST", "127.0.0.1"),
				"port":     config.Env("DB_PORT", "3306"),
				"database": config.Env("DB_DATABASE", "gohub"),
				"username": config.Env("DB_USERNAME", ""),
				"password": config.Env("DB_PASSWORD", ""),
				"charset":  "utf8mb4",
			},
			"postgresql": map[string]any{
				"host":     config.Env("DB_HOST", "127.0.0.1"),
				"port":     config.Env("DB_PORT", "5432"),
				"database": config.Env("DB_DATABASE", "gohub"),
				"username": config.Env("DB_USERNAME", ""),
				"password": config.Env("DB_PASSWORD", ""),
				"timezone": config.Env("DB_TIMEZONE", "Asia/Shanghai"),
			},
			"sqlite": map[string]any{
				"database": config.Env("DB_SQL_FILE", "database/database.db"),
			},
		}
	})
}
