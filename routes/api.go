// Package routes Register route
package routes

import (
	"github.com/gin-gonic/gin"
	"gohub/app/http/controllers/api/v1/auth"
)

// RegisterAPIRoutes Registration page related routing
func RegisterAPIRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			suc := new(auth.SignupController)
			// Determine whether the phone number is registered
			authGroup.POST("/signup/phone/exist", suc.IsPhoneExist)
			// Determine whether the email is registered
			authGroup.POST("/signup/email/exist", suc.IsEmailExist)
			// Send verification code
			vcc := new(auth.VerifyController)
			// Image verification code, need to add current limit
			authGroup.GET("/verify-codes/captcha", vcc.ShowCaptcha)
		}
	}
}
