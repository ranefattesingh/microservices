package handler

import (
	"errors"
	"net/http"

	"github.com/ranefattesingh/microservices/user/json"
	"github.com/ranefattesingh/microservices/user/service/models"
)

var (
	ErrUnableToDecode = errors.New("unable to decode the request")
)

func (h *userHandle) CreateUser(w http.ResponseWriter, r *http.Request) {
	req := CreateUserRequest{}
	err := json.Decode(r, &req)
	if err != nil {
		json.Respond(w).BadRequest(err)

		return
	}

	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}

	id, err := h.s.CreateUser(r.Context(), user)
	if err != nil {
		json.Respond(w).InternalServerError()

		return
	}

	json.Respond(w).Created(map[string]int64{"id": id})
}
