package auth

import (
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/models/user"
	"gohub/app/requests"
	"gohub/pkg/jwt"
	"gohub/pkg/response"
)

type SignupController struct {
	v1.BaseAPIController
}

func (sc *SignupController) IsPhoneExist(c *gin.Context) {
	request := requests.SignupPhoneExistRequest{}
	if ok := requests.Validate(c, &request, requests.SignupPhoneExist); !ok {
		return
	}

	exist, err := user.IsPhoneExist(c.Request.Context(), request.Phone)
	if err != nil {
		response.Abort500(c, "Failed to check phone number")
		return
	}
	response.Data(c, gin.H{
		"exist": exist,
	})
}

func (sc *SignupController) IsEmailExist(c *gin.Context) {
	request := requests.SignupEmailExistRequest{}
	if ok := requests.Validate(c, &request, requests.SignupEmailExist); !ok {
		return
	}

	exist, err := user.IsEmailExist(c.Request.Context(), request.Email)
	if err != nil {
		response.Abort500(c, "Failed to check email")
		return
	}
	response.Data(c, gin.H{
		"exist": exist,
	})
}

func (sc *SignupController) SignupUsingPhone(c *gin.Context) {
	request := requests.SignupUsingPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.SignupUsingPhone); !ok {
		return
	}

	sc.createUserAndRespond(c, user.User{
		Name:     request.Name,
		Phone:    request.Phone,
		Password: request.Password,
	})
}

func (sc *SignupController) SignupUsingEmail(c *gin.Context) {
	request := requests.SignupUsingEmailRequest{}
	if ok := requests.Validate(c, &request, requests.SignupUsingEmail); !ok {
		return
	}

	sc.createUserAndRespond(c, user.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
}

func (sc *SignupController) createUserAndRespond(c *gin.Context, userModel user.User) {
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
