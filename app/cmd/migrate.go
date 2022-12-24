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

func init() {
	Migrate.AddCommand(
		MigrateUp,
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
