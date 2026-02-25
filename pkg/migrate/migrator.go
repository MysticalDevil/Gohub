// Package migrate Handling database migration
package migrate

import (
	"os"

	"gohub/pkg/console"
	"gohub/pkg/database"
	"gohub/pkg/file"
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
			logger.ErrorString("migrate", "error", err.Error())
		}
	}
}

// Up Execute all files that have not been migrated
func (migrator *Migrator) Up() {
	// Read all migration files to ensure they are sorted by time
	migrateFiles := migrator.readAllMigrationFiles()

	// Get the current batch value
	batch := migrator.getBatch()

	// Get all migration data
	var migrations []Migration
	migrator.DB.Find(&migrations)

	// This value can be used to determine if the database is up-to-date
	runed := false

	// Iterate over the migration file and execute the up callback if it has not been executed before
	for _, mfile := range migrateFiles {
		// Compare file names to determine if they have been run
		if mfile.isNotMigrated(migrations) {
			migrator.runUpMigration(mfile, batch)
			runed = true
		}
	}

	if !runed {
		console.Success("database is up-to-date.")
	}
}

// Rollback of the previous operation
func (migrator *Migrator) Rollback() {
	// Get the last batch of migration data
	lastMigration := Migration{}
	migrator.DB.Order("id DESC").First(&lastMigration)

	var migrations []Migration
	migrator.DB.Where("batch = ?", lastMigration.Batch).Order("id DESC").Find(&migrations)

	// Rollback of the last migration
	if !migrator.rollbackMigrations(migrations) {
		console.Success("[migrations] table is empty, nothing to rollback.")
	}
}

// Fallback migration, execute the down method of migration in reverse order
func (migrator *Migrator) rollbackMigrations(migrations []Migration) bool {
	// Flag whether a migration fallback has actually been performed
	runed := false

	for _, _migration := range migrations {
		console.Warning("rollback " + _migration.Migration)

		// Execute the down method of migrating files
		mfile := getMigrationFile(_migration.Migration)
		if mfile.Down != nil {
			mfile.Down(database.DB.Migrator(), database.SQLDB)
		}

		runed = true

		// Delete this record if the rollback is successful
		migrator.DB.Delete(&_migration)

		console.Success("finish " + mfile.FileName)
	}

	return runed
}

// Reset all migrations
func (migrator *Migrator) Reset() {
	var migrations []Migration

	// Read all migration files in reverse order
	migrator.DB.Order("id DESC").Find(&migrations)

	// Rollback of all migrations
	if !migrator.rollbackMigrations(migrations) {
		console.Success("[migrations] table is empty, nothing to reset.")
	}
}

// Refresh Roll back all migrations and run all migrations
func (migrator *Migrator) Refresh() {
	// Rollback of all migrations
	migrator.Reset()
	// Run all migrations
	migrator.Up()
}

// Fresh Drop all tables and rerun all migrations
func (migrator *Migrator) Fresh() {
	// Get the database name to prompt for
	dbName := database.CurrentDatabase()

	// Delete all tables
	err := database.DeleteAllTables()
	console.ExitIf(err)
	console.Success("clear up database " + dbName)

	// Re-create the migrates table
	migrator.createMigrationsTable()
	console.Success("[migrations] table created.")

	migrator.Up()
}

// Get the current value of this batch
func (migrator *Migrator) getBatch() int {
	// Default value is 1
	batch := 1

	// Take the last migration data executed
	lastMigration := Migration{}
	migrator.DB.Order("id DESC").First(&lastMigration)

	// Add one if there is a value
	if lastMigration.ID > 0 {
		batch = lastMigration.Batch + 1
	}
	return batch
}

// Read files from file directories to ensure correct time ordering
func (migrator *Migrator) readAllMigrationFiles() []MigrationFile {
	// Retrieve all files in the database/migrations/ directory
	// By default, the files are sorted by name
	files, err := os.ReadDir(migrator.Folder)
	console.ExitIf(err)

	var migrateFiles []MigrationFile
	for _, f := range files {
		// Remove the file suffix '.go'
		fileName := file.NameWithoutExtension(f.Name())

		// Get the [MigrationFile] object by the name of the migration file
		mfile := getMigrationFile(fileName)

		// Make sure the migration files are available,
		// then put them in the migrateFiles array
		if len(mfile.FileName) > 0 {
			migrateFiles = append(migrateFiles, mfile)
		}
	}

	// Returns the sorted [MigrationFile] array
	return migrateFiles
}

// Execute migration, execute the up method of migration
func (migrator *Migrator) runUpMigration(mfile MigrationFile, batch int) {
	// Execute sql for up block
	if mfile.Up != nil {
		console.Warning("migrating " + mfile.FileName)
		// Execute up method
		mfile.Up(database.DB.Migrator(), database.SQLDB)
		// Prompts for that file to be migrated
		console.Success("migrated " + mfile.FileName)
	}

	err := migrator.DB.Create(&Migration{
		Migration: mfile.FileName,
		Batch:     batch,
	}).Error
	console.ExitIf(err)
}
