package handler

import (
	"net/http"

	"github.com/ranefattesingh/microservices/auth/handler/dto"
	"github.com/ranefattesingh/microservices/pkg/encoding/json"
	"go.uber.org/zap"
)

// Login handles user login
// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} json.Error
// @Router /v1/auth/login [post]
func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	req := dto.LoginRequest{}
	if err := json.Decode(r, &req); err != nil {
		h.logger.Error("login: unable to decode request", zap.Error(err))
		handleError(w, ErrBadRequest)
		return
	}

	if fieldErrs := dto.ValidateRequest(h.validator, req); len(fieldErrs) > 0 {
		h.logger.Error("login: invalid request")
		json.Respond(w).BadRequest(ErrBadRequest, fieldErrs...)
		return
	}

	token, refresh, err := h.service.Login(r.Context(), req)
	if err != nil {
		h.logger.Error("login: service error", zap.Error(err))
		handleError(w, err)
		return
	}

	json.Respond(w).JSON(map[string]string{"token": token, "refresh_token": refresh})
}
