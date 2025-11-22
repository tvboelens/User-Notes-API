package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"user-notes-api/services"
)

func JwtMiddleware(jwt_secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header format must be Bearer {token}"})
			return
		}

		token_string := strings.TrimPrefix(header, "Bearer ")
		if token_string == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty token"})
			return
		}

		token, err := jwt.ParseWithClaims(token_string, &services.JwtClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(jwt_secret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "failed to parse token: " + err.Error()})
			return
		}

		claims, ok := token.Claims.(*services.JwtClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		notBefore, err := claims.GetNotBefore()

		if err != nil || notBefore == nil || !time.Now().After(notBefore.Time) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is not valid yet"})
			return
		}

		issuer, err := claims.GetIssuer()
		if err != nil || issuer != "auth.user-notes-api.local" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid issuer"})
			return
		}

		subject, err := claims.GetSubject() //username
		if err != nil || subject == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid subject"})
			return
		}

		c.Set("username", subject)

		expirationTime, err := claims.GetExpirationTime()
		if err != nil || expirationTime == nil || time.Now().After(expirationTime.Time) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
			return
		}

		user_id, err := claims.GetUserId()
		if err != nil || user_id == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
			return
		}
		c.Set("user_id", user_id)

		c.Next()
	}
}
