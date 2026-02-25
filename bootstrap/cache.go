package bootstrap

import (
	"fmt"

	"gohub/pkg/app"
	"gohub/pkg/cache"
	"gohub/pkg/config"
)

// SetupCache Set up cache
func SetupCache() {
	if app.IsTesting() {
		cache.InitWithCacheStore(cache.NewMemoryStore())
		return
	}

	// Initialize the cache-specific redis client and use the dedicated cache DB
	rds := cache.NewRedisStore(
		fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port")),
		config.GetString("redis.username"),
		config.GetString("redis.password"),
		config.GetInt("redis.database_cache"),
	)

	cache.InitWithCacheStore(rds)
}
