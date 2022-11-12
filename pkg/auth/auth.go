// Package auth Authorization related logic
package auth

import (
	"errors"
	"gohub/app/models/user"
)

// Attempt Try to log in
func Attempt(email, password string) (user.User, error) {
	userModel := user.GetByUtil(email)
	if userModel.ID == 0 {
		return user.User{}, errors.New("account does not exist")
	}

	if !userModel.ComparePassword(password) {
		return user.User{}, errors.New("wrong password")
	}

	return userModel, nil
}

// LoginByPhone Login specified user
func LoginByPhone(phone string) (user.User, error) {
	userModel := user.GetByPhone(phone)
	if userModel.ID == 0 {
		return user.User{}, errors.New("mobile number is not registered")
	}

	return userModel, nil
}
