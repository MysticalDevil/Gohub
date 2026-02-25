package link

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gohub/app/models"
	"gohub/pkg/cache"
	"gohub/pkg/helpers"
	"gohub/pkg/paginator"
)

func Get(ctx context.Context, idStr string) (link Link) {
	return models.Get[Link](ctx, idStr)
}

func GetBy(ctx context.Context, field, value string) (link Link) {
	return models.GetBy[Link](ctx, field, value)
}

func All(ctx context.Context) (links []Link) {
	return models.All[Link](ctx)
}

func IsExist(ctx context.Context, field, value string) bool {
	return models.Exists[Link](ctx, field, value)
}

func Paginate(ctx context.Context, c *gin.Context, limit int) (links []Link, paging paginator.Paging) {
	return models.Paginate[Link](ctx, c, limit)
}

func AllCached(ctx context.Context) (links []Link) {
	// Set cache key
	cacheKey := "links:all"
	// Set expire time
	expireTime := 120 * time.Minute
	// Get data
	cache.GetObject(cacheKey, &links)

	if helpers.Empty(links) {
		links = All(ctx)
		if helpers.Empty(links) {
			return
		}
		// Set cache
		cache.Set(cacheKey, links, expireTime)
	}
	return
}
