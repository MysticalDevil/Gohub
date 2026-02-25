// Package paginator Handling paging logic
package paginator

import (
	"fmt"
	"math"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gohub/pkg/config"
	"gohub/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Paging Data
type Paging struct {
	CurrentPage int    // Current page
	PerPage     int    // Number of entries per page
	TotalPage   int    // Total page
	TotalCount  int64  // Total number of entries
	NextPageURL string // Links on the next page
	PrevPageURL string // Links on the previous page
}

// Paginator Page operation class
type Paginator struct {
	BaseURL    string // Used to splice URLs
	PerPage    int    // Number of entries per page
	Page       int    // Current page
	Offset     int    // The value of offset when the database reads data
	TotalCount int64  // Total number of entries
	TotalPage  int    // Total page = TotalCount/PerPage
	Sort       string // Sorting rules
	Order      string // Sort order

	query *gorm.DB // db query handle
	ctx   *gin.Context
}

// Paginate
// c - gin.context, Used to get the URL parameters for paging
// db - Gorm, Query handle to query the data set and get the total number of data
// baseURL - Use to paginate links
// data - Model arrays, pass address to get data
// PerPage - Number of entries per page, preferably from the url parameter, otherwise use the perPage value
//
// for example:
//
// query := database.DB.Model(Topic{}).Where("category_id = ?", cid)
// var topics []Topic
// paging := paginator.Paginate(
//
//			c,
//	     query,
//	     &topics,
//	     app.APIURL(database.TableName(&Topic{})),
//	     perPage,
//
// )
func Paginate(c *gin.Context, db *gorm.DB, data any, baseURL string, perPage int) Paging {
	// Initialize Paginator instance
	p := &Paginator{
		query: db,
		ctx:   c,
	}
	p.initProperties(perPage, baseURL)

	// Query database
	err := p.query.Preload(clause.Associations). // Read Associations
							Order(p.Sort + " " + p.Order). // Sort
							Limit(p.PerPage).
							Offset(p.Offset).
							Find(data).
							Error
	if err != nil {
		logger.LogIf(err)
		return Paging{}
	}

	return Paging{
		CurrentPage: p.Page,
		PerPage:     p.PerPage,
		TotalPage:   p.TotalPage,
		TotalCount:  p.TotalCount,
		NextPageURL: p.getNextPageURL(),
		PrevPageURL: p.getPrevPageURL(),
	}
}

// The properties that must be used to initialize paging and query the database based on these properties
func (p *Paginator) initProperties(perPage int, baseURL string) {
	p.BaseURL = p.formatBaseURL(baseURL)
	p.PerPage = p.getPerPage(perPage)

	// Sort parameters
	p.Order = p.ctx.DefaultQuery(config.Get("paging.url_query_order"), "asc")
	p.Sort = p.ctx.DefaultQuery(config.Get("paging.url_query_sort"), "id")

	p.TotalCount = p.getTotalCount()
	p.TotalPage = p.getTotalPage()
	p.Page = p.getCurrentPage()
	p.Offset = (p.Page - 1) * p.PerPage
}

func (p *Paginator) getPerPage(perPage int) int {
	// Preferred use of per_page parameter
	queryPerPage := p.ctx.Query(config.Get("paging.url_query_per_page"))
	if len(queryPerPage) > 0 {
		perPage = cast.ToInt(queryPerPage)
	}

	// No parameters are passed, using the default
	if perPage <= 0 {
		perPage = config.GetInt("paging.perpage")
	}

	return perPage
}

func (p *Paginator) getCurrentPage() int {
	// Preferred page for user requests
	page := cast.ToInt(p.ctx.Query(config.Get("paging.url_query_page")))
	if page <= 0 {
		page = 1
	}
	// TotalPage is equal to 0, which means there is not enough data for paging
	if p.TotalPage == 0 {
		return 0
	}
	// If the number of requested pages is greater than the total number of pages,
	// return the total number of pages
	if page > p.TotalPage {
		return p.TotalPage
	}

	return page
}

func (p *Paginator) getTotalCount() int64 {
	var count int64
	if err := p.query.Count(&count).Error; err != nil {
		return 0
	}
	return count
}

func (p *Paginator) getTotalPage() int {
	if p.TotalCount == 0 {
		return 0
	}

	nums := int64(math.Ceil(float64(p.TotalCount) / float64(p.PerPage)))
	if nums == 0 {
		nums = 1
	}
	return int(nums)
}

// Compatible URLs with and without `?` case
func (p *Paginator) formatBaseURL(baseURL string) string {
	if strings.Contains(baseURL, "?") {
		baseURL = baseURL + "&" + config.Get("paging.url_query_page") + "="
	} else {
		baseURL = baseURL + "?" + config.Get("paging.url_query_page") + "="
	}
	return baseURL
}

// Splice pagination links
func (p *Paginator) getPageLink(page int) string {
	return fmt.Sprintf("%v%v&%s=%s&%s=%s&%s=%v",
		p.BaseURL,
		page,
		config.Get("paging.url_query_sort"),
		p.Sort,
		config.Get("paging.url_query_order"),
		p.Order,
		config.Get("paging.url_query_per_page"),
		p.PerPage,
	)
}

func (p *Paginator) getNextPageURL() string {
	if p.TotalPage > p.Page {
		return p.getPageLink(p.Page + 1)
	}
	return ""
}

func (p *Paginator) getPrevPageURL() string {
	if p.Page <= 1 || p.Page > p.TotalPage {
		return ""
	}
	return p.getPageLink(p.Page - 1)
}
