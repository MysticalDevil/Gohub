package cache

import (
	"gohub/pkg/config"
	"gohub/pkg/redis"
	"time"
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
	s.RedisClient.Increment(parameters...)
}

func (s *RedisStore) Decrement(parameters ...any) {
	s.RedisClient.Decrement(parameters...)
}
