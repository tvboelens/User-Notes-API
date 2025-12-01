package services

import (
	"context"
	"fmt"
	"time"
	"user-notes-api/auth"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	UserId uint `json:"user_id,omitempty"`
	jwt.RegisteredClaims
}

type LoginServiceIfc interface {
	Login(ctx context.Context, credentials auth.Credentials) (string, error)
}

type RegistrationServiceIfc interface {
	Register(ctx context.Context, credentials auth.Credentials) (string, error)
}

type LoginService struct {
	LoginManager auth.LoginManagerIfc
	jwt_secret   string
}

type RegistrationService struct {
	RegistrationManager auth.RegistrationManagerIfc
	jwt_secret          string
}

type ErrorWrongPassword struct {
	Username string
}

func (e *ErrorWrongPassword) Error() string {
	return fmt.Sprintf("wrong password for user %q:", e.Username)
}

func (c *JwtClaims) GetUserId() (uint, error) {
	return c.UserId, nil
}

func NewLoginService(login_manager auth.LoginManagerIfc, jwt_secret string) *LoginService {
	login_service := LoginService{LoginManager: login_manager, jwt_secret: jwt_secret}
	return &login_service
}

func NewRegistrationService(registration_manager auth.RegistrationManagerIfc, jwt_secret string) *RegistrationService {
	registration_service := RegistrationService{RegistrationManager: registration_manager, jwt_secret: jwt_secret}
	return &registration_service
}

func (s *LoginService) Login(ctx context.Context, credentials auth.Credentials) (string, error) {
	user_id, isValid, err := s.LoginManager.LoginUser(ctx, &credentials)
	if err != nil {
		return "", err
	}

	if !isValid {
		myErr := ErrorWrongPassword{Username: credentials.Username}
		return "", &myErr
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		UserId: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth.user-notes-api.local",
			Subject:   credentials.Username,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour)),
		},
	})

	token_string, err := token.SignedString([]byte(s.jwt_secret))

	return token_string, err
}

func (s *RegistrationService) Register(ctx context.Context, credentials auth.Credentials) (string, error) {
	user_id, err := s.RegistrationManager.RegisterUser(ctx, &credentials)

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		UserId: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth.user-notes-api.local",
			Subject:   credentials.Username,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour)),
		},
	})

	token_string, err := token.SignedString([]byte(s.jwt_secret))

	return token_string, err
}
