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

func runUp(_ *cobra.Command, _ []string) {
	migrator().Up()
}

var MigrateRollback = &cobra.Command{
	Use:     "down",
	Aliases: []string{"rollback"},
	Short:   "Reverse the up command",
	Run:     runDown,
}

func runDown(_ *cobra.Command, _ []string) {
	migrator().Rollback()
}

var MigrateReset = &cobra.Command{
	Use:   "reset",
	Short: "Rollback all database migrations",
	Run:   runReset,
}

func runReset(_ *cobra.Command, _ []string) {
	migrator().Reset()
}

var MigrateRefresh = &cobra.Command{
	Use:   "refresh",
	Short: "Reset and re-run all migrations",
	Run:   runRefresh,
}

func runRefresh(_ *cobra.Command, _ []string) {
	migrator().Refresh()
}

var MigrateFresh = &cobra.Command{
	Use:   "fresh",
	Short: "Drop all tables and re-run all migrations",
	Run:   runFresh,
}

func runFresh(_ *cobra.Command, _ []string) {
	migrator().Fresh()
}
