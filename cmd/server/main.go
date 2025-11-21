package main

import (
	"AvitoTech/internal/domain/teams"
	"AvitoTech/internal/http"
	"AvitoTech/internal/http/handlers"
	"AvitoTech/internal/infrastructure/postgres"
	"log"
	"os"
)

func main() {
	db, err := postgres.NewDB()
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных: ", err)
	}
	defer db.CloseDB()

	teamRepo := postgres.NewTeamRepo(db)

	teamService := teams.NewService(teamRepo, nil)
	teamHandler := handlers.NewTeamHandler(teamService)

	router := http.NewRouter(teamHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.StartServer(router, ":"+port); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
