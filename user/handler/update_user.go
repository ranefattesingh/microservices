package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ranefattesingh/microservices/user/handler/dto"
	"github.com/ranefattesingh/microservices/user/json"
	"go.uber.org/zap"
)

// UpdateUser updates an existing user.
// @Summary Update a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body dto.UpdateUserRequest true "Update user payload"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} json.Error
// @Failure 404 {object} json.Error
// @Failure 409 {object} json.Error
// @Router /v1/users/{id} [put]
func (h *userHandle) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("update user: unable to parse query param user_id", zap.Error(err))
		handleError(w, ErrInvalidUserID)

		return
	}

	req := dto.UpdateUserRequest{}
	if err := json.Decode(r, &req); err != nil {
		h.logger.Error("update user: unable to decode request", zap.Error(err))
		handleError(w, ErrBadRequest)

		return
	}

	if fieldErrs := dto.ValidateUser(h.validator, req); fieldErrs != nil {
		h.logger.Error("update user: invalid request", zap.Error(err))
		json.Respond(w).BadRequest(ErrBadRequest, fieldErrs...)

		return
	}

	user := req.ToServiceModel()
	if err = h.service.UpdateUser(r.Context(), id, user); err != nil {
		h.logger.Error("update user: fail", zap.Error(err))

		handleError(w, err)

		return
	}

	json.Respond(w).NoContent()
}
