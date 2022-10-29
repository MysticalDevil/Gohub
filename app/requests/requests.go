package requests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"net/http"
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
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Request parsing error, please confirm the format is correct. " +
				"Please use the `multipart` header for uploading files, and use json format for parameters",
			"error": err.Error(),
		})
		fmt.Println(err.Error())
		return false
	}

	// Validate form
	errs := handler(obj, c)
	if len(errs) > 0 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Request verification failed, please see errors for details",
			"errors":  errs,
		})
		return false
	}

	return true
}

func validate(data any, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
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
