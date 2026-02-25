// Package validators Store custom rules and validators
package validators

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
	"gohub/app/captcha"
	"gohub/app/verifycode"
	"gohub/pkg/database"
)

// RegisterCustomValidations registers custom validation rules.
func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidationCtx("not_exists", ValidateFieldNotExist)
	v.RegisterValidationCtx("exists", ValidateFieldExist)
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
func ValidateFieldNotExist(ctx context.Context, field validator.FieldLevel) bool {
	rng := splitParam(field.Param())
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
	query := database.DBWithContext(ctx).Table(tableName).Where(dbField+" = ?", requestValue)

	if len(exceptID) > 0 {
		query.Where("id != ?", exceptID)
	}

	var count int64
	query.Count(&count)

	return count == 0
}

// ValidateFieldExist Customize rules, verify that the field exists in the table
func ValidateFieldExist(ctx context.Context, field validator.FieldLevel) bool {
	rng := splitParam(field.Param())
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
	database.DBWithContext(ctx).Table(tableName).Where(dbField+"= ?", requestValue).Count(&count)

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
	parts := splitParam(field.Param())
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
	parts := splitParam(field.Param())
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
	return slices.Contains(splitParam(field.Param()), value)
}

func ValidateNotIn(field validator.FieldLevel) bool {
	value := fmt.Sprint(field.Field().Interface())
	return !slices.Contains(splitParam(field.Param()), value)
}

func ValidateFileExt(field validator.FieldLevel) bool {
	return ValidateFileRuleValue("ext", field.Field().Interface(), field.Param())
}

func ValidateFileSize(field validator.FieldLevel) bool {
	return ValidateFileRuleValue("size", field.Field().Interface(), field.Param())
}

func splitParam(param string) []string {
	if strings.Contains(param, "_") {
		return strings.Split(param, "_")
	}
	if strings.Contains(param, "|") {
		return strings.Split(param, "|")
	}
	return strings.Split(param, ",")
}

// ValidateFileRuleValue validates file rules without relying on validator tags.
func ValidateFileRuleValue(ruleName string, value any, param string) bool {
	file := normalizeFileHeader(value)
	if file == nil {
		return true
	}

	switch ruleName {
	case "ext":
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(file.Filename)), ".")
		return slices.ContainsFunc(splitParam(param), func(allowed string) bool {
			return strings.ToLower(allowed) == ext
		})
	case "size":
		limit, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return false
		}
		return file.Size <= limit
	default:
		return true
	}
}

func normalizeFileHeader(value any) *multipart.FileHeader {
	if value == nil {
		return nil
	}
	if file, ok := value.(*multipart.FileHeader); ok {
		return file
	}
	if file, ok := value.(multipart.FileHeader); ok {
		return &file
	}
	return nil
}
