package cache

import (
	"time"

	"gohub/pkg/config"
	"gohub/pkg/redis"
)

// RedisStore Implement the cache.Store interface
type RedisStore struct {
	RedisClient *redis.Client
	KeyPrefix   string
}

func NewRedisStore(address string, username string, password string, db int) *RedisStore {
	rs := &RedisStore{}
	rs.RedisClient = redis.NewClient(address, username, password, db)
	rs.KeyPrefix = config.GetString("app.name") + ":cache:"
	return rs
}

func (s *RedisStore) Set(key, value string, expireTime time.Duration) {
	s.RedisClient.Set(s.KeyPrefix+key, value, expireTime)
}

func (s *RedisStore) Get(key string) string {
	return s.RedisClient.Get(s.KeyPrefix + key)
}

func (s *RedisStore) Has(key string) bool {
	return s.RedisClient.Has(s.KeyPrefix + key)
}

func (s *RedisStore) Forget(key string) {
	s.RedisClient.Del(s.KeyPrefix + key)
}

func (s *RedisStore) Forever(key, value string) {
	s.RedisClient.Set(s.KeyPrefix+key, value, 0)
}

func (s *RedisStore) Flush() {
	s.RedisClient.FlushDB()
}

func (s *RedisStore) IsAlive() error {
	return s.RedisClient.Ping()
}

func (s *RedisStore) Increment(parameters ...any) {
	s.RedisClient.Increment(s.prefixedParameters(parameters...)...)
}

func (s *RedisStore) Decrement(parameters ...any) {
	s.RedisClient.Decrement(s.prefixedParameters(parameters...)...)
}

func (s *RedisStore) prefixedParameters(parameters ...any) []any {
	if len(parameters) == 0 {
		return parameters
	}

	prefixed := make([]any, len(parameters))
	copy(prefixed, parameters)
	if key, ok := parameters[0].(string); ok {
		prefixed[0] = s.KeyPrefix + key
	}
	return prefixed
}
