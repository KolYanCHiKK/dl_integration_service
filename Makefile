.DEFAULT_GOAL := help

.PHONY: help db auth ingest up down test lint

help: ## Показать список команд
	@echo Доступные команды: db auth ingest up down test lint

db: ## Поднять Postgres
	docker compose up -d postgres

auth: ## Запустить auth через air
	air -c build/air/.air.auth.toml

grpc: ## Запустить ingest через air
	air -c build/air/.air.grpc.toml

up: ## Всё в Docker с air
	docker compose up -d --build

down: ## Остановить контейнеры
	docker compose down

test: ## Запустить тесты
	go test ./...

lint: ## Проверить код
	go vet ./...