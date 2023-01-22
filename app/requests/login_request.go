package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"gohub/app/requests/validators"
)

type LoginByPhoneRequest struct {
	Phone      string `json:"phone,omitempty" valid:"phone"`
	VerifyCode string `json:"verify_code,omitempty" valid:"verify_code"`
}

// LoginByPhone Validate the form
func LoginByPhone(data any, _ *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"phone":       []string{"required", "digits:11"},
		"verify_code": []string{"required", "digits:6"},
	}

	messages := govalidator.MapData{
		"phone": []string{
			"required:The mobile phone number is required, and the parameter name is 'phone'",
			"digits:Mobile number must be 11 digits long",
		},
		"verify_code": []string{
			"required:Verification code answer is required",
			"digits:The verification code must be a 6-digit number",
		},
	}

	errs := validate(data, rules, messages)

	_data := data.(*LoginByPhoneRequest)
	errs = validators.ValidateVerifyCode(_data.Phone, _data.VerifyCode, errs)

	return errs
}

type LoginByPasswordRequest struct {
	CaptchaID     string `json:"captcha_id,omitempty" valid:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer,omitempty" valid:"captcha_answer"`

	LoginID  string `json:"login_id" valid:"login_id"`
	Password string `json:"password,omitempty" valid:"password"`
}

// LoginByPassword Validate the form
func LoginByPassword(data any, _ *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"login_id":       []string{"required", "min:3"},
		"password":       []string{"required", "min:6"},
		"captcha_id":     []string{"required"},
		"captcha_answer": []string{"required", "digits:6"},
	}

	messages := govalidator.MapData{
		"login_id": []string{
			"required:The login ID is required, and supports mobile phone number, email address and user name",
			"min:Login ID length must be greater than 3",
		},
		"password": []string{
			"required:Password is required",
			"min:Password length must be greater than 6",
		},
		"captcha_id": []string{
			"requires:Image verification code ID is required",
		},
		"captcha_answer": []string{
			"required:Image verification code answer is required",
			"digits:The picture verification code must be 6 digits in length",
		},
	}

	errs := validate(data, rules, messages)

	_data := data.(*LoginByPasswordRequest)
	errs = validators.ValidateCaptcha(_data.CaptchaID, _data.CaptchaAnswer, errs)

	return errs
}
