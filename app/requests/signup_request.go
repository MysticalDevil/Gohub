// Package requests Handling request data and form validation
package requests

import (
	"github.com/gin-gonic/gin"
	"gohub/app/requests/validators"
)

type SignupPhoneExistRequest struct {
	Phone string `json:"phone,omitempty" valid:"phone"`
}

// SignupPhoneExist Phone number registration validator
func SignupPhoneExist(data any, c *gin.Context) map[string][]string {
	// Custom validation rules
	rules := MapData{
		"phone": []string{"required", "digits:11"},
	}

	// Customize the prompt when there is an error in the validation
	messages := MapData{
		"phone": []string{
			"required:Phone number is required, parameter name: phone",
			"digits:Phone number must be 11 digits long",
		},
	}

	return validate(c, data, rules, messages)
}

type SignupEmailExistRequest struct {
	Email string `json:"email,omitempty" valid:"email"`
}

// SignupEmailExist Email registration validator
func SignupEmailExist(data any, c *gin.Context) map[string][]string {
	// Custom validation rules
	rules := MapData{
		"email": []string{"required", "min:4", "max:30", "email"},
	}

	// Customize the prompt when there is an error in the validation
	messages := MapData{
		"email": []string{
			"required:Email is required, parameter name: email",
			"min:Email length needs to be greater than 4",
			"max:Email length needs to be less than 30",
			"email:The email format is incorrect, please provide a valid email address",
		},
	}

	return validate(c, data, rules, messages)
}

// SignupUsingPhoneRequest Request information via mobile phone registration
type SignupUsingPhoneRequest struct {
	Phone           string `json:"phone,omitempty" valid:"phone"`
	VerifyCode      string `json:"verify_code,omitempty" valid:"verify_code"`
	Name            string `json:"name" valid:"name"`
	Password        string `json:"password,omitempty" valid:"password"`
	PasswordConfirm string `json:"password_confirm,omitempty" valid:"password_confirm"`
}

func SignupUsingPhone(data any, c *gin.Context) map[string][]string {
	rules := MapData{
		"phone":            []string{"required", "digits:11", "not_exists:users,phone"},
		"name":             []string{"required", "alphanum", "between:3,20", "not_exists:users,name"},
		"password":         []string{"required", "min:6"},
		"password_confirm": []string{"required"},
		"verify_code":      []string{"required", "digits:6"},
	}

	messages := MapData{
		"phone": []string{
			"required:The mobile phone number is required, and the parameter name is 'phone'",
			"digits:Mobile number must be 11 digits long",
		},
		"name": []string{
			"required:Username is required",
			"alphanum:Username is malformed, only numbers and English are allowed",
			"between:Username length must be between 3 and 20",
		},
		"password": []string{
			"required:Password is required",
			"min:Password length must be greater than 6",
		},
		"password_confirm": []string{
			"required:Password confirm is required",
		},
		"verify_code": []string{
			"required:Verification code answer is required",
			"digits:The verification code must be a 6-digit number",
		},
	}

	errs := validate(c, data, rules, messages)

	_data := data.(*SignupUsingPhoneRequest)
	errs = validators.ValidatePasswordConfirm(_data.Password, _data.PasswordConfirm, errs)
	errs = validators.ValidateVerifyCode(_data.Phone, _data.VerifyCode, errs)

	return errs
}

// SignupUsingEmailRequest Request information via email registration
type SignupUsingEmailRequest struct {
	Email           string `json:"email,omitempty" valid:"email"`
	VerifyCode      string `json:"verify_code,omitempty" valid:"verify_code"`
	Name            string `json:"name" valid:"name"`
	Password        string `json:"password,omitempty" valid:"password"`
	PasswordConfirm string `json:"password_confirm,omitempty" valid:"password_confirm"`
}

func SignupUsingEmail(data any, c *gin.Context) map[string][]string {
	rules := MapData{
		"email":            []string{"required", "min:4", "max:30", "email", "not_exists:users,email"},
		"name":             []string{"required", "alphanum", "between:3,20", "not_exists:users,name"},
		"password":         []string{"required", "min:6"},
		"password_confirm": []string{"required"},
		"verify_code":      []string{"required", "digits:6"},
	}

	messages := MapData{
		"email": []string{
			"required:Email is required",
			"min:Email length must be greater than 4",
			"max:Email length must be less than 30",
			"email:The email format is incorrect, please provide a valid email address",
			"not_exists:Email is occupied",
		},
		"name": []string{
			"required:Username is required",
			"alphanum:Username is malformed, only numbers and English are allowed",
			"between:Username length must be between 3 and 20",
		},
		"password": []string{
			"required:Password is required",
			"min:Password length must be greater than 6",
		},
		"password_confirm": []string{
			"required:Password confirm is required",
		},
		"verify_code": []string{
			"required:Verification code answer is required",
			"digits:The verification code must be a 6-digit number",
		},
	}

	errs := validate(c, data, rules, messages)

	_data := data.(*SignupUsingEmailRequest)
	errs = validators.ValidatePasswordConfirm(_data.Password, _data.PasswordConfirm, errs)
	errs = validators.ValidateVerifyCode(_data.Email, _data.VerifyCode, errs)

	return errs
}
