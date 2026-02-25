// Package paginator Handling paging logic
package paginator

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gohub/pkg/config"
	"gohub/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Paging Data
// Offset-based pagination
// Example: {"offset": 0, "limit": 20, "total": 200}
type Paging struct {
	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
	Total  int64 `json:"total"`
}

// Paginator Page operation class
// Offset-based pagination
// Use offset/limit in the URL query for pagination
// Example: ?offset=0&limit=20&sort=id&order=asc
type Paginator struct {
	Limit      int
	Offset     int
	TotalCount int64
	Sort       string
	Order      string

	query *gorm.DB
	ctx   *gin.Context
}

// Paginate
// c - gin.context, Used to get the URL parameters for paging
// db - Gorm, Query handle to query the data set and get the total number of data
// data - Model arrays, pass address to get data
// limit - Number of entries per page, preferably from the url parameter, otherwise use the limit value
//
// for example:
//
// query := database.DB.Model(Topic{}).Where("category_id = ?", cid)
// var topics []Topic
// paging := paginator.Paginate(
//
//	c,
//	query,
//	&topics,
//	limit,
//
// )
func Paginate(_ context.Context, c *gin.Context, db *gorm.DB, data any, limit int) Paging {
	// Initialize Paginator instance
	p := &Paginator{
		query: db,
		ctx:   c,
	}
	p.initProperties(limit)

	// Query database
	err := p.query.Preload(clause.Associations). // Read Associations
							Order(p.Sort + " " + p.Order). // Sort
							Limit(p.Limit).
							Offset(p.Offset).
							Find(data).
							Error
	if err != nil {
		logger.LogIf(err)
		return Paging{}
	}

	return Paging{
		Offset: p.Offset,
		Limit:  p.Limit,
		Total:  p.TotalCount,
	}
}

// The properties that must be used to initialize paging and query the database based on these properties
func (p *Paginator) initProperties(limit int) {
	p.Limit = p.getLimit(limit)
	p.Offset = p.getOffset()

	// Sort parameters
	p.Order = p.ctx.DefaultQuery(config.Get("paging.url_query_order"), "asc")
	p.Sort = p.ctx.DefaultQuery(config.Get("paging.url_query_sort"), "id")

	p.TotalCount = p.getTotalCount()
}

func (p *Paginator) getLimit(limit int) int {
	// Preferred use of limit parameter
	queryLimit := p.ctx.Query(config.Get("paging.url_query_limit"))
	if len(queryLimit) > 0 {
		limit = cast.ToInt(queryLimit)
	}

	// No parameters are passed, using the default
	if limit <= 0 {
		limit = config.GetInt("paging.limit")
	}

	return limit
}

func (p *Paginator) getOffset() int {
	offset := cast.ToInt(p.ctx.Query(config.Get("paging.url_query_offset")))
	if offset < 0 {
		offset = 0
	}
	return offset
}

func (p *Paginator) getTotalCount() int64 {
	var count int64
	if err := p.query.Count(&count).Error; err != nil {
		return 0
	}
	return count
}
