package cmd

import (
    "fmt"
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

var CacheForget = &cobra.Command{
    Use:   "forget",
    Short: "Delete redis key, example: cache forget cache-key",
    Run:   runCacheForget,
}

// Options for the forget Command
var cacheKey string

func init() {
    Cache.AddCommand(CacheClear, CacheForget)

    // Set options for the cache forget command
    CacheForget.Flags().StringVarP(&cacheKey, "key", "k", "", "KEY of the cache")
    _ = CacheForget.MarkFlagRequired("key")
}

func runCacheClear(cmd *cobra.Command, args []string) {
    cache.Flush()
    console.Success("Cache cleared.")
}

func runCacheForget(cmd *cobra.Command, args []string) {
    cache.Forget(cacheKey)
    console.Success(fmt.Sprintf("Cache key [%s] deleted.", cacheKey))
}
