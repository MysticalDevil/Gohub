// Package validators Store custom rules and validators
package validators

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gohub/pkg/captcha"
	"gohub/pkg/database"
	"gohub/pkg/verifycode"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

// RegisterCustomValidations registers custom validation rules.
func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("not_exists", ValidateFieldNotExist)
	v.RegisterValidation("exists", ValidateFieldExist)
	v.RegisterValidation("max_cn", ValidateMaxCn)
	v.RegisterValidation("min_cn", ValidateMinCn)
	v.RegisterValidation("digits", ValidateDigits)
	v.RegisterValidation("between", ValidateBetween)
	v.RegisterValidation("numeric_between", ValidateNumericBetween)
	v.RegisterValidation("in", ValidateIn)
	v.RegisterValidation("not_in", ValidateNotIn)
	v.RegisterValidation("ext", ValidateFileExt)
	v.RegisterValidation("size", ValidateFileSize)
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

// ValidateFieldNotExist Customize rules, verify that the field already exists in the table
func ValidateFieldNotExist(field validator.FieldLevel) bool {
	rng := strings.Split(field.Param(), ",")
	if len(rng) < 2 {
		return false
	}

	// The first parameter, the table name, such as users
	tableName := rng[0]
	// The second parameter, the field name, such as email or phone
	dbField := rng[1]

	// The third parameter, the exclusion ID
	var exceptID string
	if len(rng) > 2 {
		exceptID = rng[2]
	}

	requestValue := fmt.Sprint(field.Field().Interface())

	// Splicing SQL
	query := database.DB.Table(tableName).Where(dbField+" = ?", requestValue)

	if len(exceptID) > 0 {
		query.Where("id != ?", exceptID)
	}

	var count int64
	query.Count(&count)

	return count == 0
}

// ValidateFieldExist Customize rules, verify that the field exists in the table
func ValidateFieldExist(field validator.FieldLevel) bool {
	rng := strings.Split(field.Param(), ",")
	if len(rng) < 2 {
		return false
	}

	// The first parameter, the table name
	tableName := rng[0]
	// THe second parameter, the field name, eg: id
	dbField := rng[1]

	// The data requested by the user
	requestValue := fmt.Sprint(field.Field().Interface())

	// Query database
	var count int64
	database.DB.Table(tableName).Where(dbField+"= ?", requestValue).Count(&count)

	return count > 0
}

// ValidateMaxCn Customize rules, verify the maximum length of Chinese characters
func ValidateMaxCn(field validator.FieldLevel) bool {
	valLength := utf8.RuneCountInString(field.Field().String())
	maxLen, err := strconv.Atoi(field.Param())
	if err != nil {
		return false
	}
	return valLength <= maxLen
}

// ValidateMinCn Customize rules, verify the minimum length of Chinese characters
func ValidateMinCn(field validator.FieldLevel) bool {
	valLength := utf8.RuneCountInString(field.Field().String())
	minLen, err := strconv.Atoi(field.Param())
	if err != nil {
		return false
	}
	return valLength >= minLen
}

func ValidateDigits(field validator.FieldLevel) bool {
	raw := fmt.Sprint(field.Field().Interface())
	if raw == "" {
		return false
	}
	length, err := strconv.Atoi(field.Param())
	if err != nil {
		return false
	}
	if len(raw) != length {
		return false
	}
	_, err = strconv.Atoi(raw)
	return err == nil
}

func ValidateBetween(field validator.FieldLevel) bool {
	parts := strings.Split(field.Param(), ",")
	if len(parts) != 2 {
		return false
	}
	minLen, err := strconv.Atoi(parts[0])
	if err != nil {
		return false
	}
	maxLen, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}
	valLength := utf8.RuneCountInString(field.Field().String())
	return valLength >= minLen && valLength <= maxLen
}

func ValidateNumericBetween(field validator.FieldLevel) bool {
	parts := strings.Split(field.Param(), ",")
	if len(parts) != 2 {
		return false
	}
	minVal, err := strconv.Atoi(parts[0])
	if err != nil {
		return false
	}
	maxVal, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}
	value, err := strconv.Atoi(fmt.Sprint(field.Field().Interface()))
	if err != nil {
		return false
	}
	return value >= minVal && value <= maxVal
}

func ValidateIn(field validator.FieldLevel) bool {
	value := fmt.Sprint(field.Field().Interface())
	for _, item := range strings.Split(field.Param(), ",") {
		if value == item {
			return true
		}
	}
	return false
}

func ValidateNotIn(field validator.FieldLevel) bool {
	value := fmt.Sprint(field.Field().Interface())
	for _, item := range strings.Split(field.Param(), ",") {
		if value == item {
			return false
		}
	}
	return true
}

func ValidateFileExt(field validator.FieldLevel) bool {
	file, ok := field.Field().Interface().(*multipart.FileHeader)
	if !ok || file == nil {
		return true
	}
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(file.Filename)), ".")
	for _, allowed := range strings.Split(field.Param(), ",") {
		if strings.ToLower(allowed) == ext {
			return true
		}
	}
	return false
}

func ValidateFileSize(field validator.FieldLevel) bool {
	file, ok := field.Field().Interface().(*multipart.FileHeader)
	if !ok || file == nil {
		return true
	}
	limit, err := strconv.ParseInt(field.Param(), 10, 64)
	if err != nil {
		return false
	}
	return file.Size <= limit
}
