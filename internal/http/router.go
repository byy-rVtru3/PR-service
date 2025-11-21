package http

import (
	"AvitoTech/internal/http/handlers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func NewRouter(teamHandler *handlers.TeamHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/team", func(r chi.Router) {
		r.Post("/add", teamHandler.CreateTeam)
		r.Get("/get", teamHandler.GetTeam)
	})

	return r
}

func StartServer(router *chi.Mux, addr string) error {
	log.Printf("Сервер запущен на %s", addr)
	return http.ListenAndServe(addr, router)
}
