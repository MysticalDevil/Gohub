// Package response Response handler
package response

import (
	"errors"
	"fmt"
	"maps"
	"net/http"

	"github.com/gin-gonic/gin"
	"gohub/pkg/logger"
	"gohub/pkg/paginator"
	"gorm.io/gorm"
)

const (
	CodeOK            = "OK"
	CodeCreated       = "CREATED"
	CodeBadRequest    = "ERR_BAD_REQUEST"
	CodeUnauthorized  = "ERR_UNAUTHORIZED"
	CodeForbidden     = "ERR_FORBIDDEN"
	CodeNotFound      = "ERR_NOT_FOUND"
	CodeValidation    = "ERR_VALIDATION"
	CodeUnprocessable = "ERR_UNPROCESSABLE"
	CodeInternal      = "ERR_INTERNAL"
)

type envelope struct {
	Data   any                 `json:"data,omitempty"`
	Msg    string              `json:"msg"`
	Code   string              `json:"code"`
	Errors map[string][]string `json:"errors,omitempty"`
}

// JSON response 200 and JSON data (standard envelope)
func JSON(c *gin.Context, data any) {
	Data(c, data)
}

// Success
// In response to 200 and the JSON data of the preset [operation successful].
// It is called after a [change] operation without [specific return data] is successful,
// such as deletion, password modification, and mobile phone number modification.
func Success(c *gin.Context) {
	respond(c, http.StatusOK, CodeOK, "Successful operation", nil, nil)
}

// Data
// Response 200 and JSON data with data key.
// Called after the execution of the [update operation] is successful,
// such as updating the topic, and returning the updated topic after success
func Data(c *gin.Context, data any) {
	respond(c, http.StatusOK, CodeOK, "OK", data, nil)
}

// Paginated
// Response 200 and JSON data in offset/limit pagination format
func Paginated(c *gin.Context, items any, paging paginator.Paging) {
	respond(c, http.StatusOK, CodeOK, "OK", gin.H{
		"items":  items,
		"offset": paging.Offset,
		"limit":  paging.Limit,
		"total":  paging.Total,
	}, nil)
}

// Created
// Response 201 and JSON data with data key.
// Called after the execution of the [update operation] is successful,
// such as updating the topic, and returning the updated topic after success
func Created(c *gin.Context, data any) {
	respond(c, http.StatusCreated, CodeCreated, "Created", data, nil)
}

// CreatedJSON
// Response 201 and JSON data (standard envelope)
func CreatedJSON(c *gin.Context, data any) {
	Created(c, data)
}

// Abort404
// Response 404
// Use the default message when no msg parameter is passed
func Abort404(c *gin.Context, msg ...string) {
	errorResponse(c, http.StatusNotFound, CodeNotFound,
		defaultMessage("The data does not exist, please confirm that the request is correct", msg...), nil)
}

// Abort403
// Response 403
// Use the default message when no msg parameter is passed
func Abort403(c *gin.Context, msg ...string) {
	errorResponse(c, http.StatusForbidden, CodeForbidden,
		defaultMessage("Insufficient permissions, please confirm that you have the corresponding permissions", msg...), nil)
}

// Abort500
// Response 500
// Use the default message when no msg parameter is passed
func Abort500(c *gin.Context, msg ...string) {
	errorResponse(c, http.StatusInternalServerError, CodeInternal,
		defaultMessage("Internal server error, please try again later", msg...), nil)
}

// BadRequest
// Response 400, use the default message when no msg parameter is passed
// Called when parsing a user request, the format or method of the request is not as expected
func BadRequest(c *gin.Context, err error, msg ...string) {
	logger.LogIf(err)
	errorResponse(c, http.StatusBadRequest, CodeBadRequest,
		defaultMessage(
			"Request parsing error, please confirm whether the request format is correct. "+
				"Please use the 'multipart' header for uploading files and use JSON format for parameters",
			msg...,
		), map[string][]string{"error": {err.Error()}})
}

// Error
// Response 404 or 422, use the default message when no msg parameter is passed
// An error err occurs when processing the request,
// and an error message will be returned, such as a login error, and the corresponding Model cannot be found.
func Error(c *gin.Context, err error, msg ...string) {
	logger.LogIf(err)

	if err == gorm.ErrRecordNotFound {
		Abort404(c)
		return
	}

	errorResponse(c, http.StatusUnprocessableEntity, CodeUnprocessable,
		defaultMessage("Request processing failed, please check the value of error", msg...),
		map[string][]string{"error": {err.Error()}})
}

// ValidationError
// Handling form validation failure errors
func ValidationError(c *gin.Context, errs map[string][]string) {
	joined := joinValidationErrors(errs)
	logger.LogIf(joined)
	errorResponse(c, http.StatusUnprocessableEntity, CodeValidation,
		"Request verification failed, please see errors for details",
		maps.Clone(errs))
}

// Unauthorized
// Response 401, use the default message when no msg parameter is passed
// Called when login fails and jwt parsing fails
func Unauthorized(c *gin.Context, msg ...string) {
	errorResponse(c, http.StatusUnauthorized, CodeUnauthorized,
		defaultMessage("Unauthorized", msg...), nil)
}

func respond(c *gin.Context, status int, code string, msg string, data any, errs map[string][]string) {
	payload := envelope{Msg: msg, Code: code, Data: data, Errors: errs}
	c.JSON(status, payload)
}

func errorResponse(c *gin.Context, status int, code string, msg string, errs map[string][]string) {
	if errs == nil {
		errs = map[string][]string{}
	}
	payload := envelope{Msg: msg, Code: code, Errors: errs}
	c.AbortWithStatusJSON(status, payload)
}

func joinValidationErrors(errs map[string][]string) error {
	if len(errs) == 0 {
		return nil
	}

	joined := make([]error, 0, len(errs))
	for field, messages := range errs {
		for _, message := range messages {
			joined = append(joined, fmt.Errorf("%s: %s", field, message))
		}
	}

	return errors.Join(joined...)
}

// defaultMessage
func defaultMessage(defaultMsg string, msg ...string) (message string) {
	if len(msg) > 0 {
		message = msg[0]
	} else {
		message = defaultMsg
	}
	return
}
