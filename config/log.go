package config

import "gohub/pkg/config"

func init() {
	config.Add("log", func() map[string]any {
		return map[string]any{
			// Log level
			"level": config.Env("LOG_LEVEL", "debug"),
			// Log type
			"type": config.Env("LOG_TYPE", "single"),

			/*---------- Rolling log configuration ----------*/
			// Log file path
			"filename": config.Env("LOG_NAME", "storage/logs/logs.log"),
			// The maximum size of each log file, unit: M
			"max_size": config.Env("LOG_MAX_SIZE", 64),
			// Maximum number of saved files
			"max_backup": config.Env("LOG_MAX_BACKUP", 5),
			// Maximum storage time, 0 means not to delete
			"max_age": config.Env("LOG_MAX_AGE", 30),
			// Whether to compress
			"compress": config.Env("LOG_COMPRESS", false),
		}
	})
}
