package migrations

import (
	"database/sql"

	"gohub/app/models"
	"gohub/pkg/migrate"

	"gorm.io/gorm"
)

func init() {
	type Category struct {
		models.BaseModel

		UserID string `gorm:"type:bigint;default:0;index"`
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&Category{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropColumn(&Category{}, "UserID")
	}

	migrate.Add("2026_04_23_070105_add_user_id_to_categories", up, down)
}
