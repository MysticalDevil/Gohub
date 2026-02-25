package requests

import (
	"github.com/gin-gonic/gin"
	"gohub/app/requests/validators"
)

type VerifyCodePhoneRequest struct {
	CaptchaID     string `json:"captcha_id,omitempty" valid:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer,omitempty" valid:"captcha_answer"`
	Phone         string `json:"phone,omitempty" valid:"phone"`
}

type VerifyCodeEmailRequest struct {
	CaptchaID     string `json:"captcha_id,omitempty" valid:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer,omitempty" valid:"captcha_answer"`
	Email         string `json:"email,omitempty" valid:"email"`
}

// VerifyCodePhone Validate the form, return the length equal to zero to pass
func VerifyCodePhone(data any, _ *gin.Context) map[string][]string {
	rules := MapData{
		"phone":          []string{"required", "digits:11"},
		"captcha_id":     []string{"required"},
		"captcha_answer": []string{"required", "digits:6"},
	}

	messages := MapData{
		"phone": []string{
			"required:The mobile phone number is required, and the parameter name is 'phone'",
			"digits:Mobile number must be 11 digits",
		},
		"captcha_id": []string{
			"required:The ID of the image verification code is required",
		},
		"captcha_answer": []string{
			"required:Image verification code answer is required",
			"digits:The image verification code must be a 6-digit number",
		},
	}

	errs := validate(data, rules, messages)

	// Image verify code
	_data := data.(*VerifyCodePhoneRequest)
	errs = validators.ValidateCaptcha(_data.CaptchaID, _data.CaptchaAnswer, errs)

	return errs
}

// VerifyCodeEmail Validate the form, return the length equal to zero to pass
func VerifyCodeEmail(data any, _ *gin.Context) map[string][]string {
	rules := MapData{
		"email":          []string{"required", "min:4", "max:30", "email"},
		"captcha_id":     []string{"required"},
		"captcha_answer": []string{"required", "digits:6"},
	}

	messages := MapData{
		"phone": []string{
			"required:The mobile phone number is required, and the parameter name is 'phone'",
			"min:Email length must be greater than 4",
			"max:Email length must be less than 30",
			"email:Email format is incorrect, please provide a valid valid address",
		},
		"captcha_id": []string{
			"required:The ID of the image verification code is required",
		},
		"captcha_answer": []string{
			"required:Image verification code answer is required",
			"digits:The image verification code must be a 6-digit number",
		},
	}

	errs := validate(data, rules, messages)

	// Image verify code
	_data := data.(*VerifyCodeEmailRequest)
	errs = validators.ValidateCaptcha(_data.CaptchaID, _data.CaptchaAnswer, errs)

	return errs
}
