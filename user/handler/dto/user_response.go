package dto

import (
	"time"

	"github.com/ranefattesingh/microservices/user/service/models"
)

type UserResponse struct {
	ID        int64      `json:"id"`
	FirstName string     `json:"first_name" `
	LastName  string     `json:"last_name" `
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func FromServiceModel(user *models.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: &user.CreatedAt,
		UpdatedAt: &user.UpdatedAt,
	}
}
