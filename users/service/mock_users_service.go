package service

import (
	"context"

	"github.com/stretchr/testify/mock"
	handlersmodels "github.com/ranefattesingh/ecommerce-platform/users/handlers/models"
)

// MockUsersService provides testify-based expectations for handler tests.
type MockUsersService struct {
	mock.Mock
}

func NewMockUsersService() *MockUsersService {
	return &MockUsersService{}
}

func (m *MockUsersService) CreateUser(ctx context.Context, req handlersmodels.CreateUserRequest) (int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUsersService) GetUser(ctx context.Context, id int64) (handlersmodels.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(handlersmodels.User), args.Error(1)
}

func (m *MockUsersService) UpdateUser(ctx context.Context, id int64, req handlersmodels.UpdateUserRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockUsersService) DeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
