package dto

import (
	"github.com/ranefattesingh/microservices/user/service/models"
)

type UpdateUserRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	Email     string `json:"email" validate:"required,email"`
}

func (u *UpdateUserRequest) ToServiceModel() *models.User {
	return &models.User{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}
