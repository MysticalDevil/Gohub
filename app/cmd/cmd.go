// Package cmd Store all sub-commands of the program
package cmd

import (
	"github.com/spf13/cobra"
	"gohub/pkg/helpers"
	"os"
)

// Env Store the value of the global option --env
var Env string

// RegisterGlobalFlags Register global options (flag)
func RegisterGlobalFlags(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringVarP(
		&Env,
		"env",
		"e",
		"",
		"load .env file, example: --env=testing will use .env.testing file",
	)
}

// RegisterDefaultCmd Register default commands
func RegisterDefaultCmd(rootCmd *cobra.Command, subCmd *cobra.Command) {
	cmd, _, err := rootCmd.Find(os.Args[1:])
	firstArg := helpers.FirstElement(os.Args[1:])
	if err == nil && cmd.Use == rootCmd.Use && firstArg != "-h" && firstArg != "--help" {
		args := append([]string{subCmd.Use}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}
}
