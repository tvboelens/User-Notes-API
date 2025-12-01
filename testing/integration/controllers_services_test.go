package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"user-notes-api/auth"
	"user-notes-api/controllers"
	"user-notes-api/services"
	"user-notes-api/testing/testutils/authmocks"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type JwtToken struct {
	Token string `json:"token"`
}

func TestAuthControllerRegistrationAndLoginSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwt_secret := "jwt_secret"

	login_manager := new(authmocks.MockLoginManager)
	registration_manager := new(authmocks.MockRegistrationManager)

	login_service := services.NewLoginService(login_manager, jwt_secret)
	registration_service := services.NewRegistrationService(registration_manager, jwt_secret)

	authController := controllers.NewAuthController(login_service, registration_service)

	body := []byte(`{"username": "Alice", "password": "secret_pwd"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	req_ctx := c.Request.Context()
	registration_manager.On("RegisterUser", req_ctx, &auth.Credentials{Username: "Alice", Password: "secret_pwd"}).Return(1, nil)

	authController.Register(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
	registration_manager.AssertExpectations(t)

	body, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	var ttoken JwtToken
	err = json.Unmarshal(body, &ttoken)
	if err != nil {
		t.Fatal(err)
	}

	token, err := jwt.ParseWithClaims(ttoken.Token, &services.JwtClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(jwt_secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		t.Fatal(err)
	}

	claims, ok := token.Claims.(*services.JwtClaims)
	assert.True(t, ok)

	issuer, err := claims.GetIssuer()
	assert.NoError(t, err)
	assert.Equal(t, "auth.user-notes-api.local", issuer)

	subject, err := claims.GetSubject()
	assert.NoError(t, err)
	assert.Equal(t, "Alice", subject)

	issuedAt, err := claims.GetIssuedAt()
	assert.NoError(t, err)
	assert.True(t, time.Now().After(issuedAt.Time))

	notBefore, err := claims.GetNotBefore()
	assert.NoError(t, err)
	assert.True(t, time.Now().After(notBefore.Time))

	expirationTime, err := claims.GetExpirationTime()
	assert.NoError(t, err)
	assert.True(t, expirationTime.After(time.Now()))

	body = []byte(`{"username": "Alice", "password": "secret_pwd"}`)
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	req_ctx = c.Request.Context()
	login_manager.On("LoginUser", req_ctx, &auth.Credentials{Username: "Alice", Password: "secret_pwd"}).Return(1, true, nil)

	authController.Login(c)

	err = json.Unmarshal(body, &ttoken)
	if err != nil {
		t.Fatal(err)
	}

	token, err = jwt.ParseWithClaims(ttoken.Token, &services.JwtClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(jwt_secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		t.Fatal(err)
	}

	claims, ok = token.Claims.(*services.JwtClaims)
	assert.True(t, ok)

	issuer, err = claims.GetIssuer()
	assert.NoError(t, err)
	assert.Equal(t, "auth.user-notes-api.local", issuer)

	subject, err = claims.GetSubject()
	assert.NoError(t, err)
	assert.Equal(t, "Alice", subject)

	issuedAt, err = claims.GetIssuedAt()
	assert.NoError(t, err)
	assert.True(t, time.Now().After(issuedAt.Time))

	notBefore, err = claims.GetNotBefore()
	assert.NoError(t, err)
	assert.True(t, time.Now().After(notBefore.Time))

	expirationTime, err = claims.GetExpirationTime()
	assert.NoError(t, err)
	assert.True(t, expirationTime.After(time.Now()))

}
