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
func LoginByPhone(data any, c *gin.Context) map[string][]string {
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
