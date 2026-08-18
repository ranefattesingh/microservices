package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validate *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	v.RegisterValidation("strong_password", validateStrongPassword)

	return &Validator{
		Validate: v,
	}
}

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 5 {
		return false
	}

	return true
}
