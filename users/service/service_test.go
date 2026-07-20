package service_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	handlersmodels "github.com/ranefattesingh/ecommerce-platform/users/handlers/models"
	"github.com/ranefattesingh/ecommerce-platform/users/repository"
	mockRepo "github.com/ranefattesingh/ecommerce-platform/users/repository/mock"
	dbmodels "github.com/ranefattesingh/ecommerce-platform/users/repository/models"
	"github.com/ranefattesingh/ecommerce-platform/users/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_Success(t *testing.T) {
	t.Parallel()

	repo := mockRepo.NewMockUsersRepository()
	repo.On("GetUsersHavingEmailOrPhone", mock.Anything, "ada@example.com", "1234567890").Return([]dbmodels.User{}, nil).Once()
	repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user dbmodels.User) bool {
		return user.Email == "ada@example.com" && user.Phone == "1234567890"
	})).Return(int64(42), nil).Once()

	svc := service.NewUsersService(repo)
	id, err := svc.CreateUser(context.Background(), handlersmodels.CreateUserRequest{
		FirstName:  "Ada",
		LastName:   "Lovelace",
		Email:      "ada@example.com",
		Phone:      "1234567890",
		AccessType: handlersmodels.AccessType("admin"),
	})

	require.NoError(t, err)
	assert.Equal(t, int64(42), id)
	repo.AssertExpectations(t)
}

func TestCreateUser_ReturnsConflictForDuplicateEmailOrPhone(t *testing.T) {
	t.Parallel()

	repo := mockRepo.NewMockUsersRepository()
	repo.On("GetUsersHavingEmailOrPhone", mock.Anything, "ada@example.com", "1234567890").Return([]dbmodels.User{{
		Email: "ada@example.com",
		Phone: "1234567890",
	}}, nil).Once()

	svc := service.NewUsersService(repo)
	_, err := svc.CreateUser(context.Background(), handlersmodels.CreateUserRequest{
		FirstName:  "Ada",
		LastName:   "Lovelace",
		Email:      "ada@example.com",
		Phone:      "1234567890",
		AccessType: handlersmodels.AccessType("admin"),
	})

	require.Error(t, err)
	var svcErr *service.Error
	require.ErrorAs(t, err, &svcErr)
	assert.Equal(t, service.Conflict, svcErr.Code)
	assert.Len(t, svcErr.Violations, 2)
	repo.AssertExpectations(t)
}

func TestGetUser_Success(t *testing.T) {
	t.Parallel()

	repo := mockRepo.NewMockUsersRepository()
	repo.On("GetUser", mock.Anything, int64(7)).Return(dbmodels.User{
		ID:         7,
		FirstName:  "Grace",
		LastName:   "Hopper",
		Email:      "grace@example.com",
		Phone:      "5550001",
		AccessType: dbmodels.AccessType("admin"),
	}, nil).Once()

	svc := service.NewUsersService(repo)
	user, err := svc.GetUser(context.Background(), 7)

	require.NoError(t, err)
	assert.Equal(t, int64(7), user.ID)
	assert.Equal(t, "grace@example.com", user.Email)
	assert.Equal(t, handlersmodels.AccessType("admin"), user.AccessType)
	repo.AssertExpectations(t)
}

func TestGetUser_ReturnsNotFoundWhenRepositoryHasNoRows(t *testing.T) {
	t.Parallel()

	repo := mockRepo.NewMockUsersRepository()
	repo.On("GetUser", mock.Anything, int64(7)).Return(dbmodels.User{}, pgx.ErrNoRows).Once()

	svc := service.NewUsersService(repo)
	_, err := svc.GetUser(context.Background(), 7)

	require.Error(t, err)
	var svcErr *service.Error
	require.ErrorAs(t, err, &svcErr)
	assert.Equal(t, service.NotFound, svcErr.Code)
	repo.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	t.Parallel()

	repo := mockRepo.NewMockUsersRepository()
	repo.On("GetUser", mock.Anything, int64(11)).Return(dbmodels.User{ID: 11, AccessType: dbmodels.AccessType("member")}, nil).Once()
	repo.On("GetUsersHavingEmailOrPhone", mock.Anything, "new@example.com", "222").Return([]dbmodels.User{{ID: 11, Email: "new@example.com", Phone: "222"}}, nil).Once()
	repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(user dbmodels.User) bool {
		return user.ID == 11 && user.FirstName == "Linus" && user.LastName == "Torvalds"
	})).Return(nil).Once()

	svc := service.NewUsersService(repo)
	err := svc.UpdateUser(context.Background(), 11, handlersmodels.UpdateUserRequest{
		FirstName: "Linus",
		LastName:  "Torvalds",
		Email:     "new@example.com",
		Phone:     "222",
	})

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUpdateUser_ReturnsConflictForTakenEmailOrPhone(t *testing.T) {
	t.Parallel()

	repo := mockRepo.NewMockUsersRepository()
	repo.On("GetUser", mock.Anything, int64(11)).Return(dbmodels.User{ID: 11}, nil).Once()
	repo.On("GetUsersHavingEmailOrPhone", mock.Anything, "taken@example.com", "333").Return([]dbmodels.User{{ID: 22, Email: "taken@example.com", Phone: "333"}}, nil).Once()

	svc := service.NewUsersService(repo)
	err := svc.UpdateUser(context.Background(), 11, handlersmodels.UpdateUserRequest{
		Email: "taken@example.com",
		Phone: "333",
	})

	require.Error(t, err)
	var svcErr *service.Error
	require.ErrorAs(t, err, &svcErr)
	assert.Equal(t, service.Conflict, svcErr.Code)
	repo.AssertExpectations(t)
}

func TestUpdateUser_ReturnsNoUserUpdatedWhenRepositoryReportsNoRows(t *testing.T) {
	t.Parallel()

	repo := mockRepo.NewMockUsersRepository()
	repo.On("GetUser", mock.Anything, int64(11)).Return(dbmodels.User{ID: 11}, nil).Once()
	repo.On("GetUsersHavingEmailOrPhone", mock.Anything, "", "").Return([]dbmodels.User{}, nil).Once()
	repo.On("UpdateUser", mock.Anything, mock.Anything).Return(repository.ErrNoRowsUpdated).Once()

	svc := service.NewUsersService(repo)
	err := svc.UpdateUser(context.Background(), 11, handlersmodels.UpdateUserRequest{})

	require.ErrorIs(t, err, service.ErrNoUserUpdated)
	repo.AssertExpectations(t)
}

func TestDeleteUser_ReturnsNotFoundWhenRepositoryReportsNoRows(t *testing.T) {
	t.Parallel()

	repo := mockRepo.NewMockUsersRepository()
	repo.On("DeleteUser", mock.Anything, int64(99)).Return(repository.ErrNoRowsUpdated).Once()

	svc := service.NewUsersService(repo)
	err := svc.DeleteUser(context.Background(), 99)

	require.Error(t, err)
	var svcErr *service.Error
	require.ErrorAs(t, err, &svcErr)
	assert.Equal(t, service.NotFound, svcErr.Code)
	repo.AssertExpectations(t)
}
