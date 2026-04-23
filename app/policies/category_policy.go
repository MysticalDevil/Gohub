// Package policies User authorization
package policies

import (
	"github.com/gin-gonic/gin"
	"gohub/app/models/category"
	"gohub/pkg/auth"
)

func CanModifyCategory(c *gin.Context, _category category.Category) bool {
	if _category.UserID == "" || _category.UserID == "0" {
		return true
	}
	return auth.CurrentUID(c) == _category.UserID
}
