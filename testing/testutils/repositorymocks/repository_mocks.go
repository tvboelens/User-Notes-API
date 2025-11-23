package repositorymocks

import (
	"context"

	"user-notes-api/models"

	"github.com/stretchr/testify/mock"
)

type NoteCreatorMock struct {
	mock.Mock
}

type NoteReaderMock struct {
	mock.Mock
}

type UserRepoMock struct {
	mock.Mock
}

func (m *NoteReaderMock) FindNoteById(ctx context.Context, id uint) (*models.Note, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Note), args.Error(1)
}
func (m *NoteReaderMock) FindNotesByUserId(ctx context.Context, userId uint) (*[]models.Note, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*[]models.Note), args.Error(1)
}

func (m *NoteCreatorMock) CreateNote(ctx context.Context, note *models.Note) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *UserRepoMock) FindUserById(ctx context.Context, id uint) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *UserRepoMock) FindUserByName(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*models.User), args.Error(1)
}
