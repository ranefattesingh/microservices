package handler

import (
	"net/http"

	"github.com/ranefattesingh/microservices/user/handler/dto"
	"github.com/ranefattesingh/microservices/user/json"
	"go.uber.org/zap"
)

func (h *userHandle) CreateUser(w http.ResponseWriter, r *http.Request) {
	req := dto.CreateUserRequest{}
	err := json.Decode(r, &req)
	if err != nil {
		h.logger.Error("create user: unable to decode request", zap.Error(err))
		handleError(w, ErrBadRequest)

		return
	}

	if fieldErrs := dto.ValidateUser(h.validator, req); fieldErrs != nil {
		h.logger.Error("create user: invalid request", zap.Error(err))
		json.Respond(w).BadRequest(ErrBadRequest, fieldErrs...)

		return
	}

	user := req.ToServiceModel()
	id, err := h.service.CreateUser(r.Context(), user)
	if err != nil {
		h.logger.Error("create user: fail", zap.Error(err))
		handleError(w, err)

		return
	}

	json.Respond(w).Created(map[string]int64{"id": id})
}
