package dto

import (
	"errors"

	"github.com/go-playground/validator/v10"
	validatorWrap "github.com/ranefattesingh/microservices/auth/validator"
	fieldErr "github.com/ranefattesingh/microservices/pkg/errors"
)

type ModificationRequest interface {
	LoginRequest | RegisterRequest | RefreshRequest
}

func ValidateRequest[T ModificationRequest](v *validatorWrap.Validator, req T) []fieldErr.FieldError {
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
			fields := make([]fieldErr.FieldError, 0, len(ve))

			for _, ve := range ve {
				fields = append(fields, fieldErr.FieldError{
					Field:   ve.Field(),
					Message: parse(ve),
				})
			}

			return fields
		}
	}

	return nil
}
