package http

import (
	"AvitoTech/internal/http/handlers"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(teamHandler *handlers.TeamHandler, userHandler *handlers.UserHandler, prHandler *handlers.PRHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/team", func(r chi.Router) {
		r.Post("/add", teamHandler.CreateTeam)
		r.Get("/get", teamHandler.GetTeam)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", userHandler.SetUserActive)
		r.Get("/getReview", userHandler.GetUserReviews)
	})

	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", prHandler.CreatePR)
		r.Post("/merge", prHandler.MergePR)
		r.Post("/reassign", prHandler.ReassignPR)
	})

	return r
}

func StartServer(router *chi.Mux, addr string) error {
	log.Printf("Сервер запущен на %s", addr)
	return http.ListenAndServe(addr, router)
}
