package migrate

import (
	"database/sql"
	"slices"

	"gorm.io/gorm"
)

// migrationFunc Define the type of up and down callback methods
type migrationFunc func(migrator gorm.Migrator, db *sql.DB)

// migrationFiles Array of all migration files
var migrationFiles []MigrationFile

// MigrationFile Represents a single migration file
type MigrationFile struct {
	Up       migrationFunc
	Down     migrationFunc
	FileName string
}

// Add a new migration file, all migration files need this method to register
func Add(name string, up migrationFunc, down migrationFunc) {
	migrationFiles = append(migrationFiles, MigrationFile{
		FileName: name,
		Up:       up,
		Down:     down,
	})
}

// getMigrationFile Get the MigrationFile object by the name of the migration file
func getMigrationFile(name string) MigrationFile {
	if idx := slices.IndexFunc(migrationFiles, func(mfile MigrationFile) bool {
		return mfile.FileName == name
	}); idx >= 0 {
		return migrationFiles[idx]
	}
	return MigrationFile{}
}

// isNotMigrated Determine if the migration has been executed
func (mfile MigrationFile) isNotMigrated(migrations []Migration) bool {
	return !slices.ContainsFunc(migrations, func(migration Migration) bool {
		return migration.Migration == mfile.FileName
	})
}
