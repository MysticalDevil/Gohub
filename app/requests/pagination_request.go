package requests

import (
	"github.com/gin-gonic/gin"
)

type PaginationRequest struct {
	Sort   string `valid:"sort" form:"sort"`
	Order  string `valid:"order" form:"order"`
	Offset string `valid:"offset" form:"offset"`
	Limit  string `valid:"limit" form:"limit"`
}

func Pagination(data any, c *gin.Context) map[string][]string {
	rules := MapData{
		"sort":   []string{"in:id,created_at,updated_at"},
		"order":  []string{"in:asc,desc"},
		"offset": []string{"numeric_between:0,1000000"},
		"limit":  []string{"numeric_between:1,100"},
	}

	messages := MapData{
		"sort": []string{
			"in:Sort fields only support id, created_at, updated_at",
		},
		"order": []string{
			"in:Sort fields only support asc (positive order), desc (reverse order)",
		},
		"offset": []string{
			"numeric_between:Offset must be between 0 and 1,000,000",
		},
		"limit": []string{
			"numeric_between:Limit must be between 1 and 100",
		},
	}

	return validate(c, data, rules, messages)
}
