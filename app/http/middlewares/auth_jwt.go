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
		if err != nil {
			response.Unauthorized(c,
				fmt.Sprintf("Please view the interface certification documents related to %v",
					config.GetString("app.name"),
				),
			)
			return
		}

		userModel, err := user.Get(c.Request.Context(), claims.UserID)
		if err != nil {
			response.Abort500(c, "Failed to query user")
			return
		}
		if userModel.ID == 0 {
			response.Unauthorized(c, "Could not find corresponding user, user may have been deleted")
			return
		}

		c.Set("current_user_id", userModel.GetStringID())
		c.Set("current_user_name", userModel.Name)
		c.Set("current_user", userModel)

		c.Next()
	}
}
