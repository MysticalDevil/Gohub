package requests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gohub/app/requests/validators"
	"gohub/pkg/response"
	"reflect"
	"strings"
	"sync"
)

// ValidatorFunc Validate function type
type ValidatorFunc func(any, *gin.Context) map[string][]string

type MapData map[string][]string

var validateOnce sync.Once
var validateEngine *validator.Validate

// Validate Controller call example:
//
//	if ok := requests.Validate(c, &requests.UserSaveRequest{}, requests.UserSave); ! ok {
//	    return
//	}
func Validate(c *gin.Context, obj any, handler ValidatorFunc) bool {
	// Parse request, support json data, form request and url query
	if err := c.ShouldBind(obj); err != nil {
		response.BadRequest(
			c, err,
			"Request parsing error, please confirm the format is correct. "+
				"Please use the `multipart` header for uploading files, and use json format for parameters",
		)
		fmt.Println(err.Error())
		return false
	}

	// Validate form
	errs := handler(obj, c)
	if len(errs) > 0 {
		response.ValidationError(c, errs)
		return false
	}

	return true
}

func validate(data any, rules, messages MapData) map[string][]string {
	return validateWithRules(data, rules, messages)
}

func validateFile(_ *gin.Context, data any, rules, messages MapData) map[string][]string {
	return validateWithRules(data, rules, messages)
}

func validateWithRules(data any, rules, messages MapData) map[string][]string {
	errs := make(map[string][]string)
	msgLookup := buildMessageLookup(messages)
	v := getValidator()

	for field, fieldRules := range rules {
		fieldName := strings.TrimPrefix(field, "file:")
		value, ok := findFieldValue(data, fieldName)
		if !ok {
			continue
		}

		for _, rule := range fieldRules {
			ruleName, ruleParam := splitRule(rule)
			if isEmptyValue(value) && ruleName != "required" {
				continue
			}
			if err := v.Var(value, buildTag(ruleName, ruleParam)); err != nil {
				msg := msgLookup[fieldName][ruleName]
				if msg == "" {
					msg = err.Error()
				}
				errs[fieldName] = append(errs[fieldName], msg)
			}
		}
	}

	return errs
}

func getValidator() *validator.Validate {
	validateOnce.Do(func() {
		validateEngine = validator.New()
		validators.RegisterCustomValidations(validateEngine)
	})

	return validateEngine
}

func splitRule(rule string) (string, string) {
	parts := strings.SplitN(rule, ":", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

func buildTag(name, param string) string {
	if param == "" {
		return name
	}
	return name + "=" + param
}

func buildMessageLookup(messages MapData) map[string]map[string]string {
	lookup := make(map[string]map[string]string)
	for field, fieldMessages := range messages {
		if _, ok := lookup[field]; !ok {
			lookup[field] = make(map[string]string)
		}
		for _, message := range fieldMessages {
			name, text := splitRule(message)
			if text == "" {
				continue
			}
			lookup[field][name] = text
		}
	}

	return lookup
}

func findFieldValue(data any, fieldName string) (any, bool) {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil, false
	}

	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		if fieldTagMatch(field, fieldName) {
			fieldValue := value.Field(i)
			if fieldValue.Kind() == reflect.Pointer && fieldValue.IsNil() {
				return nil, true
			}
			return fieldValue.Interface(), true
		}
	}

	return nil, false
}

func fieldTagMatch(field reflect.StructField, name string) bool {
	if tagValue := tagName(field.Tag.Get("valid")); tagValue == name {
		return true
	}
	if tagValue := tagName(field.Tag.Get("json")); tagValue == name {
		return true
	}
	if tagValue := tagName(field.Tag.Get("form")); tagValue == name {
		return true
	}
	return false
}

func tagName(tag string) string {
	if tag == "" {
		return ""
	}
	parts := strings.Split(tag, ",")
	return parts[0]
}

func isEmptyValue(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		return v.IsNil()
	case reflect.String:
		return strings.TrimSpace(v.String()) == ""
	}
	return false
}
