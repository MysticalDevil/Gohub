package models

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"gohub/pkg/database"
	"gohub/pkg/paginator"
	"gorm.io/gorm"
)

func Query(ctx context.Context) *gorm.DB {
	return database.DBWithContext(ctx)
}

func GetBy[T any](ctx context.Context, field string, value any) (model T, err error) {
	result := database.DBWithContext(ctx).Where(field+" = ?", value).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model, nil
		}
		return model, result.Error
	}
	return model, nil
}

func Get[T any](ctx context.Context, id any) (model T, err error) {
	result := database.DBWithContext(ctx).Where("id", id).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model, nil
		}
		return model, result.Error
	}
	return model, nil
}

func All[T any](ctx context.Context) (models []T, err error) {
	result := database.DBWithContext(ctx).Find(&models)
	if result.Error != nil {
		return models, result.Error
	}
	return models, nil
}

func Exists[T any](ctx context.Context, field string, value any) (bool, error) {
	var count int64
	result := database.DBWithContext(ctx).Model(new(T)).Where(field+" = ?", value).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func Paginate[T any](ctx context.Context, c *gin.Context, limit int) (models []T, paging paginator.Paging) {
	query := database.DBWithContext(ctx).Model(new(T))
	paging = paginator.Paginate(ctx, c, query, &models, limit)
	return
}
