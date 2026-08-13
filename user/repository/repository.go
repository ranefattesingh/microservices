package repository

import (
	"context"

	"github.com/ranefattesingh/microservices/user/service/models"
)

var _ UserRepository = (*userRepository)(nil)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (int64, error)
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
