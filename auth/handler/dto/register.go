package dto

import "github.com/ranefattesingh/microservices/auth/json"

// RegisterRequest holds incoming registration payload
type RegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (r RegisterRequest) Validate() []json.FieldError {
	var errs []json.FieldError
	if r.Email == "" {
		errs = append(errs, json.FieldError{Field: "email", Message: "email is required"})
	}
	if r.Password == "" {
		errs = append(errs, json.FieldError{Field: "password", Message: "password is required"})
	}
	if r.ConfirmPassword == "" {
		errs = append(errs, json.FieldError{Field: "confirm_password", Message: "confirm password is required"})
	}
	if r.Password != "" && r.ConfirmPassword != "" && r.Password != r.ConfirmPassword {
		errs = append(errs, json.FieldError{Field: "confirm_password", Message: "passwords do not match"})
	}

	return errs
}
