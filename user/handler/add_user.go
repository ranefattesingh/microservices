package handler

import (
	"net/http"

	"github.com/ranefattesingh/microservices/user/handler/dto"
	"github.com/ranefattesingh/microservices/user/json"
	"github.com/ranefattesingh/microservices/user/pkg"
	"github.com/ranefattesingh/microservices/user/service/models"
)

var (
	ErrUnableToDecode = "unable to decode the request"
	ErrInvalidRequest = "invalid request"
)

func (h *userHandle) CreateUser(w http.ResponseWriter, r *http.Request) {
	req := dto.CreateUserRequest{}
	err := json.Decode(r, &req)
	if err != nil {
		json.Respond(w).BadRequest(pkg.NewHTTPError(ErrUnableToDecode))

		return
	}

	if fieldErrs := h.v.ValidateCreateUser(req); fieldErrs != nil {
		json.Respond(w).BadRequest(pkg.NewHTTPError(ErrInvalidRequest, fieldErrs...))

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
