package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ranefattesingh/ecommerce-platform/user-service/models"
)

type Handle struct{}

// Handler for creating user.
func (h *Handle) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Marshal CreateUser model
	var request models.CreateUser

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("fail to unmarshal the incoming request."))

		return
	}

	w.WriteHeader(http.StatusCreated)
}
