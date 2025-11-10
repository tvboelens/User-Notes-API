package auth

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"user-notes-api/models"
	"user-notes-api/utils"
)

/*
type MockPwdComparer struct {
	Username string
	Password string
}

func (c *MockPwdComparer) Compare(hash, salt, password []byte) (bool, error) {
	pwd := []byte(c.Password)
	return bytes.Equal(pwd, password), nil
}
*/

type MockUserCreatorReader struct {
	User       *models.User
	Registered bool
}

type MockPwdHasher struct {
	Hash []byte
}

func (m *MockUserCreatorReader) CreateUser(ctx context.Context, user *models.User) error {
	ph := utils.ParsedHashString{Hash: []byte(user.Password)}
	hash_string, err := utils.EncodeHashString(&ph)
	if err != nil {
		return err
	}

	m.User = &models.User{Username: user.Username, Password: hash_string}
	m.Registered = true
	return nil
}

func (m *MockUserCreatorReader) CreateUserByNameAndPassword(ctx context.Context, username string, password string) (models.User, error) {
	ph := utils.ParsedHashString{Hash: []byte(password)}
	hash_string, err := utils.EncodeHashString(&ph)
	if err != nil {
		return models.User{Username: "fake_name", Password: ""}, err
	}
	m.User = &models.User{Username: username, Password: hash_string}
	m.Registered = true
	return *m.User, nil
}

func (m *MockUserCreatorReader) FindUserById(ctx context.Context, id uint) (models.User, error) {
	return models.User{Username: "fake_username", Password: ""}, errors.New("function not supported")
}

func (m *MockUserCreatorReader) FindUserByName(ctx context.Context, username string) (models.User, error) {
	if username == m.User.Username && m.Registered {
		return *m.User, nil
	}
	return models.User{Username: "fake_username", Password: ""}, errors.New("wrong user")
}

func (m *MockPwdHasher) GenerateHash(password, salt []byte) ([]byte, error) {
	return password, nil
}
func (m *MockPwdHasher) GenerateSalt() ([]byte, error) {
	var arr []byte
	return arr, nil
}
func (m *MockPwdHasher) Compare(hash, salt, password []byte) (bool, error) {
	return bytes.Equal(hash, password), nil
}

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
	repo := MockUserCreatorReader{User: &user, Registered: false}
	pwd_hasher := MockPwdHasher{Hash: []byte(password)}

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
