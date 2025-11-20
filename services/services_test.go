package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"user-notes-api/auth"
	"user-notes-api/models"
	"user-notes-api/testing/testutils"
	"user-notes-api/testing/testutils/repositorymocks"
)

func TestAuthServices(t *testing.T) {
	var username string = "Alice"
	var jwt_secret string = "jwt_secret"
	var password string = "secret_password"
	var wrong_pwd string = "wrong_password"

	creds := auth.Credentials{Username: username, Password: password}
	ctx := context.Background()

	user := models.User{Username: username, Password: password}
	repo := testutils.MockUserCreatorReader{User: &user, Registered: false}
	pwd_hasher := testutils.MockPwdHasher{Hash: []byte(password)}

	login_service := NewLoginService(&pwd_hasher, &repo, jwt_secret)

	// Login fails if user does not exist and we get a NotFound error
	token_string, err := login_service.Login(ctx, creds)

	assert.Error(t, err)
	assert.Equal(t, 0, len(token_string))
	var errNotFound *auth.ErrorNotFound
	assert.True(t, errors.As(err, &errNotFound))

	registration_service := NewRegistrationService(&pwd_hasher, &repo, jwt_secret)

	// First registration succesful and jwt token is not empty
	token_string, err = registration_service.Register(ctx, creds)
	assert.NoError(t, err)
	assert.True(t, len(token_string) > 0)

	// check the claims in the token
	token, err := jwt.Parse(token_string, func(token *jwt.Token) (any, error) {
		return []byte(jwt_secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	assert.NoError(t, err)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	issuer, err := claims.GetIssuer()
	assert.NoError(t, err)
	assert.Equal(t, "auth.user-notes-api.local", issuer)

	subject, err := claims.GetSubject()
	assert.NoError(t, err)
	assert.Equal(t, creds.Username, subject)

	issuedAt, err := claims.GetIssuedAt()
	assert.NoError(t, err)
	assert.True(t, time.Now().After(issuedAt.Time))

	expirationTime, err := claims.GetExpirationTime()
	assert.NoError(t, err)
	assert.True(t, expirationTime.After(time.Now()))

	// After registration login is possible
	token_string, err = login_service.Login(ctx, creds)
	assert.NoError(t, err)
	assert.True(t, len(token_string) > 0)

	// check the claims in the token
	token, err = jwt.Parse(token_string, func(token *jwt.Token) (any, error) {
		return []byte(jwt_secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	assert.NoError(t, err)

	claims, ok = token.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	issuer, err = claims.GetIssuer()
	assert.NoError(t, err)
	assert.Equal(t, "auth.user-notes-api.local", issuer)

	subject, err = claims.GetSubject()
	assert.NoError(t, err)
	assert.Equal(t, creds.Username, subject)

	issuedAt, err = claims.GetIssuedAt()
	assert.NoError(t, err)
	assert.True(t, time.Now().After(issuedAt.Time))

	expirationTime, err = claims.GetExpirationTime()
	assert.NoError(t, err)
	assert.True(t, expirationTime.After(time.Now()))

	// Registration fails if user already exists
	token_string, err = registration_service.Register(ctx, creds)
	assert.Error(t, err)
	assert.False(t, len(token_string) > 0)

	// Login fails with the wrong password
	wrong_creds := auth.Credentials{Username: username, Password: wrong_pwd}
	token_string, err = login_service.Login(ctx, wrong_creds)

	assert.Error(t, err)
	assert.Equal(t, 0, len(token_string))
}

func TestNoteServiceGetNoteSuccess(t *testing.T) {
	note_reader := new(repositorymocks.NoteReaderMock)
	note_creator := new(repositorymocks.NoteCreatorMock)
	user_repo := new(repositorymocks.UserRepoMock)

	note_service := NewNoteService(note_reader, note_creator, user_repo)

	noteId := uint(1)
	userId := uint(2)
	ctx := context.Background()
	note_reader.On("FindNoteById", ctx, noteId).Return(&models.Note{UserID: 2, Title: "Title", Body: "Content"}, nil)

	note, err := note_service.GetNote(ctx, noteId, userId)
	assert.NoError(t, err)
	assert.Equal(t, "Title", note.Title)
	assert.Equal(t, "Content", note.Content)
}

func TestNoteServiceGetNoteWrongOwner(t *testing.T) {
	note_reader := new(repositorymocks.NoteReaderMock)
	note_creator := new(repositorymocks.NoteCreatorMock)
	user_repo := new(repositorymocks.UserRepoMock)

	note_service := NewNoteService(note_reader, note_creator, user_repo)

	noteId := uint(1)
	userId := uint(2)
	ctx := context.Background()
	note_reader.On("FindNoteById", ctx, noteId).Return(&models.Note{UserID: 1, Title: "Title", Body: "Content"}, nil)

	_, err := note_service.GetNote(ctx, noteId, userId)
	assert.Error(t, err)
	var errWrongOwner *ErrorWrongOwner
	assert.True(t, errors.As(err, &errWrongOwner))
}

func TestNoteServiceGetNoteNotFound(t *testing.T) {
	note_reader := new(repositorymocks.NoteReaderMock)
	note_creator := new(repositorymocks.NoteCreatorMock)
	user_repo := new(repositorymocks.UserRepoMock)

	note_service := NewNoteService(note_reader, note_creator, user_repo)

	noteId := uint(1)
	userId := uint(2)
	ctx := context.Background()
	note_reader.On("FindNoteById", ctx, noteId).Return(&models.Note{}, errors.New("note not found"))

	_, err := note_service.GetNote(ctx, noteId, userId)
	assert.Error(t, err)
	var errNotFound *ErrorNoteNotFound
	assert.True(t, errors.As(err, &errNotFound))
}

func TestNoteServiceCreateNoteUser(t *testing.T) {
	note_reader := new(repositorymocks.NoteReaderMock)
	note_creator := new(repositorymocks.NoteCreatorMock)
	user_repo := new(repositorymocks.UserRepoMock)

	note_service := NewNoteService(note_reader, note_creator, user_repo)

	username := "Alice"
	password := "secret_password"
	note := Note{Title: "title", Content: "content"}
	ctx := context.Background()
	user_repo.On("FindUserByName", ctx, username).
		Return(&models.User{Model: gorm.Model{ID: 2}, Username: username, Password: password}, nil)

	note_model := models.Note{User: models.User{Model: gorm.Model{ID: 2}, Username: username, Password: password},
		UserID: 2, Title: note.Title, Body: note.Content}
	note_creator.On("CreateNote", ctx, &note_model).
		Run(func(args mock.Arguments) {
			note := args.Get(1).(*models.Note)
			note.ID = 4
		}).
		Return(nil)

	id, err := note_service.CreateNote(ctx, note, username)
	assert.NoError(t, err)
	assert.Equal(t, uint(4), id)

}

func TestNoteServiceCreateNoteUserNotFound(t *testing.T) {
	note_reader := new(repositorymocks.NoteReaderMock)
	note_creator := new(repositorymocks.NoteCreatorMock)
	user_repo := new(repositorymocks.UserRepoMock)

	note_service := NewNoteService(note_reader, note_creator, user_repo)

	username := "Alice"
	note := Note{Title: "title", Content: "content"}
	ctx := context.Background()
	user_repo.On("FindUserByName", ctx, username).Return(&models.User{}, errors.New("user not found"))

	_, err := note_service.CreateNote(ctx, note, username)
	assert.Error(t, err)
	var errNotFound *ErrorUserNotFound
	assert.True(t, errors.As(err, &errNotFound))
}
