package auth

import (
	"context"
	"fmt"

	"user-notes-api/repositories"
	"user-notes-api/utils"
)

type LoginManagerIfc interface {
	LoginUser(ctx context.Context, credentials *Credentials) (uint, bool, error)
}

type RegistrationManagerIfc interface {
	RegisterUser(ctx context.Context, credentials *Credentials) (uint, error)
}

type LoginManager struct {
	UserReader  repositories.UserReader
	PwdComparer utils.PasswordComparer
}

type RegistrationManager struct {
	UserCreator repositories.UserCreator
	PwdHasher   utils.PasswordHasher
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ErrorNotFound struct {
	Username string
	Err      error
}

func (e *ErrorNotFound) Error() string {
	return fmt.Sprintf("user %q not found: %v", e.Username, e.Err)
}

func (e *ErrorNotFound) Unwrap() error {
	return e.Err
}

func NewLoginManager(user_reader repositories.UserReader, pwd_comparer utils.PasswordComparer) *LoginManager {
	login_manager := LoginManager{UserReader: user_reader, PwdComparer: pwd_comparer}
	return &login_manager
}

func NewRegistrationManager(user_creator repositories.UserCreator, pwd_hasher utils.PasswordHasher) *RegistrationManager {
	registration_manager := RegistrationManager{UserCreator: user_creator, PwdHasher: pwd_hasher}
	return &registration_manager
}

func (m *LoginManager) LoginUser(ctx context.Context, credentials *Credentials) (uint, bool, error) {
	user, err := m.UserReader.FindUserByName(ctx, credentials.Username)
	if err != nil {
		return 0, false, &ErrorNotFound{Username: credentials.Username, Err: err}
	}

	p, err := utils.ParseHashString(user.Password)
	if err != nil {
		return 0, false, fmt.Errorf("login user: could not parse hash string: %w", err)
	}

	isValid, err := m.PwdComparer.Compare(p.Hash, p.Salt, []byte(credentials.Password))
	return user.ID, isValid, err
}

func (m *RegistrationManager) RegisterUser(ctx context.Context, credentials *Credentials) (uint, error) {
	salt, err := m.PwdHasher.GenerateSalt()
	if err != nil {
		return 0, fmt.Errorf("register user: could not generate salt: %w", err)
	}

	hash, err := m.PwdHasher.GenerateHash([]byte(credentials.Password), salt)
	if err != nil {
		return 0, fmt.Errorf("register user: could not generate hash: %w", err)
	}

	p := utils.ParsedHashString{Id: "Argon2id", Version: 19, Hash: hash, Salt: salt}

	hash_string, err := utils.EncodeHashString(&p)
	if err != nil {
		return 0, fmt.Errorf("register user: failed to encode hash string: %w", err)
	}

	user, err := m.UserCreator.CreateUserByNameAndPassword(ctx, credentials.Username, hash_string)

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
