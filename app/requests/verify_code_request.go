package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"gohub/pkg/captcha"
)

type VerifyCodePhoneRequest struct {
	CaptchaID     string `json:"captcha_id,omitempty" valid:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer,omitempty" valid:"captcha_answer"`
	Phone         string `json:"phone,omitempty" valid:"phone"`
}

// VerifyCodePhone Validate the form, return the length equal to zero to pass
func VerifyCodePhone(data any, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"phone":          []string{"required", "digits:11"},
		"captcha_id":     []string{"required"},
		"captcha_answer": []string{"required", "digits:6"},
	}

	messages := govalidator.MapData{
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
	if ok := captcha.NewCaptcha().VerifyCaptcha(_data.CaptchaID, _data.CaptchaAnswer); !ok {
		errs["captcha_answer"] = append(errs["captcha_answer"], "Image verification code error")
	}

	return errs
}
