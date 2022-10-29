// Package requests Handling request data and form validation
package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type SignupPhoneExistRequest struct {
	Phone string `json:"phone,omitempty" valid:"phone"`
}

// ValidateSignupPhoneExist Phone number registration validator
func ValidateSignupPhoneExist(data any, _ *gin.Context) map[string][]string {
	// Custom validation rules
	rules := govalidator.MapData{
		"phone": []string{"required", "digits:11"},
	}

	// Customize the prompt when there is an error in the validation
	messages := govalidator.MapData{
		"phone": []string{
			"required:Phone number is required, parameter name: phone",
			"digits:Phone number must be 11 digits long",
		},
	}

	// Configuration initialization
	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		TagIdentifier: "valid",
		Messages:      messages,
	}

	// Start validate
	return govalidator.New(opts).ValidateStruct()
}

type SignupEmailExistRequest struct {
	Email string `json:"email,omitempty" valid:"email"`
}

// ValidateSignupEmailExist Email registration validator
func ValidateSignupEmailExist(data any, _ *gin.Context) map[string][]string {
	// Custom validation rules
	rules := govalidator.MapData{
		"email": []string{"required", "min:4", "max:30", "email"},
	}

	// Customize the prompt when there is an error in the validation
	messages := govalidator.MapData{
		"email": []string{
			"required:Email is required, parameter name: email",
			"min:Email length needs to be greater than 4",
			"max:Email length needs to be less than 30",
			"email:The email format is incorrect, please provide a valid email address",
		},
	}

	// Configuration initialization
	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		TagIdentifier: "valid",
		Messages:      messages,
	}

	// Start validate
	return govalidator.New(opts).ValidateStruct()
}
