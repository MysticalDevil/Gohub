package cache

import "time"

type Store interface {
	Set(key, value string, expireTime time.Duration)
	Get(key string) string
	Has(key string) bool
	Forget(key string)
	Forever(key, value string)
	Flush()

	IsAlive() error

	// Increment
	// When there is only one parameter, add 1 to the key
	// When there are two parameters, the first parameter is the key,
	// and the second parameter is the value to be added (int64 type)
	Increment(parameters ...any)

	// Decrement
	// When there is only one parameter, sub 1 to the key
	// When there are two parameters, the first parameter is the key,
	// and the second parameter is the value to be subtracted (int64 type)
	Decrement(parameters ...any)
}
