package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"user-notes-api/auth"
	"user-notes-api/models"
	"user-notes-api/testing/testutils"
)

func TestAuthServices(t *testing.T) {
	var username string = "Alice"
	var jwt_secret string = "jwt_secret"
	var password string = "secret_password"
	var wrong_pwd string = "wrong_password"

	creds := auth.Credentials{Username: username, Password: password}
	ctx := context.Background()

	user := models.User{Username: username, Password: password}
	repo := testutils.MockUserCreatorReader{User: &user, Registered: false}
	pwd_hasher := testutils.MockPwdHasher{Hash: []byte(password)}

	login_service := NewLoginService(&pwd_hasher, &repo)

	// Login fails if user does not exist and we get a NotFound error
	token_string, err := login_service.Login(ctx, jwt_secret, creds)

	assert.Error(t, err)
	assert.Equal(t, 0, len(token_string))
	var errNotFound *auth.ErrorNotFound
	assert.True(t, errors.As(err, &errNotFound))

	registration_service := NewRegistrationService(&pwd_hasher, &repo)

	// First registration succesful and jwt token is not empty
	token_string, err = registration_service.Register(ctx, jwt_secret, creds)
	assert.NoError(t, err)
	assert.True(t, len(token_string) > 0)

	// check the claims in the token
	token, err := jwt.Parse(token_string, func(token *jwt.Token) (any, error) {
		return []byte(jwt_secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	assert.NoError(t, err)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	issuer, err := claims.GetIssuer()
	assert.NoError(t, err)
	assert.Equal(t, "auth.user-notes-api.local", issuer)

	subject, err := claims.GetSubject()
	assert.NoError(t, err)
	assert.Equal(t, creds.Username, subject)

	issuedAt, err := claims.GetIssuedAt()
	assert.NoError(t, err)
	assert.True(t, time.Now().After(issuedAt.Time))

	expirationTime, err := claims.GetExpirationTime()
	assert.NoError(t, err)
	assert.True(t, expirationTime.After(time.Now()))

	// After registration login is possible
	token_string, err = login_service.Login(ctx, jwt_secret, creds)
	assert.NoError(t, err)
	assert.True(t, len(token_string) > 0)

	// check the claims in the token
	token, err = jwt.Parse(token_string, func(token *jwt.Token) (any, error) {
		return []byte(jwt_secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	assert.NoError(t, err)

	claims, ok = token.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	issuer, err = claims.GetIssuer()
	assert.NoError(t, err)
	assert.Equal(t, "auth.user-notes-api.local", issuer)

	subject, err = claims.GetSubject()
	assert.NoError(t, err)
	assert.Equal(t, creds.Username, subject)

	issuedAt, err = claims.GetIssuedAt()
	assert.NoError(t, err)
	assert.True(t, time.Now().After(issuedAt.Time))

	expirationTime, err = claims.GetExpirationTime()
	assert.NoError(t, err)
	assert.True(t, expirationTime.After(time.Now()))

	// Registration fails if user already exists
	token_string, err = registration_service.Register(ctx, jwt_secret, creds)
	assert.Error(t, err)
	assert.False(t, len(token_string) > 0)

	// Login fails with the wrong password
	wrong_creds := auth.Credentials{Username: username, Password: wrong_pwd}
	token_string, err = login_service.Login(ctx, jwt_secret, wrong_creds)

	assert.Error(t, err)
	assert.Equal(t, 0, len(token_string))
}
