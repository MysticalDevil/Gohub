package migrations

import (
	"database/sql"
	"gohub/pkg/migrate"

	"gorm.io/gorm"
)

func init() {
	type User struct {
		City         string `gorm:"type:varchar(10);"`
		Introduction string `gorm:"type:varchar(255);"`
		Avatar       string `gorm:"type:varchar(255);default:null"`
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&User{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropColumn(&User{}, "City")
		_ = migrator.DropColumn(&User{}, "Introduction")
		_ = migrator.DropColumn(&User{}, "Avatar")
	}

	migrate.Add("2023_01_16_160250_add_fields_to_user", up, down)
}
