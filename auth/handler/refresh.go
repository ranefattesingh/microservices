package handler

import (
	"net/http"

	"github.com/ranefattesingh/microservices/auth/handler/dto"
	"github.com/ranefattesingh/microservices/pkg/encoding/json"
	"go.uber.org/zap"
)

// Refresh handles token refresh
// @Summary Refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} json.Error
// @Router /v1/auth/refresh [post]
func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	req := dto.RefreshRequest{}
	if err := json.Decode(r, &req); err != nil {
		h.logger.Error("refresh: unable to decode request", zap.Error(err))
		handleError(w, ErrBadRequest)
		return
	}

	if fieldErrs := dto.ValidateRequest(h.validator, req); len(fieldErrs) > 0 {
		h.logger.Error("refresh: invalid request")
		json.Respond(w).BadRequest(ErrBadRequest, fieldErrs...)
		return
	}

	token, err := h.service.Refresh(r.Context(), req)
	if err != nil {
		h.logger.Error("refresh: service error", zap.Error(err))
		handleError(w, err)
		return
	}

	json.Respond(w).JSON(map[string]string{"token": token})
}
