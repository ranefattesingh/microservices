package repository

import (
	"context"

	"github.com/ranefattesingh/microservices/user/service/models"
)

var _ UserRepository = (*userRepository)(nil)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetUser(ctx context.Context, id int64) (*models.User, error)
}

type userRepository struct {
	// Add your database connection or any other dependencies here
	// e.g., *sql.DB, *gorm.DB, etc.
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	// TODO: Implement database insert logic here
	return 0, nil
}

func (r *userRepository) GetUser(ctx context.Context, id int64) (*models.User, error) {
	// TODO: Implement database query logic here
	return nil, nil
}
