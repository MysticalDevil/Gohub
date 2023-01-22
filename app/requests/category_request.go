package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type CategoryRequest struct {
	Name        string `valid:"name" json:"name"`
	Description string `valid:"description" json:"description,omitempty"`
}

func CategorySave(data any, _ *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"name":        []string{"required", "min_cn:2", "max_cn:8", "not_exists:categories,name"},
		"description": []string{"min_cn:3", "max_cn:255"},
	}
	messages := govalidator.MapData{
		"name": []string{
			"required:Name is required",
			"min_cn:Name length should be at least 2 words",
			"max_cn:Name length cannot exceed 8 words",
			"not_exists:Name already exists",
		},
		"description": []string{
			"min_cn:Description length should be at least 3 words",
			"max_cn:Description length cannot exceed 255 characters",
		},
	}

	return validate(data, rules, messages)
}
