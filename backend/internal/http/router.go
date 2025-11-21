package http

import (
	"AvitoTech/internal/http/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(teamHandler *handlers.TeamHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/team/add", teamHandler.CreateTeam)
	r.Get("/team/get", teamHandler.GetTeam)
	return r
}
