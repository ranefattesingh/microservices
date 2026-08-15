package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ranefattesingh/microservices/user/json"
)

type RouteProvider interface {
	Routes(r chi.Router)
}

func NewRouter(providers ...RouteProvider) chi.Router {
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.Respond(w).JSON("pong!")
	})

	for _, provider := range providers {
		provider.Routes(r)
	}

	return r
}
