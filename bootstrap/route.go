// Package bootstrap Handler initialization logic
package bootstrap

import (
	"github.com/gin-gonic/gin"
	"gohub/routes"
	"net/http"
	"strings"
)

// SetupRoute Initialize route
func SetupRoute(router *gin.Engine) {
	// Register global middleware
	registerGlobalMiddleware(router)

	// Register APi routes
	routes.RegisterAPIRoutes(router)

	// Configure 404 routing
	setup404Handler(router)
}

func registerGlobalMiddleware(router *gin.Engine) {
	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)
}

func setup404Handler(router *gin.Engine) {
	// Handling 404 request
	router.NoRoute(func(c *gin.Context) {
		// Get the 'Accept' information of the header
		acceptString := c.Request.Header.Get("Accept")
		if strings.Contains(acceptString, "text/html") {
			// If it is HTML
			c.String(http.StatusNotFound, "Page return 404")
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code": 404,
				"error_message": "This route is not defined, " +
					"please confirm whether the url and request method ar correct",
			})
		}
	})
}
