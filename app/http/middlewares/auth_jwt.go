package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gohub/app/models/user"
	"gohub/pkg/config"
	"gohub/pkg/jwt"
	"gohub/pkg/response"
)

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := jwt.NewJWT().ParseToken(c)

		// JWT parsing failed, an error occurred
		if err != nil {
			response.Unauthorized(c,
				fmt.Sprintf("Please view the interface certification documents related to %v",
					config.GetString("app.name"),
				),
			)
			return
		}

		// JWT parsed successfully, set user information
		var userModel user.User = user.Get(claims.UserID)
		if userModel.ID == 0 {
			response.Unauthorized(c, "Could not find corresponding user, user may have been deleted")
			return
		}

		// Store the user information in gin.Context,
		// and the subsequent auth package will get the current user data from here
		c.Set("current_user_id", userModel.GetStringID())
		c.Set("current_user_name", userModel.Name)
		c.Set("current_user", userModel)

		c.Next()
	}
}
