package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ranefattesingh/microservices/user/service"
	"github.com/ranefattesingh/microservices/user/validator"
)

var _ UserHandler = (*userHandle)(nil)

type UserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
}

type userHandle struct {
	s service.UserService
	v *validator.Validator
}

func NewUserHandler(s service.UserService) *userHandle {
	return &userHandle{
		s: s,
		v: validator.NewValidator(),
	}
}

func (h *userHandle) Routes(r chi.Router) {
	r.Route("/v1/users", func(r chi.Router) {
		r.Post("/", h.CreateUser)
		r.Get("/{id}", h.GetUser)
	})
}
