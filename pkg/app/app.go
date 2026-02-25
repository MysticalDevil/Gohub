// Package app Application information
package app

import (
	"time"

	"gohub/pkg/config"
)

func IsLocal() bool {
	return config.Get("app.env") == "local"
}

func IsProduction() bool {
	return config.Get("app.env") == "production"
}

func IsTesting() bool {
	return config.Get("app.env") == "testing"
}

// TimenowInTimezone Signature effective time
func TimenowInTimezone() time.Time {
	chinaTimezone, _ := time.LoadLocation(config.GetString("app.timezone"))
	return time.Now().In(chinaTimezone)
}

// URL Pass the path parameter to splice the URL of the site
func URL(path string) string {
	return config.Get("app.url") + path
}

// V1URL Splicing with v1 marker URL
func V1URL(path string) string {
	return URL("/v1/" + path)
}
