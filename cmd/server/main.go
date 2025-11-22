package main

import (
	"AvitoTech/internal/domain/teams"
	"AvitoTech/internal/domain/user"
	"AvitoTech/internal/http"
	"AvitoTech/internal/http/handlers"
	"AvitoTech/internal/infrastructure/postgres"
	"AvitoTech/pkg/logger"
	"log"
	"os"

	"go.uber.org/zap"
)

func main() {
	isDev := os.Getenv("ENVIRONMENT") != "production"
	if err := logger.Init(isDev); err != nil {
		log.Fatal("Не удалось инициализировать логгер: ", err)
	}
	defer logger.Sync()

	logger.Log.Info("Запуск сервиса",
		zap.Bool("development", isDev),
		zap.String("version", "1.0.0"),
	)

	db, err := postgres.NewDB()
	if err != nil {
		logger.Log.Fatal("Ошибка подключения к базе данных", zap.Error(err))
	}
	defer db.CloseDB()

	teamRepo := postgres.NewTeamRepo(db)
	userRepo := postgres.NewUserRepo(db)

	teamService := teams.NewService(teamRepo, userRepo)
	teamHandler := handlers.NewTeamHandler(teamService)

	userService := user.NewService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	router := http.NewRouter(teamHandler, userHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Log.Info("Запуск HTTP сервера", zap.String("port", port))

	if err := http.StartServer(router, ":"+port); err != nil {
		logger.Log.Fatal("Ошибка запуска сервера", zap.Error(err))
	}
}
