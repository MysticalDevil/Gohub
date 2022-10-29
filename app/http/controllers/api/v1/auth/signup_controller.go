// Package auth Handle user authentication related logic
package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/app/models/user"
	"gohub/app/requests"
	"net/http"
)

// SignupController Register controller
type SignupController struct {
	v1.BaseAPIController
}

// IsPhoneExist Check if the phone number is registered
func (sc *SignupController) IsPhoneExist(c *gin.Context) {
	request := requests.SignupPhoneExistRequest{}

	// Parse Json request
	if err := c.ShouldBindJSON(&request); err != nil {
		// Parsing failed, returning 422 status code and error message
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		// print error message
		fmt.Println(err.Error())
		// interrupt request
		return
	}

	// form validation
	errs := requests.ValidateSignupPhoneExist(&request, c)

	if len(errs) > 0 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"errors": errs,
		})
		return
	}

	// Check database and return response
	c.JSON(http.StatusOK, gin.H{
		"exist": user.IsPhoneExist(request.Phone),
	})
}
