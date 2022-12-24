// Package migrate Handling database migration
package migrate

import (
	"gohub/pkg/database"
	"gohub/pkg/logger"
	"gorm.io/gorm"
)

// Migrator Data migration operation class
type Migrator struct {
	Folder   string
	DB       *gorm.DB
	Migrator gorm.Migrator
}

// Migration A data in the migrations table of the corresponding data
type Migration struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement;"`
	Migration string `gorm:"type:varchar(255);not null;unique;"`
	Batch     int
}

// NewMigrator Create Migrator instances to perform migration operations
func NewMigrator() *Migrator {
	migrator := &Migrator{
		Folder:   "database/migrations/",
		DB:       database.DB,
		Migrator: database.DB.Migrator(),
	}
	// If migrations does not exist, create it
	migrator.createMigrationsTable()

	return migrator
}

// Create migrations table
func (migrator *Migrator) createMigrationsTable() {
	migration := Migration{}

	if !migrator.Migrator.HasTable(&migration) {
		err := migrator.Migrator.CreateTable(&migration)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}
