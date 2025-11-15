package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"user-notes-api/auth"
	"user-notes-api/services"
	"user-notes-api/testing/testutils/servicemocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

/*
	1. username already exists
	2. unknown error
	3. empty password?
*/

func TestAuthControllerRegistrationSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"username": "Alice", "password": "pwd"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	mockLoginService := new(servicemocks.MockLoginService)
	mockRegistrationService := new(servicemocks.MockRegistrationService)
	mockRegistrationService.On("Register", c.Request.Context(), auth.Credentials{Username: "Alice", Password: "pwd"}).Return("jwt", nil)

	authController := NewAuthController(mockLoginService, mockRegistrationService)

	authController.Register(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "jwt")
	mockLoginService.AssertExpectations(t)

}

func TestAuthControllerLoginSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"username": "Alice", "password": "pwd"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	mockLoginService := new(servicemocks.MockLoginService)
	mockLoginService.On("Login", c.Request.Context(), auth.Credentials{Username: "Alice", Password: "pwd"}).Return("jwt", nil)

	mockRegistrationService := new(servicemocks.MockRegistrationService)

	authController := NewAuthController(mockLoginService, mockRegistrationService)

	authController.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "jwt")
	mockLoginService.AssertExpectations(t)
}

func TestAuthControllerLoginFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"username": "Alice", "password": "wrong_pwd"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	wrongPwdError := new(services.ErrorWrongPassword)
	wrongPwdError.Username = "Alice"

	mockLoginService := new(servicemocks.MockLoginService)
	mockLoginService.On("Login", c.Request.Context(), auth.Credentials{Username: "Alice", Password: "wrong_pwd"}).Return("", wrongPwdError)

	mockRegistrationService := new(servicemocks.MockRegistrationService)

	authController := NewAuthController(mockLoginService, mockRegistrationService)

	authController.Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "wrong password")
	mockLoginService.AssertExpectations(t)
}

func TestAuthControllerLoginUserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := []byte(`{"username": "Unknown user", "password": "pwd"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	notFoundError := new(auth.ErrorNotFound)
	notFoundError.Username = "Unknown user"

	mockLoginService := new(servicemocks.MockLoginService)
	mockLoginService.On("Login", c.Request.Context(), auth.Credentials{Username: "Unknown user", Password: "pwd"}).Return("", notFoundError)

	mockRegistrationService := new(servicemocks.MockRegistrationService)

	authController := NewAuthController(mockLoginService, mockRegistrationService)

	authController.Login(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "user not found")
	mockLoginService.AssertExpectations(t)
}
