package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gohub/app/cmd"
	"gohub/bootstrap"
	btsConfig "gohub/config"
	"gohub/pkg/config"
	"gohub/pkg/console"
	"os"
)

func init() {
	// Load the configuration information in the config directory
	btsConfig.Initialize()
}

func main() {
	// Application entry, the command cmd.CmdServer is called by default
	var rootCmd = &cobra.Command{
		Use:   "Gohub",
		Short: "A simple forum project",
		Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,

		// All subcommands of rootCmd execute the following code
		PersistentPreRun: func(command *cobra.Command, args []string) {
			config.InitConfig(cmd.Env)

			// Initialize Logger
			bootstrap.SetupLogger()

			// Initialize DB
			bootstrap.SetupDB()

			// Initialize Redis
			bootstrap.SetupRedis()

			// Initialize the cache
		},
	}

	// register subcommand
	rootCmd.AddCommand(
		cmd.Serve,
		cmd.Key,
		cmd.Play,
	)

	// Configure the web service to run by default
	cmd.RegisterDefaultCmd(rootCmd, cmd.Serve)

	// Register global parameters, --env
	cmd.RegisterGlobalFlags(rootCmd)

	// Execute the main command
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %x", os.Args, err.Error()))
	}
}
