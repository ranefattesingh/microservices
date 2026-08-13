package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ranefattesingh/microservices/user/service"
)

var _ UserHandler = (*userHandle)(nil)

type UserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
}

type userHandle struct {
	s service.UserService
}

func NewUserHandler(s service.UserService) *userHandle {
	return &userHandle{
		s: s,
	}
}

func (h *userHandle) Routes(r chi.Router) {
	r.Post("/", h.CreateUser)
}
