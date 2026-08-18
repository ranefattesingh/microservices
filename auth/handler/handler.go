package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ranefattesingh/microservices/auth/service"
	"github.com/ranefattesingh/microservices/auth/validator"
	"github.com/ranefattesingh/microservices/pkg/encoding/json"
	"go.uber.org/zap"
)

var (
	ErrBadRequest = errors.New("server fail to handle bad request")
)

type UserHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
}

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Routes(r chi.Router)
}

type authHandler struct {
	logger    *zap.Logger
	service   service.AuthService
	validator *validator.Validator
}

func NewAuthHandler(l *zap.Logger, s service.AuthService) *authHandler {
	return &authHandler{
		logger:    l,
		service:   s,
		validator: validator.NewValidator(),
	}
}

func (h *authHandler) Routes(r chi.Router) {
	r.Route("/v1/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/refresh", h.Refresh)
	})
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrBadRequest):
		json.Respond(w).BadRequest(err)
	case errors.Is(err, service.ErrInvalidCredentials):
		json.Respond(w).BadRequest(err)
	case errors.Is(err, service.ErrEmailAlreadyTaken):
		json.Respond(w).Conflict(err)
	default:
		json.Respond(w).InternalServerError()
	}
}
