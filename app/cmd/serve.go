package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gohub/bootstrap"
	"gohub/pkg/config"
	"gohub/pkg/console"
	"gohub/pkg/logger"
)

// Serve represents the available web sub-command
var Serve = &cobra.Command{
	Use:   "serve",
	Short: "Start web server",
	Run:   runWeb,
	Args:  cobra.NoArgs,
}

func runWeb(cmd *cobra.Command, args []string) {
	// Set the running mode of gin, support debug, release, test
	// release mod will block debugging information
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	bootstrap.SetupRoute(router)

	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		logger.ErrorString("CMD", "serve", err.Error())
		console.Exit("Unable to start server, error:" + err.Error())
	}
}
