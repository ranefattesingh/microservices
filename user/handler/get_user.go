package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ranefattesingh/microservices/user/handler/dto"
	"github.com/ranefattesingh/microservices/user/json"
	"go.uber.org/zap"
)

func (h *userHandle) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("get user: unable to parse query param user_id", zap.Error(err))
		handleError(w, ErrInvalidUserID)

		return
	}

	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		h.logger.Error("get user: fail", zap.Error(err))

		handleError(w, err)

		return
	}

	json.Respond(w).JSON(dto.FromServiceModel(user))
}
