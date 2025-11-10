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
	Password_comparer utils.PasswordComparer
	User_repo         repositories.UserReader
	jwt_secret        string
}

type RegistrationService struct {
	Password_hasher utils.PasswordHasher
	User_repo       repositories.UserCreator
	jwt_secret      string
}

type NoteService struct {
	UserRepo    repositories.UserReader
	NoteCreator repositories.NoteCreator
	NoteReader  repositories.NoteReader
}

type ErrorWrongPassword struct {
	Username string
}

func (e *ErrorWrongPassword) Error() string {
	return fmt.Sprintf("wrong password for user %q:", e.Username)
}

func NewLoginService(password_comparer utils.PasswordComparer, user_repo repositories.UserReader, jwt_secret string) *LoginService {
	login_service := LoginService{Password_comparer: password_comparer, User_repo: user_repo, jwt_secret: jwt_secret}
	return &login_service
}

func NewRegistrationService(password_hasher utils.PasswordHasher, user_repo repositories.UserCreator, jwt_secret string) *RegistrationService {
	registration_service := RegistrationService{Password_hasher: password_hasher, User_repo: user_repo, jwt_secret: jwt_secret}
	return &registration_service
}

func NewNoteService(note_reader repositories.NoteReader, note_creator repositories.NoteCreator, user_repo repositories.UserReader) *NoteService {
	note_service := NoteService{NoteReader: note_reader, NoteCreator: note_creator, UserRepo: user_repo}
	return &note_service
}

func (s *LoginService) Login(ctx context.Context, credentials auth.Credentials) (string, error) {
	isValid, err := auth.LoginUser(ctx, &credentials, s.User_repo, s.Password_comparer)
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

	token_string, err := token.SignedString([]byte(s.jwt_secret))

	return token_string, err
}

func (s *RegistrationService) Register(ctx context.Context, credentials auth.Credentials) (string, error) {
	err := auth.RegisterUser(ctx, &credentials, s.User_repo, s.Password_hasher)

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "auth.user-notes-api.local",
		Subject:   credentials.Username,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour)),
	})

	token_string, err := token.SignedString([]byte(s.jwt_secret))

	return token_string, err
}
