package bootstrap

import (
	"fmt"
	"gohub/pkg/config"
	"gohub/pkg/redis"
)

// SetupRedis Initialize Redis
func SetupRedis() {
	// Establish a Redis connection
	redis.ConnectRedis(
		fmt.Sprintf("%v:%v",
			config.GetString("redis.host"),
			config.GetString("redis.port")),
		config.GetString("redis.username"),
		config.GetString("redis.password"),
		config.GetInt("redis.database"),
	)
}
