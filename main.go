package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"gohub/bootstrap"
	btsConfig "gohub/config"
	"gohub/pkg/config"
	"log"
)

func init() {
	// Load the configuration information in the config directory
	btsConfig.Initialize()
}

func main() {
	// Configuration initialization, depends on the command line --env parameter
	var env string
	flag.StringVar(&env, "env", "", "Load .env file. " +
		"For example --env=testing, which loads the .env.testing file")
	flag.Parse()
	config.InitConfig(env)

	// Create a new Gin Engine instance
	router := gin.New()
	// Initialize route binding
	bootstrap.SetupRoute(router)
	// Run serve
	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		// Error handling, port occupied or other errors
		log.Println(err.Error())
	}
}
