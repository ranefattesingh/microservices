package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ranefattesingh/microservices/user/json"
)

type Router interface {
	Routes(chi chi.Router)
}

func NewRouter(userRouter Router) chi.Router {
	router := chi.NewRouter()

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.Respond(w).ResponseJSON("pong!")
	})

	router.Route("/users", func(r chi.Router) { userRouter.Routes(r) })

	return router
}
