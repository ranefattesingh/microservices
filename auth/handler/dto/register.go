package dto

// RegisterRequest holds incoming registration payload
type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,strong_password,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirm_password,omitempty" validate:"eqfield=Password"`
}
