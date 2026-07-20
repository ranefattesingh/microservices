package models

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,min=10"`
	Phone           string `json:"phone" validate:"required,numeric,min=10"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"eqfield=Password"`
}
