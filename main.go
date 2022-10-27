package main

import (
	"github.com/gin-gonic/gin"
	"gohub/bootstrap"
	"log"
)

func main() {
	router := gin.New()
	// Initialize route binding
	bootstrap.SetupRoute(router)
	// Run serve
	err := router.Run(":3000")
	if err != nil {
		// Error handling, port occupied or other errors
		log.Println(err.Error())
	}
}
