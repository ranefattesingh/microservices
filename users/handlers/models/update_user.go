package models

type UpdateUserRequest struct {
	FirstName string `json:"firstName" validate:"required,min=3"`
	LastName  string `json:"lastName" validate:"required,min=3"`
	Email     string `json:"email" validate:"required,min=10"`
	Phone     string `json:"phone" validate:"required,min=10"`
}
