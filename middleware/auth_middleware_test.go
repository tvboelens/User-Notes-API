package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"user-notes-api/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// test missing fields in jwt?

func TestAuthMiddlewareSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	jwt_secret := "jwt_secret"
	router.Use(JwtMiddleware(jwt_secret))

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, services.JwtClaims{
		UserId: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth.user-notes-api.local",
			Subject:   "Alice",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour))},
	})

	token_string, err := token.SignedString([]byte(jwt_secret))
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token_string)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")

}

func TestAuthMiddlewareTokenExpired(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	jwt_secret := "jwt_secret"
	router.Use(JwtMiddleware(jwt_secret))

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, services.JwtClaims{
		UserId: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth.user-notes-api.local",
			Subject:   "Alice",
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-8 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-8 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-4 * time.Hour))},
	})

	token_string, err := token.SignedString([]byte(jwt_secret))
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token_string)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "token is expired")
}

func TestAuthMiddlewareInvalidNotBefore(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	jwt_secret := "jwt_secret"
	router.Use(JwtMiddleware(jwt_secret))

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, services.JwtClaims{
		UserId: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth.user-notes-api.local",
			Subject:   "Alice",
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(0 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(4 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour))},
	})

	token_string, err := token.SignedString([]byte(jwt_secret))
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token_string)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "token is not valid yet")

}
