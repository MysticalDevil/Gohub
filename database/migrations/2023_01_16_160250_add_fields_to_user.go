package migrations

import (
	"database/sql"
	"errors"

	"gohub/pkg/migrate"

	"gorm.io/gorm"
)

func init() {
	type User struct {
		City         string `gorm:"type:varchar(10);"`
		Introduction string `gorm:"type:varchar(255);"`
		Avatar       string `gorm:"type:varchar(255);default:null"`
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) error {
		return migrator.AutoMigrate(&User{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) error {
		return errors.Join(
			migrator.DropColumn(&User{}, "City"),
			migrator.DropColumn(&User{}, "Introduction"),
			migrator.DropColumn(&User{}, "Avatar"),
		)
	}

	migrate.Add("2023_01_16_160250_add_fields_to_user", up, down)
}
