package controllers

import (
	"errors"
	"net/http"

	"user-notes-api/auth"
	"user-notes-api/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	LoginService        services.LoginServiceIfc
	RegistrationService services.RegistrationServiceIfc
}

func NewAuthController(login_service services.LoginServiceIfc, registration_service services.RegistrationServiceIfc) *AuthController {
	controller := AuthController{LoginService: login_service, RegistrationService: registration_service}
	return &controller
}

/* func (a *AuthController) Register(c *gin.Context) {
	request_ctx := c.Request.Context()
	token_string, err := a.RegistrationService.Register(request_ctx, credentials)
} */

func (a *AuthController) Login(c *gin.Context) {
	var credentials auth.Credentials
	err := c.Bind(&credentials)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	request_ctx := c.Request.Context()
	token_string, err := a.LoginService.Login(request_ctx, credentials)

	if err != nil {
		var wrongPwdError *services.ErrorWrongPassword
		var notFoundError *auth.ErrorNotFound

		if errors.As(err, &wrongPwdError) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
			return
		} else if errors.As(err, &notFoundError) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"token": token_string})
}
