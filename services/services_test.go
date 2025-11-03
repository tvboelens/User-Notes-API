package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"user-notes-api/auth"
	"user-notes-api/models"
	"user-notes-api/testing/testutils"
)

/*
	Mocks for
		1. utils.PasswordComparer
		2. repositories.UserReader
		3. utils.PasswordHasher
		4. repositories.UserCreator

	I could simply use the mocks from the auth tests?

	Login:
		1. error if user does not exist
		2. error if wrong password
		3. in both these case need the jwt to be empty
		4. success with correct credentials, need to check validity of jwt
			1. subject is user(name)
			2. Issuer not empty
			3. Issued at and expires at have meaningful (?) values (difference of 4 hours and issued at is somewhere in the last minute)
			4. signature is correct

	Register
		1. error if user already exists -> empty jwt
		2. success if not
		3.	validity of jwt as above
*/

func TestAuthServices(t *testing.T) {
	var username string = "Alice"
	var jwt_secret string = "jwt_secret"
	var password string = "secret_password"
	//var wrong_pwd string = "wrong_password"

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

	// After registration login is possible
	token_string, err = login_service.Login(ctx, jwt_secret, creds)
	assert.NoError(t, err)
	assert.True(t, len(token_string) > 0)

	// Registration fails if user already exists
	token_string, err = registration_service.Register(ctx, jwt_secret, creds)
	assert.Error(t, err)
	assert.False(t, len(token_string) > 0)
}
