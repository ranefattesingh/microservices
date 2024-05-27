package core

import (
	"context"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/ranefattesingh/pkg/json"
)

const validEmailRegex = "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"

var (
	ErrNameIsRequired = &json.Error{
		HTTPStatusCode: http.StatusBadRequest,
		Message:        "name is required",
		Code:           1,
	}

	ErrEmailIsRequired = &json.Error{
		HTTPStatusCode: http.StatusBadRequest,
		Message:        "email is required",
		Code:           2,
	}

	ErrEmailIsInvalid = &json.Error{
		HTTPStatusCode: http.StatusBadRequest,
		Message:        "email is invalid",
		Code:           3,
	}

	ErrPasswordIsRequired = &json.Error{
		HTTPStatusCode: http.StatusBadRequest,
		Message:        "password is invalid",
		Code:           4,
	}
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	IsAdmin    bool      `json:"is_admin"`
	CreateDate string    `json:"created_date"`
	UpdateDate string    `json:"updated_date"`
}

func (u User) Valid(context.Context) error {
	if u.Name == "" {
		return ErrNameIsRequired
	}

	if u.Email == "" {
		return ErrEmailIsRequired
	}

	regex := regexp.MustCompile(validEmailRegex)
	if !regex.MatchString(u.Email) {
		return ErrEmailIsInvalid
	}

	if u.Password == "" {
		return ErrPasswordIsRequired
	}

	return nil
}

type Users []*User

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
