package handler

import (
	"errors"
	"net/http"

	"github.com/ranefattesingh/microservices/user/json"
)

var (
	ErrUnableToDecode = errors.New("unable to decode the request")
)

func (h *userHandle) AddUser(w http.ResponseWriter, r *http.Request) {
	req := CreateUserRequest{}
	err := json.Decode(r, &req)
	if err != nil {
		json.Respond(w).BadRequest(err)

		return
	}

	json.Respond(w).ResponseJSON(req)
}
