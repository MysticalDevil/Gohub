// Package topic model
package topic

import (
	"context"

	"gohub/app/models"
	"gohub/app/models/category"
	"gohub/app/models/user"
	"gohub/pkg/database"
)

type Topic struct {
	models.BaseModel

	Title      string `json:"title,omitempty"`
	Body       string `json:"body,omitempty"`
	UserID     string `json:"user_id,omitempty"`
	CategoryID string `json:"category_id,omitempty"`

	// Associate users by user_id
	User user.User `json:"user"`

	// Associate categories by category_id
	Category category.Category `json:"category"`

	models.CommonTimestampsField
}

func (topic *Topic) Create(ctx context.Context) {
	database.DBWithContext(ctx).Create(&topic)
}

func (topic *Topic) Save(ctx context.Context) (rowsAffected int64) {
	result := database.DBWithContext(ctx).Save(&topic)
	return result.RowsAffected
}

func (topic *Topic) Delete(ctx context.Context) (rowsAffected int64) {
	result := database.DBWithContext(ctx).Delete(&topic)
	return result.RowsAffected
}
