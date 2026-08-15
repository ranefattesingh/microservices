package dto

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/ranefattesingh/microservices/user/pkg"
	validatorWrap "github.com/ranefattesingh/microservices/user/validator"
)

type ModificationRequest interface {
	CreateUserRequest | UpdateUserRequest
}

func ValidateUser[T ModificationRequest](v *validatorWrap.Validator, req T) []pkg.FieldError {
	parse := func(fe validator.FieldError) string {
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

	err := v.Validate.Struct(req)
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
