// Package routes Register route
package routes

import (
	"github.com/gin-gonic/gin"
	controllers "gohub/app/http/controllers/api/v1"
	"gohub/app/http/controllers/api/v1/auth"
	"gohub/app/http/middlewares"
)

// RegisterAPIRoutes Registration page related routing
func RegisterAPIRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	// Global middleware: rate limit per hour. Here is where all API requests add up.
	v1.Use(middlewares.LimitIP("200-H"))
	{
		authGroup := v1.Group("/auth")
		authGroup.Use(middlewares.LimitIP("1000-H"))
		{
			// Sign up
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

			// Login
			lgc := new(auth.LoginController)
			// Use phone, SMS verify code to login
			authGroup.POST("/login/using-phone", lgc.LoginByPhone)
			// Support phone, username, email
			authGroup.POST("/login/using-password", lgc.LoginByPassword)
			authGroup.POST("/login/refresh-token", lgc.RefreshToken)

			// Reset password
			pwc := new(auth.PasswordController)
			authGroup.POST("/password-reset/using-phone", pwc.ResetByPhone)
			authGroup.POST("/password-reset/using-email", pwc.ResetByEmail)
		}
	}

	uc := new(controllers.UsersController)
	// Get current user
	v1.GET("/user", middlewares.AuthJWT(), uc.CurrentUser)
	userGroup := v1.Group("/users")
	{
		userGroup.GET("", uc.Index)
	}

	cgc := new(controllers.CategoriesController)
	cgcGroup := v1.Group("/categories")
	{
		cgcGroup.GET("", cgc.Index)
		cgcGroup.POST("", middlewares.AuthJWT(), cgc.Store)
		cgcGroup.PUT("/:id", middlewares.AuthJWT(), cgc.Update)
		cgcGroup.DELETE("/:id", middlewares.AuthJWT(), cgc.Delete)
	}

	tpc := new(controllers.TopicsController)
	tpcGroup := v1.Group("/topics")
	{
		tpcGroup.POST("", middlewares.AuthJWT(), tpc.Store)
		tpcGroup.PUT("/:id", middlewares.AuthJWT(), tpc.Update)
	}
}
