// Package routes Register route
package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterAPIRoutes Registration page related routing
func RegisterAPIRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"Hello": "World",
			})
		})
	}
}
