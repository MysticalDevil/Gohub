package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gohub/pkg/app"
	"gohub/pkg/limiter"
	"gohub/pkg/logger"
	"gohub/pkg/response"
)

// LimitIP Global current limiting middleware, limiting current for IP
// limit is a format string, such as "5-S", example:
//
// * 5 reqs/second: "5-S"
// * 10 reqs/minute: "10-M"
// * 1000 reqs/hour: "1000-H"
// * 2000 reqs/day: "2000-D"
func LimitIP(limit string) gin.HandlerFunc {
	if app.IsTesting() {
		limit = "1000000-H"
	}

	return func(c *gin.Context) {
		// Current limit for IP
		key := limiter.GetKeyIP(c)
		if ok := limitHandler(c, key, limit); !ok {
			return
		}
		c.Next()
	}
}

// LimitPerRoute Throttle middleware, used in a separate route
func LimitPerRoute(limit string) gin.HandlerFunc {
	if app.IsTesting() {
		limit = "1000000-H"
	}
	return func(c *gin.Context) {
		// For a single route, increase the number of visits
		c.Set("limiter-once", false)

		// Limit traffic for IP + routing
		key := limiter.GetKeyRouteWithIP(c)
		if ok := limitHandler(c, key, limit); !ok {
			return
		}
		c.Next()
	}
}

func limitHandler(c *gin.Context, key string, limit string) bool {
	// get excess
	rate, err := limiter.CheckRate(c, key, limit)
	if err != nil {
		logger.LogIf(err)
		response.Abort500(c)
		return false
	}

	// ---- Set header information ----
	// X-RateLimit-Limit :10000 Maximum number of visits
	// X-RateLimit-Remaining :9993 Visits remaining
	// X-RateLimit-Reset :1513784506 At this point in time, the number of visits will be reset to X-RateLimit-Limit
	c.Header("X-RateLimit-Limit", cast.ToString(rate.Limit))
	c.Header("X-RateLimit-Remaining", cast.ToString(rate.Remaining))
	c.Header("X-RateLimit-Reset", cast.ToString(rate.Reset))

	// excess
	if rate.Reached {
		// Notify the user that the quota is exceeded
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"message": "Interface requests are too frequent",
		})
		return false
	}
	return true
}
