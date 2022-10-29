package bootstrap

import (
	"errors"
	"fmt"
	"gohub/pkg/config"
	"gohub/pkg/database"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

// SetupDB Initialize database and ORM
func SetupDB() {
	var dbConfig gorm.Dialector
	switch config.Get("database.connection") {
	case "mysql":
		dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
			config.Get("database.mysql.username"),
			config.Get("database.mysql.password"),
			config.Get("database.mysql.host"),
			config.Get("database.mysql.port"),
			config.Get("database.mysql.database"),
			config.Get("database.mysql.charset"),
		)
		dbConfig = mysql.New(mysql.Config{
			DSN: dsn,
		})
	case "postgresql":
		dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=%v",
			config.Get("database.postgresql.host"),
			config.Get("database.postgresql.username"),
			config.Get("database.postgresql.password"),
			config.Get("database.postgresql.database"),
			config.Get("database.postgresql.port"),
			config.Get("database.postgresql.timezone"),
		)
		log.Println(dsn)
		dbConfig = postgres.New(postgres.Config{
			DSN: dsn,
		})
	case "sqlite":
		_database := config.Get("database.sqlite.database")
		dbConfig = sqlite.Open(_database)
	default:
		panic(errors.New("database connection not supported"))
	}
	// Connect to the database and set the log mode of Gorm
	database.Connect(dbConfig, logger.Default.LogMode(logger.Info))

	database.SQLDB.SetMaxOpenConns(config.GetInt("database.max_open_connections"))
	database.SQLDB.SetMaxIdleConns(config.GetInt("database.max_idle_connections"))
	database.SQLDB.SetConnMaxLifetime(
		time.Duration(config.GetInt("database.max_life_seconds")) * time.Second,
	)
}
