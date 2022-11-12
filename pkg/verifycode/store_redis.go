package verifycode

import (
	"gohub/pkg/app"
	"gohub/pkg/config"
	"gohub/pkg/redis"
	"time"
)

// RedisStore Implement the verifycode.Store interface
type RedisStore struct {
	RedisClient *redis.Client
	KeyPrefix   string
}

// Set Implement the Set method of verifycode.Store
func (s *RedisStore) Set(key string, value string) bool {
	ExpireTime := time.Minute * time.Duration(config.GetInt64("verifycode.expire_time"))

	if app.IsLocal() {
		ExpireTime = time.Minute * time.Duration(config.GetInt64("verifycode.debug_expire_time"))
	}

	return s.RedisClient.Set(s.KeyPrefix+key, value, ExpireTime)
}

// Get Implement the Get method of verifycode.Store
func (s *RedisStore) Get(key string, clear bool) (value string) {
	key = s.KeyPrefix + key
	value = s.RedisClient.Get(key)
	if clear {
		s.RedisClient.Del(key)
	}
	return
}

// Verify Implement the Verify method of verifycode.Store
func (s *RedisStore) Verify(key, answer string, clear bool) bool {
	v := s.Get(key, clear)
	return v == answer
}
