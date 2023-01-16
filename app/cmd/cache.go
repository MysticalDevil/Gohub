package cmd

import (
	"github.com/spf13/cobra"
	"gohub/pkg/cache"
	"gohub/pkg/console"
)

var Cache = &cobra.Command{
	Use:   "cache",
	Short: "Cache management",
}

var CacheClear = &cobra.Command{
	Use:   "clear",
	Short: "Clear cache",
	Run:   runCacheClear,
}

func init() {
	Cache.AddCommand(CacheClear)
}

func runCacheClear(cmd *cobra.Command, args []string) {
	cache.Flush()
	console.Success("Cache cleared.")
}
