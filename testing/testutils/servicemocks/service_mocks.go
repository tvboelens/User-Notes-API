package servicemocks

import (
	"context"

	"user-notes-api/auth"
	"user-notes-api/services"

	"github.com/stretchr/testify/mock"
)

type MockLoginService struct {
	mock.Mock
}

type MockRegistrationService struct {
	mock.Mock
}

type MockNoteModificationService struct {
	mock.Mock
}

type MockNoteReaderService struct {
	mock.Mock
}

func (m *MockLoginService) Login(ctx context.Context, credentials auth.Credentials) (string, error) {
	args := m.Called(ctx, credentials)
	return args.String(0), args.Error(1)
}

func (m *MockRegistrationService) Register(ctx context.Context, credentials auth.Credentials) (string, error) {
	args := m.Called(ctx, credentials)
	return args.String(0), args.Error(1)
}

func (m *MockNoteReaderService) GetNotes(ctx context.Context, userId uint) (services.GetNotesResult, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(services.GetNotesResult), args.Error(1)
}

func (m *MockNoteReaderService) GetNote(ctx context.Context, noteId uint, userId uint) (services.Note, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(services.Note), args.Error(1)
}

func (m *MockNoteModificationService) CreateNote(ctx context.Context, note services.Note, username string) (uint, error) {
	args := m.Called(ctx, note, username)
	return uint(args.Int(0)), args.Error(1)
}
func (m *MockNoteModificationService) UpdateNote(ctx context.Context, id uint, note services.Note) error {
	args := m.Called(ctx, id, note)
	return args.Error(0)
}
func (m *MockNoteModificationService) DeleteNote(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
