.PHONY: help build up down restart logs clean

help: ## Показать помощь
	@echo "Доступные команды:"
	@echo "  make build    - Собрать Docker образы"
	@echo "  make up       - Запустить сервис"
	@echo "  make down     - Остановить сервис"
	@echo "  make restart  - Перезапустить сервис"
	@echo "  make logs     - Показать логи"
	@echo "  make clean    - Остановить и удалить volumes"

build: ## Собрать Docker образы
	docker-compose build

up: ## Запустить сервис
	docker-compose up -d

down: ## Остановить сервис
	docker-compose down

restart: ## Перезапустить сервис
	docker-compose down
	docker-compose up -d

logs: ## Показать логи
	docker-compose logs -f

clean: ## Остановить и удалить volumes
	docker-compose down -v
