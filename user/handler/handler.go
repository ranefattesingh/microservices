package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

var _ UserHandler = (*userHandle)(nil)

type UserHandler interface {
	AddUser(w http.ResponseWriter, r *http.Request)
}

type userHandle struct {
}

func NewUserHandler() *userHandle {
	return &userHandle{}
}

func (h *userHandle) Routes(r chi.Router) {
	r.Post("/", h.AddUser)
}
