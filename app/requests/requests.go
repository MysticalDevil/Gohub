package requests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"gohub/pkg/response"
)

// ValidatorFunc Validate function type
type ValidatorFunc func(any, *gin.Context) map[string][]string

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

func validate(data any, rules, messages govalidator.MapData) map[string][]string {
	// Configuration initialization
	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		TagIdentifier: "valid",
		Messages:      messages,
	}
	// Start validate
	return govalidator.New(opts).ValidateStruct()
}

func validateFile(c *gin.Context, data any, rules, messages govalidator.MapData) map[string][]string {
	opts := govalidator.Options{
		Request:       c.Request,
		Rules:         rules,
		Messages:      messages,
		TagIdentifier: "valid",
	}

	return govalidator.New(opts).Validate()
}
