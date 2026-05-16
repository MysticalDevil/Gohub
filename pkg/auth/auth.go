package auth

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"gohub/app/models/user"
	"gohub/pkg/logger"
)

func Attempt(ctx context.Context, email, password string) (user.User, error) {
	userModel, err := user.GetByUtil(ctx, email)
	if err != nil {
		return user.User{}, err
	}
	if userModel.ID == 0 {
		return user.User{}, errors.New("account does not exist")
	}

	if !userModel.ComparePassword(password) {
		return user.User{}, errors.New("wrong password")
	}

	return userModel, nil
}

func LoginByPhone(ctx context.Context, phone string) (user.User, error) {
	userModel, err := user.GetByPhone(ctx, phone)
	if err != nil {
		return user.User{}, err
	}
	if userModel.ID == 0 {
		return user.User{}, errors.New("mobile number is not registered")
	}

	return userModel, nil
}

func CurrentUser(c *gin.Context) user.User {
	userModel, ok := c.MustGet("current_user").(user.User)
	if !ok {
		logger.LogIf(errors.New("could not get user"))
		return user.User{}
	}
	return userModel
}

func CurrentUID(c *gin.Context) string {
	return c.GetString("current_user_id")
}
