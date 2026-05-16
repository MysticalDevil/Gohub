package auth

import (
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/models/user"
	"gohub/app/requests"
	"gohub/pkg/response"
)

type PasswordController struct {
	v1.BaseAPIController
}

func (pc *PasswordController) ResetByPhone(c *gin.Context) {
	request := requests.ResetByPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.ResetByPhone); !ok {
		return
	}

	userModel, err := user.GetByPhone(c.Request.Context(), request.Phone)
	if err != nil {
		response.Abort500(c, "Failed to query user")
		return
	}
	if userModel.ID == 0 {
		response.Abort404(c)
		return
	}

	userModel.Password = request.Password
	userModel.Save(c.Request.Context())

	response.Success(c)
}

func (pc *PasswordController) ResetByEmail(c *gin.Context) {
	request := requests.ResetByEmailRequest{}
	if ok := requests.Validate(c, &request, requests.ResetByEmail); !ok {
		return
	}

	userModel, err := user.GetByEmail(c.Request.Context(), request.Email)
	if err != nil {
		response.Abort500(c, "Failed to query user")
		return
	}
	if userModel.ID == 0 {
		response.Abort404(c)
		return
	}

	userModel.Password = request.Password
	userModel.Save(c.Request.Context())

	response.Success(c)
}
