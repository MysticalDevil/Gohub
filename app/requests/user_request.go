package requests

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gohub/app/requests/validators"
	"gohub/pkg/auth"
)

type UserUpdateProfileRequest struct {
	Name         string `valid:"name" json:"name"`
	City         string `valid:"city" json:"city"`
	Introduction string `valid:"introduction" json:"introduction"`
}

func UserUpdateProfile(data any, c *gin.Context) map[string][]string {
	uid := auth.CurrentUID(c)
	rules := MapData{
		"name":         []string{"required", "alphanum", "between:3,20", "not_exists:users,name," + uid},
		"introduction": []string{"min_cn:4", "max_cn:240"},
		"city":         []string{"min_cn:2", "max_cn:20"},
	}

	messages := MapData{
		"name": []string{
			"required:Username is required",
			"alphanum:Username is malformed, only numbers and English are allowed",
			"between:Username length must be between 3 and 20",
			"not_exists:Username already taken",
		},
		"introduction": []string{
			"min_cn:The introduction length must be at least 4 characters",
			"max_cn:The introduction length must be at most 240 characters",
		},
		"city": []string{
			"min_cn:The city length must be at least 2 characters",
			"max_cn:The city length must be at most 20 characters",
		},
	}

	return validate(c, data, rules, messages)
}

type UserUpdateEmailRequest struct {
	Email      string `json:"email,omitempty" valid:"email"`
	VerifyCode string `json:"verify_code,omitempty" valid:"verify_code"`
}

func UserUpdateEmail(data any, c *gin.Context) map[string][]string {
	currentUser := auth.CurrentUser(c)

	rules := MapData{
		"email": []string{
			"required",
			"min:4",
			"max:30",
			"email",
			"not_exists:users,email," + currentUser.GetStringID(),
			"not_in:" + currentUser.Email,
		},
		"verify_code": []string{"required", "digits:6"},
	}

	messages := MapData{
		"email": []string{
			"required:Email is required",
			"min:Email length must be greater than 4",
			"max:Email length must be less than 30",
			"email:The email format is incorrect, please provide a valid email address",
			"not_exists:Email is occupied",
			"not_in:The new email is the same as the old email",
		},
		"verify_code": []string{
			"required:Verification code answer is required",
			"digits:The verification code must be a 6-digit number",
		},
	}

	errs := validate(c, data, rules, messages)
	_data := data.(*UserUpdateEmailRequest)
	errs = validators.ValidateVerifyCode(_data.Email, _data.VerifyCode, errs)

	return errs
}

type UserUpdatePhoneRequest struct {
	Phone      string `json:"phone,omitempty" valid:"phone"`
	VerifyCode string `json:"verify_code,omitempty" valid:"verify_code"`
}

func UserUpdatePhone(data any, c *gin.Context) map[string][]string {
	currentUser := auth.CurrentUser(c)

	rules := MapData{
		"phone": []string{
			"required",
			"digits:11",
			"not_exists:users,phone," + currentUser.GetStringID(),
			"not_in:" + currentUser.Phone,
		},
		"verify_code": []string{"required", "digits:6"},
	}

	messages := MapData{
		"phone": []string{
			"required:The mobile phone number is required, and the parameter name is 'phone'",
			"digits:Mobile number must be 11 digits long",
			"not_exists:Phone is occupied",
			"not_in:The new phone is the same as the old phone",
		},
		"verify_code": []string{
			"required:Verification code answer is required",
			"digits:The verification code must be a 6-digit number",
		},
	}

	errs := validate(c, data, rules, messages)
	_data := data.(*UserUpdatePhoneRequest)
	errs = validators.ValidateVerifyCode(_data.Phone, _data.VerifyCode, errs)

	return errs
}

type UserUpdatePasswordRequest struct {
	Password           string `json:"password,omitempty" valid:"password"`
	NewPassword        string `json:"new_password,omitempty" valid:"new_password"`
	NewPasswordConfirm string `json:"new_password_confirm,omitempty" valid:"new_password_confirm"`
}

func UserUpdatePassword(data any, c *gin.Context) map[string][]string {
	rules := MapData{
		"password":             []string{"required", "min:6"},
		"new_password":         []string{"required", "min:6"},
		"new_password_confirm": []string{"required", "min:6"},
	}

	messages := MapData{
		"password": []string{
			"required:Password is required",
			"min:Password length must be greater than 6",
		},
		"new_password": []string{
			"required:Password is required",
			"min:Password length must be greater than 6",
		},
		"new_password_confirm": []string{
			"required:Password confirm is required",
			"min:Password confirm length must be greater than 6",
		},
	}

	errs := validate(c, data, rules, messages)
	_data := data.(*UserUpdatePasswordRequest)
	errs = validators.ValidatePasswordConfirm(_data.NewPassword, _data.NewPasswordConfirm, errs)

	return errs
}

type UserUpdateAvatarRequest struct {
	Avatar *multipart.FileHeader `valid:"avatar" form:"avatar"`
}

func UserUpdateAvatar(data any, c *gin.Context) map[string][]string {
	rules := MapData{
		"file:avatar": []string{"ext:png,jpg,jpeg", "size:20971520", "required"},
	}

	messages := MapData{
		"file:avatar": []string{
			"ext:The avatar can only upload pictures in png, jpg, jpeg format",
			"size:The maximum size of the avatar cannot exceed 20MB",
			"required:Image must be uploaded",
		},
	}

	return validateFile(c, data, rules, messages)
}
