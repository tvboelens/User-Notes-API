package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"user-notes-api/models"
	"user-notes-api/testing/testutils"
)

func TestMocks(t *testing.T) {
	password := []byte("secret_password")
	wrong_pwd := []byte("wrong_password")

	var pwd_hasher testutils.MockPwdHasher

	salt, err := pwd_hasher.GenerateSalt()
	assert.NoError(t, err)

	hash, err := pwd_hasher.GenerateHash(password, salt)
	assert.NoError(t, err)

	isValid, err := pwd_hasher.Compare(hash, salt, password)
	assert.NoError(t, err)
	assert.True(t, isValid)

	hash, err = pwd_hasher.GenerateHash(wrong_pwd, salt)
	assert.NoError(t, err)

	isValid, err = pwd_hasher.Compare(hash, salt, password)
	assert.NoError(t, err)
	assert.False(t, isValid)

	var repo testutils.MockUserCreatorReader
	ctx := context.Background()
	_, err = repo.CreateUserByNameAndPassword(ctx, "Alice", string(password))
	assert.NoError(t, err)

	user, err := repo.FindUserByName(ctx, "Alice")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", user.Username)

	_, err = repo.FindUserByName(ctx, "Bob")
	assert.Error(t, err)
}

func TestRegisterAndLogin(t *testing.T) {
	var username string = "Alice"
	var password string = "secret_password"
	var wrong_pwd string = "wrong_password"

	creds := Credentials{Username: username, Password: password}

	user := models.User{Username: username, Password: password}
	repo := testutils.MockUserCreatorReader{User: &user, Registered: false}
	pwd_hasher := testutils.MockPwdHasher{Hash: []byte(password)}

	// Get NotFoundError if user not registered
	ctx := context.Background()
	_, err := LoginUser(ctx, &creds, &repo, &pwd_hasher)
	assert.Error(t, err)
	var errNotFound *ErrorNotFound
	assert.True(t, errors.As(err, &errNotFound))

	err = RegisterUser(ctx, &creds, &repo, &pwd_hasher)
	assert.NoError(t, err)

	// Login succesful with correct credentials
	logged_in, err := LoginUser(ctx, &creds, &repo, &pwd_hasher)
	assert.NoError(t, err)
	assert.True(t, logged_in)

	// Login unsuccesful with incorrect credentials
	creds_wrong_pwd := Credentials{Username: username, Password: wrong_pwd}
	isValid, err := LoginUser(ctx, &creds_wrong_pwd, &repo, &pwd_hasher)
	assert.NoError(t, err)
	assert.False(t, isValid)
}
