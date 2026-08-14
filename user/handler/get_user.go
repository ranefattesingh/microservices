package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ranefattesingh/microservices/user/json"
)

var ErrUserDoesNotExist = errors.New("user does not exist")

func (h *userHandle) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.Respond(w).BadRequest(err)
		return
	}

	user, err := h.s.GetUser(r.Context(), id)
	if err != nil {
		json.Respond(w).InternalServerError()
		return
	}

	if user == nil {
		json.Respond(w).NotFound(ErrUserDoesNotExist)
		return
	}

	json.Respond(w).ResponseJSON(user)
}
