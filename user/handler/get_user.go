package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ranefattesingh/microservices/pkg/encoding/json"
	"github.com/ranefattesingh/microservices/user/handler/dto"
	"go.uber.org/zap"
)

// GetUser fetches a single user by ID.
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} json.Error
// @Failure 404 {object} json.Error
// @Router /v1/users/{id} [get]
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
