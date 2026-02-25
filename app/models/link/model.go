// Package link model
package link

import (
	"context"

	"gohub/app/models"
	"gohub/pkg/database"
)

type Link struct {
	models.BaseModel

	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`

	models.CommonTimestampsField
}

func (link *Link) Create(ctx context.Context) {
	database.DBWithContext(ctx).Create(&link)
}

func (link *Link) Save(ctx context.Context) (rowsAffected int64) {
	result := database.DBWithContext(ctx).Save(&link)
	return result.RowsAffected
}

func (link *Link) Delete(ctx context.Context) (rowsAffected int64) {
	result := database.DBWithContext(ctx).Delete(&link)
	return result.RowsAffected
}
