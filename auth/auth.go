package auth

import (
	"context"
	"fmt"

	"user-notes-api/repositories"
	"user-notes-api/utils"
)

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

func LoginUser(ctx context.Context, credentials *Credentials, user_reader repositories.UserReader, pwd_comparer utils.PasswordComparer) (uint, bool, error) {
	user, err := user_reader.FindUserByName(ctx, credentials.Username)
	if err != nil {
		return 0, false, &ErrorNotFound{Username: credentials.Username, Err: err}
	}

	p, err := utils.ParseHashString(user.Password)
	if err != nil {
		return 0, false, fmt.Errorf("login user: could not parse hash string: %w", err)
	}

	isValid, err := pwd_comparer.Compare(p.Hash, p.Salt, []byte(credentials.Password))
	return user.ID, isValid, err
}

func RegisterUser(ctx context.Context, credentials *Credentials, user_creator repositories.UserCreator, pwd_hasher utils.PasswordHasher) (uint, error) {
	salt, err := pwd_hasher.GenerateSalt()
	if err != nil {
		return 0, fmt.Errorf("register user: could not generate salt: %w", err)
	}

	hash, err := pwd_hasher.GenerateHash([]byte(credentials.Password), salt)
	if err != nil {
		return 0, fmt.Errorf("register user: could not generate hash: %w", err)
	}

	p := utils.ParsedHashString{Id: "Argon2id", Version: 19, Hash: hash, Salt: salt}

	hash_string, err := utils.EncodeHashString(&p)
	if err != nil {
		return 0, fmt.Errorf("register user: failed to encode hash string: %w", err)
	}

	user, err := user_creator.CreateUserByNameAndPassword(ctx, credentials.Username, hash_string)

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
