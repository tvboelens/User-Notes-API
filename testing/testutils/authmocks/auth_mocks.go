package authmocks

import (
	"context"
	"user-notes-api/auth"

	"github.com/stretchr/testify/mock"
)

type MockLoginManager struct {
	mock.Mock
}

type MockRegistrationManager struct {
	mock.Mock
}

func (m *MockLoginManager) LoginUser(ctx context.Context, credentials *auth.Credentials) (uint, bool, error) {
	args := m.Called(ctx, credentials)
	return uint(args.Int(0)), args.Bool(1), args.Error(2)
}

func (m *MockRegistrationManager) RegisterUser(ctx context.Context, credentials *auth.Credentials) (uint, error) {
	args := m.Called(ctx, credentials)
	return uint(args.Int(0)), args.Error(1)
}
