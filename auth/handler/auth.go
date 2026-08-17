package handler

import (
	"net/http"

	"github.com/ranefattesingh/microservices/auth/json"
	"github.com/ranefattesingh/microservices/auth/handler/dto"
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
		json.Respond(w).BadRequest(ErrBadRequest)
		return
	}

	if fieldErrs := req.Validate(); len(fieldErrs) > 0 {
		json.Respond(w).BadRequest(ErrBadRequest, fieldErrs...)
		return
	}

	if err := h.service.Register(r.Context(), req); err != nil {
		h.logger.Error("register: service error", zap.Error(err))
		json.Respond(w).InternalServerError()
		return
	}

	json.Respond(w).Created(map[string]string{"email": req.Email})
}

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
		json.Respond(w).BadRequest(ErrBadRequest)
		return
	}

	if fieldErrs := req.Validate(); len(fieldErrs) > 0 {
		json.Respond(w).BadRequest(ErrBadRequest, fieldErrs...)
		return
	}

	token, refresh, err := h.service.Login(r.Context(), req)
	if err != nil {
		h.logger.Error("login: service error", zap.Error(err))
		json.Respond(w).InternalServerError()
		return
	}

	json.Respond(w).JSON(map[string]string{"token": token, "refresh_token": refresh})
}

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
		json.Respond(w).BadRequest(ErrBadRequest)
		return
	}

	if req.RefreshToken == "" {
		json.Respond(w).BadRequest(ErrBadRequest)
		return
	}

	token, err := h.service.Refresh(r.Context(), req)
	if err != nil {
		h.logger.Error("refresh: service error", zap.Error(err))
		json.Respond(w).InternalServerError()
		return
	}

	json.Respond(w).JSON(map[string]string{"token": token})
}

