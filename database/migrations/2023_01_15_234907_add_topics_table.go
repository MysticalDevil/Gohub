package migrations

import (
	"database/sql"
	"gohub/app/models"
	"gohub/pkg/logger"
	"gohub/pkg/migrate"

	"gorm.io/gorm"
)

func init() {
	type User struct {
		models.BaseModel
	}

	type Category struct {
		models.BaseModel
	}

	type Topic struct {
		models.BaseModel

		Title      string `gorm:"type:varchar(255);not null;index"`
		Body       string `gorm:"type:text;not null"`
		UserID     string `gorm:"type:bigint;not null;index"`
		CategoryID string `gorm:"type:bigint;not null;index"`

		User     User
		Category Category

		models.CommonTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		if err := migrator.AutoMigrate(&Topic{}); err != nil {
			logger.LogIf(err)
			return
		}
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		migrator.DropTable(&Topic{})
	}

	migrate.Add("2023_01_15_234907_add_topics_table", up, down)
}
