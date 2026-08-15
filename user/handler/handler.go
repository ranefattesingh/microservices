package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ranefattesingh/microservices/user/json"
	"github.com/ranefattesingh/microservices/user/service"
	"github.com/ranefattesingh/microservices/user/validator"
	"go.uber.org/zap"
)

var (
	ErrBadRequest    = errors.New("server fail to handle bad request")
	ErrInvalidUserID = errors.New("invalid user id")
)

var _ UserHandler = (*userHandle)(nil)

type UserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
}

type userHandle struct {
	logger    *zap.Logger
	service   service.UserService
	validator *validator.Validator
}

func NewUserHandler(l *zap.Logger, s service.UserService) *userHandle {
	return &userHandle{
		logger:    l,
		service:   s,
		validator: validator.NewValidator(),
	}
}

func (h *userHandle) Routes(r chi.Router) {
	r.Route("/v1/users", func(r chi.Router) {
		r.Post("/", h.CreateUser)
		r.Get("/{id}", h.GetUser)
		r.Put("/{id}", h.UpdateUser)
	})
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidUserID):
		json.Respond(w).BadRequest(err)
	case errors.Is(err, service.ErrEmailAlreadyTaken):
		json.Respond(w).Conflict(err)
	case errors.Is(err, service.ErrUserDoesNotExist):
		json.Respond(w).NotFound(err)
	default:
		json.Respond(w).InternalServerError()
	}
}
