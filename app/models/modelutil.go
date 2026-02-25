package models

import (
	"context"

	"github.com/gin-gonic/gin"
	"gohub/pkg/database"
	"gohub/pkg/paginator"
	"gorm.io/gorm"
)

func Query(ctx context.Context) *gorm.DB {
	return database.DBWithContext(ctx)
}

func GetBy[T any](ctx context.Context, field string, value any) (model T) {
	database.DBWithContext(ctx).Where(field+" = ?", value).First(&model)
	return
}

func Get[T any](ctx context.Context, id any) (model T) {
	database.DBWithContext(ctx).Where("id", id).First(&model)
	return
}

func All[T any](ctx context.Context) (models []T) {
	database.DBWithContext(ctx).Find(&models)
	return
}

func Exists[T any](ctx context.Context, field string, value any) bool {
	var count int64
	database.DBWithContext(ctx).Model(new(T)).Where(field+" = ?", value).Count(&count)
	return count > 0
}

func Paginate[T any](ctx context.Context, c *gin.Context, limit int) (models []T, paging paginator.Paging) {
	query := database.DBWithContext(ctx).Model(new(T))
	paging = paginator.Paginate(ctx, c, query, &models, limit)
	return
}
