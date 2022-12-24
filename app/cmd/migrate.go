package cmd

import (
	"github.com/spf13/cobra"
	"gohub/database/migrations"
	"gohub/pkg/migrate"
)

var Migrate = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migration",
}

func init() {
	Migrate.AddCommand(
		MigrateUp,
		MigrateRollback,
		MigrateReset,
		MigrateRefresh,
		MigrateFresh,
	)
}

func migrator() *migrate.Migrator {
	// Register all migration files under database/migrations
	migrations.Initialize()
	// Initialize migrator
	return migrate.NewMigrator()
}

var MigrateUp = &cobra.Command{
	Use:   "up",
	Short: "Run unmigrated migrations",
	Run:   runUp,
}

func runUp(cmd *cobra.Command, args []string) {
	migrator().Up()
}

var MigrateRollback = &cobra.Command{
	Use:     "down",
	Aliases: []string{"rollback"},
	Short:   "Reverse the up command",
	Run:     runDown,
}

func runDown(cmd *cobra.Command, args []string) {
	migrator().Rollback()
}

var MigrateReset = &cobra.Command{
	Use:   "reset",
	Short: "Rollback all database migrations",
	Run:   runReset,
}

func runReset(cmd *cobra.Command, args []string) {
	migrator().Reset()
}

var MigrateRefresh = &cobra.Command{
	Use:   "refresh",
	Short: "Reset and re-run all migrations",
	Run:   runRefresh,
}

func runRefresh(cmd *cobra.Command, args []string) {
	migrator().Refresh()
}

var MigrateFresh = &cobra.Command{
	Use:   "fresh",
	Short: "Drop all tables and re-run all migrations",
	Run:   runFresh,
}

func runFresh(cmd *cobra.Command, args []string) {
	migrator().Fresh()
}
