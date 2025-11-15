package servicemocks

import (
	"context"

	"user-notes-api/auth"

	"github.com/stretchr/testify/mock"
)

type MockLoginService struct {
	mock.Mock
}

type MockRegistrationService struct {
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
