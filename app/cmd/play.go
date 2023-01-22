package cmd

import (
	"github.com/spf13/cobra"
)

var Play = &cobra.Command{
	Use:   "play",
	Short: "Likes the Go Playground, but running at our application context",
	Run:   runPlay,
}

func runPlay(_ *cobra.Command, _ []string) {
}
