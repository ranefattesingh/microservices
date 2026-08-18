package handler

import (
	"net/http"

	"github.com/ranefattesingh/microservices/auth/handler/dto"
	"github.com/ranefattesingh/microservices/pkg/encoding/json"
	"go.uber.org/zap"
)

// Register handles user registration
// @Summary Register
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} json.Error
// @Router /v1/auth/register [post]
func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	req := dto.RegisterRequest{}
	if err := json.Decode(r, &req); err != nil {
		h.logger.Error("register: unable to decode request", zap.Error(err))
		handleError(w, ErrBadRequest)
		return
	}

	if fieldErrs := dto.ValidateRequest(h.validator, req); len(fieldErrs) > 0 {
		h.logger.Error("register: invalid request")
		json.Respond(w).BadRequest(ErrBadRequest, fieldErrs...)
		return
	}

	if err := h.service.Register(r.Context(), req); err != nil {
		h.logger.Error("register: service error", zap.Error(err))
		handleError(w, err)
		return
	}

	json.Respond(w).Created(map[string]string{"email": req.Email})
}
