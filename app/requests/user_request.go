package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"gohub/pkg/auth"
)

type UserUpdateProfileRequest struct {
	Name         string `valid:"name" json:"name"`
	City         string `valid:"city" json:"city"`
	Introduction string `valid:"introduction" json:"introduction"`
}

func UserUpdateProfile(data any, c *gin.Context) map[string][]string {
	uid := auth.CurrentUID(c)
	rules := govalidator.MapData{
		"name":         []string{"required", "alpha_num", "between:3,20", "not_exists:users,name," + uid},
		"introduction": []string{"min_cn:4", "max_cn:240"},
		"city":         []string{"min_cn:2", "max_cn:20"},
	}

	messages := govalidator.MapData{
		"name": []string{
			"required:Username is required",
			"alpha_num:Username is malformed, only numbers and English are allowed",
			"between:Username length must be between 3 and 20",
			"not_exists:Username already taken",
		},
		"introduction": []string{
			"min_cn:The introduction length must be at least 4 characters",
			"max_cn:The introduction length must be at most 240 characters",
		},
		"city": []string{
			"min_cn:The city length must be at least 2 characters",
			"max_cn:The city length must be at most 20 characters",
		},
	}

	return validate(data, rules, messages)
}
