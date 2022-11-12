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
			authGroup.POST("/signup/using-phone", suc.SignupUsingPhone)
			authGroup.POST("/signup/using-email", suc.SignupUsingEmail)
			// Send verification code
			vcc := new(auth.VerifyController)
			// Image verification code, need to add current limit
			authGroup.POST("/verify-codes/captcha", vcc.ShowCaptcha)
			authGroup.POST("/verify-codes/phone", vcc.SendUsingPhone)
			authGroup.POST("/verify-codes/email", vcc.SendUsingEmail)
			lgc := new(auth.LoginController)
			// Use phone, SMS verify code to login
			authGroup.POST("/login/using-phone", lgc.LoginByPhone)
		}
	}
}
