package seeders

import (
	"fmt"
	"gohub/database/factories"
	"gohub/pkg/console"
	"gohub/pkg/logger"
	"gohub/pkg/seed"
	"gorm.io/gorm"
)

func init() {
	// Add Seeder
	seed.Add("SeedUsersTable", func(db *gorm.DB) {
		// Create 10 user objects
		users := factories.MakeUsers(10)

		// Create users in bulk (note that bulk creation does not invoke model hooks)
		result := db.Table("users").Create(&users)

		// Record error
		if err := result.Error; err != nil {
			logger.LogIf(err)
			return
		}

		// Print runs
		console.Success(
			fmt.Sprintf(
				"Table [%v] %v rows seeded",
				result.Statement.Table,
				result.RowsAffected,
			),
		)
	})
}
