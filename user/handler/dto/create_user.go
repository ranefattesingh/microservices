package dto

import (
	"github.com/ranefattesingh/microservices/user/service/models"
)

type CreateUserRequest struct {
	FirstName       string `json:"first_name" validate:"required,min=2"`
	LastName        string `json:"last_name" validate:"required,min=2"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,strong_password,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirm_password,omitempty" validate:"eqfield=Password"`
}

func (u *CreateUserRequest) ToServiceModel() *models.User {
	return &models.User{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
	}
}
