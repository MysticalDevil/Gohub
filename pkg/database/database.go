// Package database Database operations
package database

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB
var SQLDB *sql.DB

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
