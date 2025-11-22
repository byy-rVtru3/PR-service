package main

import (
	"AvitoTech/internal/domain/teams"
	"AvitoTech/internal/domain/user"
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
	userRepo := postgres.NewUserRepo(db.GetConn())

	teamService := teams.NewService(teamRepo, userRepo)
	teamHandler := handlers.NewTeamHandler(teamService)

	userService := user.NewService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	router := http.NewRouter(teamHandler, userHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.StartServer(router, ":"+port); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
