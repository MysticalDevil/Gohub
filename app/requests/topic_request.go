package requests

import (
	"github.com/gin-gonic/gin"
)

type TopicRequest struct {
	Title      string `json:"title,omitempty" valid:"title"`
	Body       string `json:"body,omitempty" valid:"body"`
	CategoryID string `json:"category_id,omitempty" valid:"category_id"`
}

func TopicSave(data any, c *gin.Context) map[string][]string {
	rules := MapData{
		"title":       []string{"required", "min_cn:3", "max_cn:40"},
		"body":        []string{"required", "min_cn:10", "max_cn:50000"},
		"category_id": []string{"required", "exists:categories,id"},
	}
	messages := MapData{
		"title": []string{
			"required:Topic title is required",
			"min_cn:Title length must be greater than 3",
			"max_cn:Title length must be less than 40",
		},
		"body": []string{
			"required:Topic body is required",
			"min_cn:Body length must be greater than 10",
		},
		"category_id": []string{
			"required:Topic category is required",
			"exists:Category not found",
		},
	}

	return validate(c, data, rules, messages)
}
