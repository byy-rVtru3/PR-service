.PHONY: help build up down restart logs clean lint

help: ## Показать помощь
	@echo "Доступные команды:"
	@echo "  make build    - Собрать Docker образы"
	@echo "  make up       - Запустить сервис"
	@echo "  make down     - Остановить сервис"
	@echo "  make restart  - Перезапустить сервис"
	@echo "  make logs     - Показать логи"
	@echo "  make clean    - Остановить и удалить volumes"
	@echo "  make lint     - Запустить линтер (golangci-lint)"

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

lint: ## Запустить линтер
	@echo "Running go vet..."
	go vet ./...
	@echo "Running gofmt..."
	gofmt -w .
	@echo "Running staticcheck..."
	@which staticcheck > /dev/null || go install honnef.co/go/tools/cmd/staticcheck@latest
	$(shell go env GOPATH)/bin/staticcheck ./...
	@echo "✓ All checks passed!"

