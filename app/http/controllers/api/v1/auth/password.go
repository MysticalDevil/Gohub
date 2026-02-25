package auth

import (
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/models/user"
	"gohub/app/requests"
	"gohub/pkg/response"
)

// PasswordController User controller
type PasswordController struct {
	v1.BaseAPIController
}

// ResetByPhone Use phone and verify code to reset password
func (pc *PasswordController) ResetByPhone(c *gin.Context) {
	request := requests.ResetByPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.ResetByPhone); !ok {
		return
	}

	userModel := user.GetByPhone(c.Request.Context(), request.Phone)
	if userModel.ID == 0 {
		response.Abort404(c)
	} else {
		userModel.Password = request.Password
		userModel.Save(c.Request.Context())

		response.Success(c)
	}
}

// ResetByEmail Use email and verify code to reset password
func (pc *PasswordController) ResetByEmail(c *gin.Context) {
	request := requests.ResetByEmailRequest{}
	if ok := requests.Validate(c, &request, requests.ResetByEmail); !ok {
		return
	}

	userModel := user.GetByEmail(c.Request.Context(), request.Email)
	if userModel.ID == 0 {
		response.Abort404(c)
	} else {
		userModel.Password = request.Password
		userModel.Save(c.Request.Context())

		response.Success(c)
	}
}
