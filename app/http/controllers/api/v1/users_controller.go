package v1

import (
	"github.com/gin-gonic/gin"
	"gohub/app/models/user"
	"gohub/pkg/auth"
	"gohub/pkg/response"
)

type UsersController struct {
	BaseAPIController
}

// CurrentUser Information about the currently logged-in user
func (ctrl *UsersController) CurrentUser(c *gin.Context) {
	userModel := auth.CurrentUser(c)
	response.Data(c, userModel)
}

// Index All user
func (ctrl *UsersController) Index(c *gin.Context) {
	data := user.All()
	response.Data(c, data)
}
