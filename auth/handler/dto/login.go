package dto

import "github.com/ranefattesingh/microservices/auth/json"

// LoginRequest holds login payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l LoginRequest) Validate() []json.FieldError {
	var errs []json.FieldError
	if l.Email == "" {
		errs = append(errs, json.FieldError{Field: "email", Message: "email is required"})
	}
	if l.Password == "" {
		errs = append(errs, json.FieldError{Field: "password", Message: "password is required"})
	}

	return errs
}
