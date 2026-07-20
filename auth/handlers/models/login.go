package models

type LoginRequest struct {
	Email    string `json:"email" validate:"required,min=10"`
	Password string `json:"password" validate:"required,min=8"`
}
