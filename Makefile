.PHONY: help setup dev build test clean

help:
	@echo "Доступные команды:"
	@echo "  make setup     - Установить все зависимости"
	@echo "  make dev       - Запустить в режиме разработки"
	@echo "  make build     - Собрать проект"
	@echo "  make test      - Запустить тесты"
	@echo "  make clean     - Очистить временные файлы"

setup:
	cd backend && go mod tidy
	cd frontend && npm install
	pre-commit install

dev:
	docker-compose up --build

build:
	docker-compose build

test:
	cd backend && go test ./...
	cd frontend && npm test

clean:
	docker-compose down -v
	rm -rf frontend/node_modules frontend/.nuxt frontend/.output
	rm -rf backend/bin
	go clean -modcache