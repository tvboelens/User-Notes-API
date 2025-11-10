package testutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMocks(t *testing.T) {
	password := []byte("secret_password")
	wrong_pwd := []byte("wrong_password")

	var pwd_hasher MockPwdHasher

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

	var repo MockUserCreatorReader
	ctx := context.Background()
	_, err = repo.CreateUserByNameAndPassword(ctx, "Alice", string(password))
	assert.NoError(t, err)

	// registering the same user twice yields an error
	_, err = repo.CreateUserByNameAndPassword(ctx, "Alice", string(password))
	assert.Error(t, err)

	user, err := repo.FindUserByName(ctx, "Alice")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", user.Username)

	_, err = repo.FindUserByName(ctx, "Bob")
	assert.Error(t, err)
}
