// Package category model
package category

import (
	"context"

	"gohub/app/models"
	"gohub/pkg/database"
)

type Category struct {
	models.BaseModel

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`

	models.CommonTimestampsField
}

func (category *Category) Create(ctx context.Context) {
	database.DBWithContext(ctx).Create(&category)
}

func (category *Category) Save(ctx context.Context) (rowsAffected int64) {
	result := database.DBWithContext(ctx).Save(&category)
	return result.RowsAffected
}

func (category *Category) Delete(ctx context.Context) (rowsAffected int64) {
	result := database.DBWithContext(ctx).Delete(&category)
	return result.RowsAffected
}
