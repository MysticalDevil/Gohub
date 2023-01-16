package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gohub/pkg/response"
)

// ForceUA middleware, Mandatory must be accompanied by User-Agent header
func ForceUA() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get User-Agent header information
		if len(c.Request.Header["User-Agent"]) == 0 {
			response.BadRequest(
				c,
				errors.New("User-Agent header not found"),
				"Requests must be accompanied by a User-Agent header",
			)
		}

		c.Next()
	}
}
