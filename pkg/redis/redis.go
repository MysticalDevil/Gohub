// Package redis Toolkit
package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gohub/pkg/logger"
	"sync"
	"time"
)

// Client Redis serve
type Client struct {
	Client  *redis.Client
	Context context.Context
}

// once Make sure the global Redis object is only instanced once
var once sync.Once

// Redis Global Redis, use db 1
var Redis *Client

// ConnectRedis Connect to the redis database and set the global Redis object
func ConnectRedis(address, username, password string, db int) {
	once.Do(func() {
		Redis = NewClient(address, username, password, db)
	})
}

// NewClient Create a new redis connection
func NewClient(address, username, password string, db int) *Client {
	// Initialize a custom Client instance
	rds := &Client{}
	// Use default context
	rds.Context = context.Background()

	// Use the NewClient in the redis library to initialize the connection
	rds.Client = redis.NewClient(&redis.Options{
		Addr:     address,
		Username: username,
		Password: password,
		DB:       db,
	})

	// Test connection
	err := rds.Ping()
	logger.LogIf(err)

	return rds
}

// Ping Used to test whether the redis connection is normal
func (rds Client) Ping() error {
	_, err := rds.Client.Ping(rds.Context).Result()
	return err
}

// Set Store the value corresponding to the key and set the expiration time
func (rds Client) Set(key string, value any, expiration time.Duration) bool {
	if err := rds.Client.Set(rds.Context, key, value, expiration).Err(); err != nil {
		logger.ErrorString("Redis", "Set", err.Error())
		return false
	}
	return true
}

// Get the value corresponding to the key
func (rds Client) Get(key string) string {
	result, err := rds.Client.Get(rds.Context, key).Result()
	if err != nil {
		if err != redis.Nil {
			logger.ErrorString("Redis", "Get", err.Error())
		}
		return ""
	}
	return result
}

// Has Check if a key exists
func (rds Client) Has(key string) bool {
	_, err := rds.Client.Get(rds.Context, key).Result()
	if err != nil {
		if err != redis.Nil {
			logger.ErrorString("Redis", "Has", err.Error())
		}
		return false
	}
	return true
}

// Del Delete data stored in redis, support multiple key parameters
func (rds Client) Del(keys ...string) bool {
	if err := rds.Client.Del(rds.Context, keys...).Err(); err != nil {
		logger.ErrorString("Redis", "Del", err.Error())
		return false
	}
	return true
}

// FlushDB Clear all data in current redis db
func (rds Client) FlushDB() bool {
	if err := rds.Client.FlushDB(rds.Context).Err(); err != nil {
		logger.ErrorString("Redis", "FlushDB", err.Error())
		return false
	}
	return true
}

// Increment
// When there is only 1 parameter, it is the key, and its value is increased by 1.
// When there are 2 parameters, the first parameter is the key,
// and the second parameter is the int64 type of the value to be added.
func (rds Client) Increment(parameters ...any) bool {
	switch len(parameters) {
	case 1:
		key := parameters[0].(string)
		if err := rds.Client.Incr(rds.Context, key).Err(); err != nil {
			logger.ErrorString("Redis", "Increment", err.Error())
			return false
		}
	case 2:
		key := parameters[0].(string)
		value := parameters[1].(int64)
		if err := rds.Client.IncrBy(rds.Context, key, value).Err(); err != nil {
			logger.ErrorString("Redis", "Increment", err.Error())
			return false
		}
	default:
		logger.ErrorString("Redis", "Increment", "Too many parameters")
		return false
	}
	return true
}

// Decrement The value of key is decremented by one
func (rds Client) Decrement(parameters ...any) bool {
	switch len(parameters) {
	case 1:
		key := parameters[0].(string)
		if err := rds.Client.Decr(rds.Context, key).Err(); err != nil {
			logger.ErrorString("Redis", "Increment", err.Error())
			return false
		}
	case 2:
		key := parameters[0].(string)
		value := parameters[1].(int64)
		if err := rds.Client.DecrBy(rds.Context, key, value).Err(); err != nil {
			logger.ErrorString("Redis", "Increment", err.Error())
			return false
		}
	default:
		logger.ErrorString("Redis", "Increment", "Too many parameters")
		return false
	}
	return true
}
