// Package validators Store custom rules and validators
package validators

import "gohub/pkg/captcha"

// ValidateCaptcha Customize rules, verify [picture verification code]
func ValidateCaptcha(captchaID, captchaAnswer string, errs map[string][]string) map[string][]string {
	if ok := captcha.NewCaptcha().VerifyCaptcha(captchaID, captchaAnswer); !ok {
		errs["captcha_answer"] = append(errs["captcha_answer"], "Image verification code error")
	}
	return errs
}
