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

var MigrateUp = &cobra.Command{
	Use:   "up",
	Short: "Run unmigrated migrations",
	Run:   runUp,
}

var MigrateRollback = &cobra.Command{
	Use:     "down",
	Aliases: []string{"rollback"},
	Short:   "Reverse the up command",
	Run:     runDown,
}

func init() {
	Migrate.AddCommand(
		MigrateUp,
		MigrateRollback,
	)
}

func migrator() *migrate.Migrator {
	// Register all migration files under database/migrations
	migrations.Initialize()
	// Initialize migrator
	return migrate.NewMigrator()
}

func runUp(cmd *cobra.Command, args []string) {
	migrator().Up()
}

func runDown(cmd *cobra.Command, args []string) {
	migrator().Rollback()
}
