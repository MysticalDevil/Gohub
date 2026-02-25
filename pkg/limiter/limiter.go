// Package limiter Handle throttling logic
package limiter

import (
	"strings"

	"github.com/gin-gonic/gin"
	limiterLib "github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"gohub/pkg/config"
	"gohub/pkg/logger"
	"gohub/pkg/redis"
)

// GetKeyIP Get Limiter's Key, IP
func GetKeyIP(c *gin.Context) string {
	return c.ClientIP()
}

// GetKeyRouteWithIP Limiter's Key, route + IP, limits current for a single route
func GetKeyRouteWithIP(c *gin.Context) string {
	return routeToKeyString(c.FullPath() + c.ClientIP())
}

// CheckRate Detect whether the request is overclocked
func CheckRate(c *gin.Context, key, formatted string) (limiterLib.Context, error) {
	// Instantiate the limiter.Rate object of the dependent limiter package
	var context limiterLib.Context
	rate, err := limiterLib.NewRateFromFormatted(formatted)
	if err != nil {
		logger.LogIf(err)
		return context, err
	}

	// Initialize the storage, using the shared redis.Redis object in the program
	store, err := sredis.NewStoreWithOptions(redis.Redis.Client, limiterLib.StoreOptions{
		// Set a prefix for limiter to keep redis data tidy
		Prefix: config.GetString("app.name") + ":limiter",
	})
	if err != nil {
		logger.LogIf(err)
		return context, err
	}

	// Use the limiter.Rate object and storage object initialized above
	limiterObj := limiterLib.New(store, rate)

	// Get the result of the current limit
	if c.GetBool("limiter-once") {
		// Peek() takes the result without increasing the number of visits
		return limiterObj.Peek(c, key)
	} else {
		// Make sure that when LimitIP is called in multiple routing groups to limit traffic,
		// only one visit will be added
		c.Set("limiter-core", true)
		// Get() takes the result and increases the number of visits
		return limiterObj.Get(c, key)
	}
}

// routeToKeyString helper method, format '/' in the URL as '-'
func routeToKeyString(routeName string) string {
	routeName = strings.ReplaceAll(routeName, "/", "-")
	routeName = strings.ReplaceAll(routeName, ":", "_")
	return routeName
}
