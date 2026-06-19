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
	UserID      string `json:"user_id,omitempty"`

	models.CommonTimestampsField
}

func (category *Category) Create(ctx context.Context) error {
	return database.DBWithContext(ctx).Create(&category).Error
}

func (category *Category) Save(ctx context.Context) (rowsAffected int64, err error) {
	result := database.DBWithContext(ctx).Save(&category)
	return result.RowsAffected, result.Error
}

func (category *Category) UpdateFields(ctx context.Context, name, description string) (rowsAffected int64, err error) {
	result := database.DBWithContext(ctx).Model(&category).Updates(map[string]any{
		"name":        name,
		"description": description,
	})
	return result.RowsAffected, result.Error
}

func (category *Category) Delete(ctx context.Context) (rowsAffected int64, err error) {
	result := database.DBWithContext(ctx).Delete(&category)
	return result.RowsAffected, result.Error
}
