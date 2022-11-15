package middlewares

import (
	"github.com/gin-gonic/gin"
	"gohub/pkg/jwt"
	"gohub/pkg/response"
)

// GuestJWT Force users to access as guest
func GuestJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.GetHeader("Authorization")) > 0 {
			_, err := jwt.NewJWT().ParseToken(c)
			if err == nil {
				response.Unauthorized(c, "Please visit as a tourist")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
