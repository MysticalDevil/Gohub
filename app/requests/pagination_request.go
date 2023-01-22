package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type PaginationRequest struct {
	Sort    string `valid:"sort" form:"sort"`
	Order   string `valid:"order" form:"order"`
	PerPage string `valid:"per_page" from:"per_page"`
}

func Pagination(data any, _ *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"sort":     []string{"in:id,created_at,updated_at"},
		"order":    []string{"in:asc,desc"},
		"per_page": []string{"numeric_between:2,100"},
	}

	messages := govalidator.MapData{
		"sort": []string{
			"in:Sort fields only support id, created_at, updated_at",
		},
		"order": []string{
			"in:Sort fields only support asc (positive order), desc (reverse order)",
		},
		"per_page": []string{
			"numeric_between:The number of entries per page has a value between 2 and 100",
		},
	}

	return validate(data, rules, messages)
}
