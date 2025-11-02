package services

import (
	"context"
	"fmt"
	"time"
	"user-notes-api/auth"
	"user-notes-api/repositories"
	"user-notes-api/utils"

	"github.com/golang-jwt/jwt/v5"
)

type LoginService struct {
	password_comparer utils.PasswordComparer
	user_repo         repositories.UserReader
}

type RegistrationService struct {
	password_hasher utils.PasswordHasher
	user_repo       repositories.UserCreator
}

type ErrorWrongPassword struct {
	Username string
}

func (e *ErrorWrongPassword) Error() string {
	return fmt.Sprintf("wrong password for user %q:", e.Username)
}

func (s *LoginService) Login(ctx context.Context, jwt_secret string, credentials auth.Credentials) (string, error) {
	isValid, err := auth.LoginUser(ctx, &credentials, s.user_repo, s.password_comparer)
	if err != nil {
		return "", err
	}

	if !isValid {
		myErr := ErrorWrongPassword{Username: credentials.Username}
		return "", &myErr
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "auth.user-notes-api.local",
		Subject:   credentials.Username,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour)),
	})

	token_string, err := token.SignedString(jwt_secret)

	return token_string, err
}

func (s *RegistrationService) Register(ctx context.Context, jwt_secret string, credentials auth.Credentials) (string, error) {
	err := auth.RegisterUser(ctx, &credentials, s.user_repo, s.password_hasher)

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "auth.user-notes-api.local",
		Subject:   credentials.Username,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour)),
	})

	token_string, err := token.SignedString(jwt_secret)

	return token_string, err
}
