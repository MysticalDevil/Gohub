// Package auth Authorization related logic
package auth

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"gohub/app/models/user"
	"gohub/pkg/logger"
)

// Attempt Try to log in
func Attempt(ctx context.Context, email, password string) (user.User, error) {
	userModel := user.GetByUtil(ctx, email)
	if userModel.ID == 0 {
		return user.User{}, errors.New("account does not exist")
	}

	if !userModel.ComparePassword(password) {
		return user.User{}, errors.New("wrong password")
	}

	return userModel, nil
}

// LoginByPhone Login specified user
func LoginByPhone(ctx context.Context, phone string) (user.User, error) {
	userModel := user.GetByPhone(ctx, phone)
	if userModel.ID == 0 {
		return user.User{}, errors.New("mobile number is not registered")
	}

	return userModel, nil
}

// CurrentUser Get the currently logged-in user from gin.Context
func CurrentUser(c *gin.Context) user.User {
	userModel, ok := c.MustGet("current_user").(user.User)
	if !ok {
		logger.LogIf(errors.New("could not get user"))
		return user.User{}
	}
	return userModel
}

// CurrentUID Get the current login user ID from gin.Context
func CurrentUID(c *gin.Context) string {
	return c.GetString("current_user_id")
}
