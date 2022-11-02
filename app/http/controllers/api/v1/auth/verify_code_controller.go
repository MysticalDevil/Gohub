package auth

import (
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/pkg/captcha"
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
