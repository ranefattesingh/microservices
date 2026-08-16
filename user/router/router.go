package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/ranefattesingh/microservices/user/docs"
	"github.com/ranefattesingh/microservices/user/json"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type RouteProvider interface {
	Routes(r chi.Router)
}

func NewRouter(providers ...RouteProvider) chi.Router {
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.Respond(w).JSON("pong!")
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	for _, provider := range providers {
		provider.Routes(r)
	}

	return r
}
