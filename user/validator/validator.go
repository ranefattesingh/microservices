package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/ranefattesingh/microservices/user/handler/dto"
	"github.com/ranefattesingh/microservices/user/pkg"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	v.RegisterValidation("strong_password", validateStrongPassword)

	return &Validator{
		validate: v,
	}
}

func (v *Validator) ValidateCreateUser(req dto.CreateUserRequest) []pkg.FieldError {
	err := v.validate.Struct(req)
	if err != nil {
		if ve, ok := errors.AsType[validator.ValidationErrors](err); ok {
			fields := make([]pkg.FieldError, 0, len(ve))

			for _, ve := range ve {
				fields = append(fields, pkg.FieldError{
					Field:   ve.Field(),
					Message: parse(ve),
				})
			}

			return fields
		}
	}

	return nil
}

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 5 {
		return false
	}

	return true
}

func parse(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return "is too short"
	case "max":
		return "is too long"
	case "strong_password":
		return "password is too weak"
	case "eqfield":
		return "passwords doesn't match"
	}

	return "is invalid"
}
