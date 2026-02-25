package bootstrap

import (
	"strings"

	"gohub/pkg/config"
	"gohub/pkg/console"
	"gohub/pkg/logger"
)

// SetupLogger Initialize Logger
func SetupLogger() {
	// Fail fast in non-local environments if APP_KEY is missing or a placeholder.
	appKey := strings.TrimSpace(config.GetString("app.key"))
	if appKey == "" || appKey == "yourSecret" {
		env := config.GetString("app.env")
		if env != "local" && env != "development" && env != "test" && env != "testing" {
			console.Exit("APP_KEY is required for non-local environments; run `go run main.go key` and set APP_KEY in .env")
		}
	}

	logger.InitLogger(
		config.GetString("log.filename"),
		config.GetInt("log.max_size"),
		config.GetInt("log.max_backup"),
		config.GetInt("log.max_age"),
		config.GetBool("log.compress"),
		config.GetString("log.type"),
		config.GetString("log.level"),
	)
}
