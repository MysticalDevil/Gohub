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

func Get(ctx context.Context, idStr string) (link Link, err error) {
	return models.Get[Link](ctx, idStr)
}

func GetBy(ctx context.Context, field, value string) (link Link, err error) {
	return models.GetBy[Link](ctx, field, value)
}

func All(ctx context.Context) (links []Link, err error) {
	return models.All[Link](ctx)
}

func IsExist(ctx context.Context, field, value string) (bool, error) {
	return models.Exists[Link](ctx, field, value)
}

func Paginate(ctx context.Context, c *gin.Context, limit int) (links []Link, paging paginator.Paging) {
	return models.Paginate[Link](ctx, c, limit)
}

func AllCached(ctx context.Context) (links []Link, err error) {
	cacheKey := "links:all"
	expireTime := 120 * time.Minute
	cache.GetObject(cacheKey, &links)

	if helpers.Empty(links) {
		links, err = All(ctx)
		if err != nil {
			return nil, err
		}
		if helpers.Empty(links) {
			return links, nil
		}
		cache.Set(cacheKey, links, expireTime)
	}
	return links, nil
}
