package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ranefattesingh/microservices/auth/json"
	"github.com/ranefattesingh/microservices/auth/service"
	"go.uber.org/zap"
)

var (
	ErrBadRequest = errors.New("server fail to handle bad request")
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Routes(r chi.Router)
}

type authHandler struct {
	logger *zap.Logger
	service *service.AuthService
}

func NewAuthHandler(l *zap.Logger) *authHandler {
	return &authHandler{
		logger: l.Named("auth-handler"),
		service: service.NewAuthService(),
	}
}

func (h *authHandler) Routes(r chi.Router) {
	r.Route("/v1/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/refresh", h.Refresh)
	})
}
