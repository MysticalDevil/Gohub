package category

import (
	"context"

	"github.com/gin-gonic/gin"
	"gohub/app/models"
	"gohub/pkg/paginator"
)

func Get(ctx context.Context, idStr string) (category Category) {
	return models.Get[Category](ctx, idStr)
}

func GetBy(ctx context.Context, field, value string) (category Category) {
	return models.GetBy[Category](ctx, field, value)
}

func All(ctx context.Context) (categories []Category) {
	return models.All[Category](ctx)
}

func IsExist(ctx context.Context, field, value string) bool {
	return models.Exists[Category](ctx, field, value)
}

func Paginate(ctx context.Context, c *gin.Context, limit int) (categories []Category, paging paginator.Paging) {
	return models.Paginate[Category](ctx, c, limit)
}
