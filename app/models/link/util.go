package link

import (
	"time"

	"github.com/gin-gonic/gin"
	"gohub/pkg/app"
	"gohub/pkg/cache"
	"gohub/pkg/database"
	"gohub/pkg/helpers"
	"gohub/pkg/paginator"
)

func Get(idStr string) (link Link) {
	database.DB.Where("id", idStr).First(&link)
	return
}

func GetBy(field, value string) (link Link) {
	database.DB.Where("? = ?", field, value).First(&link)
	return
}

func All() (links []Link) {
	database.DB.Find(&links)
	return
}

func IsExist(field, value string) bool {
	var count int64
	database.DB.Model(Link{}).Where("? = ?", field, value).Count(&count)
	return count > 0
}

func Paginate(c *gin.Context, perPage int) (links []Link, paging paginator.Paging) {
	paging = paginator.Paginate(
		c,
		database.DB.Model(Link{}),
		&links,
		app.V1URL(database.TableName(&Link{})),
		perPage,
	)
	return
}

func AllCached() (links []Link) {
	// Set cache key
	cacheKey := "links:all"
	// Set expire time
	expireTime := 120 * time.Minute
	// Get data
	cache.GetObject(cacheKey, &links)

	if helpers.Empty(links) {
		links = All()
		if helpers.Empty(links) {
			return
		}
		// Set cache
		cache.Set(cacheKey, links, expireTime)
	}
	return
}
