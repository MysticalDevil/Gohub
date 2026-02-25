package captcha

import (
	"errors"
	"time"

	"gohub/pkg/app"
	"gohub/pkg/config"
	"gohub/pkg/redis"
)

// RedisStore Implement the base64Captcha.Store interface
type RedisStore struct {
	RedisClient *redis.Client
	KeyPrefix   string
}

// Set Implement the Set method of the base64Captcha.Store interface
func (s *RedisStore) Set(key, value string) error {
	ExpireTime := time.Minute * time.Duration(config.GetInt64("captcha.expire_time"))
	// Convenient for local development and debugging
	if app.IsLocal() {
		ExpireTime = time.Minute * time.Duration(config.GetInt64("captcha.debug_expire_time"))
	}

	if ok := s.RedisClient.Set(s.KeyPrefix+key, value, ExpireTime); !ok {
		return errors.New("unable to store image captcha answer")
	}
	return nil
}

// Get Implement the Get method of the base64Captcha.Store interface
func (s *RedisStore) Get(key string, clear bool) string {
	key = s.KeyPrefix + key
	val := s.RedisClient.Get(key)
	if clear {
		s.RedisClient.Del(key)
	}
	return val
}

// Verify Implement the Verify method of the base64Captcha.Store interface
func (s *RedisStore) Verify(key, answer string, clear bool) bool {
	v := s.Get(key, clear)
	return v == answer
}
