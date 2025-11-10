package testutils

import (
	"bytes"
	"context"
	"errors"

	"user-notes-api/models"
	"user-notes-api/utils"
)

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

func (m *MockUserCreatorReader) CreateUserByNameAndPassword(ctx context.Context, username string, password string) (*models.User, error) {
	if m.Registered {
		msg := "username " + username + " already exists"
		return nil, errors.New(msg)
	}

	ph := utils.ParsedHashString{Hash: []byte(password)}
	hash_string, err := utils.EncodeHashString(&ph)
	if err != nil {
		return nil, err
	}
	m.User = &models.User{Username: username, Password: hash_string}
	m.Registered = true
	return m.User, nil
}

func (m *MockUserCreatorReader) FindUserById(ctx context.Context, id uint) (*models.User, error) {
	return nil, errors.New("function not supported")
}

func (m *MockUserCreatorReader) FindUserByName(ctx context.Context, username string) (*models.User, error) {
	if username == m.User.Username && m.Registered {
		return m.User, nil
	}
	return nil, errors.New("wrong user")
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
