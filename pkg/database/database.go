// Package database Database operations
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"gohub/pkg/config"
	"gohub/pkg/logger"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	SQLDB *sql.DB
)

// Connect To connect database
func Connect(dbConfig gorm.Dialector, _logger gormLogger.Interface) {
	// Use gorm.Open to connect to the database
	var err error
	DB, err = gorm.Open(dbConfig, &gorm.Config{
		Logger: _logger,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	SQLDB, err = DB.DB()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func CurrentDatabase() (dbName string) {
	dbName = DB.Migrator().CurrentDatabase()
	return
}

func DeleteAllTables() (err error) {
	switch config.Get("database.connection") {
	case "mysql":
		err = deleteMySQLTables()
	case "sqlite":
		err = deleteAllSqliteTables()
	case "postgresql":
		err = deletePostgresQLTables()
	default:
		panic(errors.New("database connection not supported"))
	}

	return
}

func deleteMySQLTables() error {
	dbName := CurrentDatabase()
	var tables []string

	// Read all tables
	err := DB.Table("information_schema.tables").
		Where("table_schema = ?", dbName).
		Pluck("table_name", &tables).
		Error
	if err != nil {
		return err
	}

	// Disable foreign key detection temporarily
	DB.Exec("SET foreign_key_checks = 0;")

	// Delete all tables
	for _, table := range tables {
		err := DB.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}

	// Enable foreign key detection
	DB.Exec("SET foreign_key_checks = 1;")
	return nil
}

func deletePostgresQLTables() error {
	dbName := CurrentDatabase()
	var tables []string

	// Read all tables
	err := DB.Table("information_schema.tables").
		Where("table_schema = ? AND table_catalog = ?", "public", dbName).
		Pluck("table_name", &tables).
		Error
	if err != nil {
		return err
	}

	// Disable foreign key detection temporarily
	DB.Exec("SET CONSTRAINTS ALL DEFERRED;")

	// Delete all tables
	for _, table := range tables {
		err := DB.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}

	// Enable foreign key detection
	DB.Exec("SET CONSTRAINTS ALL IMMEDIATE;")
	return nil
}

func deleteAllSqliteTables() error {
	var tables []string

	// Read all tables
	err := DB.Select(&tables, "SELECT name FROM sqlite_master WHERE type='table'").Error
	if err != nil {
		return err
	}

	// Delete all tables
	for _, table := range tables {
		err := DB.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}

	return nil
}

func TableName(obj any) string {
	stmt := &gorm.Statement{DB: DB}
	err := stmt.Parse(obj)
	if err != nil {
		logger.LogIf(err)
	}
	return stmt.Schema.Table
}
