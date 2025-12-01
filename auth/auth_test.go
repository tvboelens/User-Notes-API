package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"user-notes-api/models"
	"user-notes-api/testing/testutils"
)

func TestRegisterAndLogin(t *testing.T) {
	var username string = "Alice"
	var password string = "secret_password"
	var wrong_pwd string = "wrong_password"

	creds := Credentials{Username: username, Password: password}

	user := models.User{Username: username, Password: password}
	repo := &testutils.MockUserCreatorReader{User: &user, Registered: false}
	pwd_hasher := &testutils.MockPwdHasher{Hash: []byte(password)}

	login_manager := NewLoginManager(repo, pwd_hasher)
	registration_manager := NewRegistrationManager(repo, pwd_hasher)

	// Get NotFoundError if user not registered
	ctx := context.Background()
	_, _, err := login_manager.LoginUser(ctx, &creds)
	assert.Error(t, err)
	var errNotFound *ErrorNotFound
	assert.True(t, errors.As(err, &errNotFound))

	id, err := registration_manager.RegisterUser(ctx, &creds)
	assert.NoError(t, err)
	assert.True(t, id > 0)

	// Login succesful with correct credentials
	_, logged_in, err := login_manager.LoginUser(ctx, &creds)
	assert.NoError(t, err)
	assert.True(t, logged_in)

	// Login unsuccesful with incorrect credentials
	creds_wrong_pwd := Credentials{Username: username, Password: wrong_pwd}
	_, isValid, err := login_manager.LoginUser(ctx, &creds_wrong_pwd)
	assert.NoError(t, err)
	assert.False(t, isValid)
}
