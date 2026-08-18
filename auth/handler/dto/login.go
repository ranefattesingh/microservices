package dto

// LoginRequest holds login payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"  validate:"required"`
}
