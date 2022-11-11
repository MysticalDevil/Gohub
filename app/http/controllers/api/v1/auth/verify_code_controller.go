package auth

import (
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/requests"
	"gohub/pkg/captcha"
	"gohub/pkg/logger"
	"gohub/pkg/response"
	"gohub/pkg/verifycode"
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
