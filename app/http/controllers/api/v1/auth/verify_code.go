package auth

import (
	"github.com/gin-gonic/gin"
	"gohub/app/captcha"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/requests"
	"gohub/app/verifycode"
	"gohub/pkg/logger"
	"gohub/pkg/response"
)

// VerifyController  User controller
type VerifyController struct {
	v1.BaseAPIController
}

// ShowCaptcha Show image verification code
func (vc *VerifyController) ShowCaptcha(c *gin.Context) {
	// Generate verification code
	id, b64s, err := captcha.NewCaptcha().GenerateCaptcha()
	logger.LogIf(err)
	response.JSON(c, gin.H{
		"captcha_id":    id,
		"captcha_image": b64s,
	})
}

// SendUsingPhone Send phone verify code
func (vc *VerifyController) SendUsingPhone(c *gin.Context) {
	// Validate the form
	request := requests.VerifyCodePhoneRequest{}
	if ok := requests.Validate(c, &request, requests.VerifyCodePhone); !ok {
		return
	}

	// Send SMS
	if ok := verifycode.NewVerifyCode().SendSMS(request.Phone); !ok {
		response.Abort500(c, "Failed to send SMS~")
	} else {
		response.Success(c)
	}
}

// SendUsingEmail Send email verify code
func (vc *VerifyController) SendUsingEmail(c *gin.Context) {
	request := requests.VerifyCodeEmailRequest{}
	if ok := requests.Validate(c, &request, requests.VerifyCodeEmail); !ok {
		return
	}

	err := verifycode.NewVerifyCode().SendEmail(request.Email)
	if err != nil {
		response.Abort500(c, "Failed to send email verification code~")
	} else {
		response.Success(c)
	}
}
