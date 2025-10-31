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

type MockUserRepository struct {
	Username string
	Password string
}

type MockPwdComparer struct {
	Username string
	Password string
}

func (r *MockUserRepository) FindUserById(ctx context.Context, id uint) (models.User, error) {
	return models.User{Username: "fake_name", Password: ""}, errors.New("Unsupported function")
}

func (r *MockUserRepository) FindUserByName(ctx context.Context, username string) (models.User, error) {
	if r.Username != username {
		return models.User{Username: username, Password: ""}, errors.New("user does not exist")
	}

	ph := utils.ParsedHashString{Hash: []byte(r.Password)}

	hash_string, err := utils.EncodeHashString(&ph)
	if err != nil {
		return models.User{Username: "fake_name", Password: ""}, errors.New("Failed to encode hash string")
	}

	return models.User{Username: r.Username, Password: hash_string}, nil
}

func (c *MockPwdComparer) Compare(hash, salt, password []byte) (bool, error) {
	pwd := []byte(c.Password)
	return bytes.Equal(pwd, password), nil
}

type MockUserCreatorReader struct {
	User       *models.User
	Registered bool
}

func (m *MockUserCreatorReader) CreateUser(ctx context.Context, user *models.User) error {
	m.User = user
	m.Registered = true
	return nil
}

func (m *MockUserCreatorReader) CreateUserByNameAndPassword(ctx context.Context, username string, password string) (models.User, error) {
	m.User = &models.User{Username: username, Password: password}
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

type MockPwdHasher struct {
	Hash []byte
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

func TestLogin(t *testing.T) {
	var username string = "Alice"
	var password string = "secret_password"
	var wrong_pwd string = "wrong_password"

	creds := Credentials{Username: username, Password: password}

	comparer := MockPwdComparer{Username: username, Password: password}
	repo := MockUserRepository{Username: username, Password: password}

	ctx := context.Background()

	isValid, err := LoginUser(ctx, &creds, &repo, &comparer)
	assert.NoError(t, err)
	assert.True(t, isValid)

	creds_wrong_pwd := Credentials{Username: username, Password: wrong_pwd}
	isValid, err = LoginUser(ctx, &creds_wrong_pwd, &repo, &comparer)
	assert.NoError(t, err)
	assert.False(t, isValid)

	creds_no_username := Credentials{Username: "Bob", Password: wrong_pwd}
	_, err = LoginUser(ctx, &creds_no_username, &repo, &comparer)
	assert.Error(t, err)
}

func TestRegister(t *testing.T) {
	var username string = "Alice"
	var password string = "secret_password"

	creds := Credentials{Username: username, Password: password}

	user := models.User{Username: username, Password: password}
	repo := MockUserCreatorReader{User: &user, Registered: false}
	pwd_hasher := MockPwdHasher{Hash: []byte(password)}

	ctx := context.Background()
	err := RegisterUser(ctx, &creds, &repo, &pwd_hasher)

	assert.NoError(t, err)

}
