package topic

import (
	"context"

	"github.com/gin-gonic/gin"
	"gohub/app/models"
	"gohub/pkg/database"
	"gohub/pkg/paginator"
	"gorm.io/gorm/clause"
)

func Get(ctx context.Context, idStr string) (topic Topic) {
	database.DBWithContext(ctx).Preload(clause.Associations).Where("id", idStr).First(&topic)
	return
}

func GetBy(ctx context.Context, field, value string) (topic Topic) {
	return models.GetBy[Topic](ctx, field, value)
}

func All(ctx context.Context) (topics []Topic) {
	return models.All[Topic](ctx)
}

func IsExist(ctx context.Context, field, value string) bool {
	return models.Exists[Topic](ctx, field, value)
}

func Paginate(ctx context.Context, c *gin.Context, limit int) (topics []Topic, paging paginator.Paging) {
	return models.Paginate[Topic](ctx, c, limit)
}
