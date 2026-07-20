package mock

import (
	"context"

	"github.com/ranefattesingh/ecommerce-platform/auth/repository/models"
	"github.com/stretchr/testify/mock"
)

type mockUsersRepository struct {
	mock.Mock
}

func NewMockUsersRepository() *mockUsersRepository {
	return new(mockUsersRepository)
}

func (m *mockUsersRepository) CreateUser(ctx context.Context, user models.Auth) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockUsersRepository) GetUsersHavingEmailOrPhone(ctx context.Context, email, phone string) ([]models.Auth, error) {
	args := m.Called(ctx, email, phone)
	return args.Get(0).([]models.Auth), args.Error(1)
}

func (m *mockUsersRepository) GetUser(ctx context.Context, id int64) (models.Auth, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Auth), args.Error(1)
}

func (m *mockUsersRepository) UpdateUser(ctx context.Context, user models.Auth) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUsersRepository) DeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
