package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ranefattesingh/ecommerce-platform/user-service/models"
	"github.com/ranefattesingh/ecommerce-platform/user-service/service"
)

type UserHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userService service.UserService
}

func NewHandler(userService service.UserService) *userHandler {
	return &userHandler{
		userService: userService,
	}
}

// Handler for creating user.
func (h *userHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Marshal CreateUser model
	var request models.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("fail to unmarshal the incoming request."))

		return
	}

	err = h.userService.CreateUser(r.Context(), request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server fail to create user."))

		return
	}

	w.WriteHeader(http.StatusCreated)
}
