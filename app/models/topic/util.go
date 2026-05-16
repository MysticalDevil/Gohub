package topic

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"gohub/app/models"
	"gohub/pkg/database"
	"gohub/pkg/paginator"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Get(ctx context.Context, idStr string) (topic Topic, err error) {
	result := database.DBWithContext(ctx).Preload(clause.Associations).Where("id", idStr).First(&topic)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return topic, nil
		}
		return topic, result.Error
	}
	return topic, nil
}

func GetBy(ctx context.Context, field, value string) (topic Topic, err error) {
	return models.GetBy[Topic](ctx, field, value)
}

func All(ctx context.Context) (topics []Topic, err error) {
	return models.All[Topic](ctx)
}

func IsExist(ctx context.Context, field, value string) (bool, error) {
	return models.Exists[Topic](ctx, field, value)
}

func Paginate(ctx context.Context, c *gin.Context, limit int) (topics []Topic, paging paginator.Paging) {
	return models.Paginate[Topic](ctx, c, limit)
}
