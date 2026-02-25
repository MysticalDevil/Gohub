// Package auth Handle user authentication related logic
package auth

import (
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/models/user"
	"gohub/app/requests"
	"gohub/pkg/jwt"
	"gohub/pkg/response"
)

// SignupController Register controller
type SignupController struct {
	v1.BaseAPIController
}

// IsPhoneExist Check if the phone number is registered
func (sc *SignupController) IsPhoneExist(c *gin.Context) {
	request := requests.SignupPhoneExistRequest{}
	// Parse Json request
	if ok := requests.Validate(c, &request, requests.SignupPhoneExist); !ok {
		return
	}

	// Check database and return response
	response.Data(c, gin.H{
		"exist": user.IsPhoneExist(c.Request.Context(), request.Phone),
	})
}

// IsEmailExist  Check if the email is registered
func (sc *SignupController) IsEmailExist(c *gin.Context) {
	request := requests.SignupEmailExistRequest{}
	if ok := requests.Validate(c, &request, requests.SignupEmailExist); !ok {
		return
	}

	response.Data(c, gin.H{
		"exist": user.IsEmailExist(c.Request.Context(), request.Email),
	})
}

// SignupUsingPhone Sign up with phone
func (sc *SignupController) SignupUsingPhone(c *gin.Context) {
	request := requests.SignupUsingPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.SignupUsingPhone); !ok {
		return
	}

	userModel := user.User{
		Name:     request.Name,
		Phone:    request.Phone,
		Password: request.Password,
	}
	userModel.Create(c.Request.Context())

	if userModel.ID > 0 {
		token := jwt.NewJWT().IssueToken(userModel.GetStringID(), userModel.Name)
		response.Created(c, gin.H{
			"token": token,
			"user":  userModel,
		})
	} else {
		response.Abort500(c, "Failed to create user, please try later~")
	}
}

// SignupUsingEmail Sign up with email
func (sc *SignupController) SignupUsingEmail(c *gin.Context) {
	request := requests.SignupUsingEmailRequest{}
	if ok := requests.Validate(c, &request, requests.SignupUsingEmail); !ok {
		return
	}

	userModel := user.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}
	userModel.Create(c.Request.Context())

	if userModel.ID > 0 {
		token := jwt.NewJWT().IssueToken(userModel.GetStringID(), userModel.Name)
		response.Created(c, gin.H{
			"token": token,
			"user":  userModel,
		})
	} else {
		response.Abort500(c, "Failed to create user, please try later~")
	}
}
