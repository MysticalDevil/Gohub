// Package validators Store custom rules and validators
package validators

import (
	"errors"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"gohub/pkg/captcha"
	"gohub/pkg/database"
	"gohub/pkg/verifycode"
	"strings"
)

func init() {
	// Custom rule not_exists, verify that the request data must not exist in the database
	// Often used to ensure that the value of a field in the database is unique
	// There are two types of not_exists parameters, one is 2 parameters and the other is 3 parameters
	// not_exists:users,email Check whether the same information exists in the database table
	// not_exists:users,email,32 Exclude users with id 32
	govalidator.AddCustomRule("not_exists", ValidateFieldExist)
}

// ValidateCaptcha Customize rules, verify [picture verification code]
func ValidateCaptcha(captchaID, captchaAnswer string, errs map[string][]string) map[string][]string {
	if ok := captcha.NewCaptcha().VerifyCaptcha(captchaID, captchaAnswer); !ok {
		errs["captcha_answer"] = append(errs["captcha_answer"], "Image verification code error")
	}
	return errs
}

// ValidatePasswordConfirm Customize rules, check if the two passwords match
func ValidatePasswordConfirm(password, passwordConfirm string, errs map[string][]string) map[string][]string {
	if password != passwordConfirm {
		errs["password_confirm"] = append(errs["password_confirm"], "The passwords entered twice do not match")
	}
	return errs
}

// ValidateVerifyCode Customize rules, verify [Mobile/Email Verification Code]
func ValidateVerifyCode(key, answer string, errs map[string][]string) map[string][]string {
	if ok := verifycode.NewVerifyCode().CheckAnswer(key, answer); !ok {
		errs["verify_code"] = append(errs["verify_code"], "Verification code error")
	}
	return errs
}

// ValidateFieldExist Customize rules, verify that the field already exists in the table
func ValidateFieldExist(field, rules, message string, value any) error {
	rng := strings.Split(strings.TrimPrefix(rules, "not_exists:"), ",")

	// The first parameter, the table name, such as users
	tableName := rng[0]
	// The second parameter, the field name, such as email or phone
	dbField := rng[1]

	// The third parameter, the exclusion ID
	var exceptID string
	if len(rng) > 2 {
		exceptID = rng[2]
	}

	// The data requested by the user
	requestValue := value.(string)

	// Splicing SQL
	query := database.DB.Table(tableName).Where(dbField+" = ?", requestValue)

	if len(exceptID) > 0 {
		query.Where("id != ?", exceptID)
	}

	var count int64
	query.Count(&count)

	if count != 0 {
		if message != "" {
			return errors.New(message)
		}
		return fmt.Errorf("%v is taken", requestValue)
	}
	return nil
}
