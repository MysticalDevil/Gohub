// Package response Response handler
package response

import (
	"github.com/gin-gonic/gin"
	"gohub/pkg/logger"
	"gorm.io/gorm"
	"net/http"
)

// JSON response 200 and JSON data
func JSON(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

// Success
// In response to 200 and the JSON data of the preset [operation successful].
// It is called after a [change] operation without [specific return data] is successful,
// such as deletion, password modification, and mobile phone number modification.
func Success(c *gin.Context) {
	JSON(c, gin.H{
		"success": true,
		"message": "Successful operation",
	})
}

// Data
// Response 200 and JSON data with data key.
// Called after the execution of the [update operation] is successful,
// such as updating the topic, and returning the updated topic after success
func Data(c *gin.Context, data any) {
	JSON(c, gin.H{
		"success": true,
		"data":    data,
	})
}

// Created
// Response 201 and JSON data with data key.
// Called after the execution of the [update operation] is successful,
// such as updating the topic, and returning the updated topic after success
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}

// CreatedJSON
// Response 201 and JSON data
func CreatedJSON(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, data)
}

// Abort404
// Response 404
// Use the default message when no msg parameter is passed
func Abort404(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"message": defaultMessage("The data does not exist, please confirm that the request is correct", msg...),
	})
}

// Abort403
// Response 403
// Use the default message when no msg parameter is passed
func Abort403(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": defaultMessage("Insufficient permissions, please confirm that you have the corresponding permissions", msg...),
	})
}

// Abort500
// Response 500
// Use the default message when no msg parameter is passed
func Abort500(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": defaultMessage("Internal server error, please try again later", msg...),
	})
}

// BadRequest
// Response 400, use the default message when no msg parameter is passed
// Called when parsing a user request, the format or method of the request is not as expected
func BadRequest(c *gin.Context, err error, msg ...string) {
	logger.LogIf(err)
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"message": defaultMessage(
			"Request parsing error, please confirm whether the request format is correct. "+
				"Please use the 'multipart' header for uploading files and use JSON format for parameters",
			msg...,
		),
		"error": err.Error(),
	})
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

	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		"message": defaultMessage("Request processing failed, please check the value of error", msg...),
		"error":   err.Error(),
	})
}

// ValidationError
// Handling form validation failure errors
func ValidationError(c *gin.Context, errors map[string][]string) {
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		"message": "Request verification failed, please see errors for details",
		"errors":  errors,
	})
}

// Unauthorized
// Response 401, use the default message when no msg parameter is passed
// Called when login fails and jwt parsing fails
func Unauthorized(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"message": defaultMessage(
			"Request parsing error, please confirm whether the request format is correct. "+
				"Please use the 'multipart' header for uploading files and use JSON format for parameters",
			msg...
		),
	})
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
