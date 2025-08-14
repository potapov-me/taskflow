# Project variables
PROJECT_NAME = taskflow
BACKEND_DIR = backend
FRONTEND_DIR = frontend
DOCKER_COMPOSE = docker-compose.yml

.PHONY: all run run-backend run-frontend build test clean docker-up docker-down docker-build

all: init run

## Инициализация проекта
init:
	@echo "Инициализация проекта..."
	cd $(BACKEND_DIR) && go mod tidy
	cd $(FRONTEND_DIR) && npm install

## Запуск всего проекта
run: run-backend run-frontend

## Бэкенд
run-backend:
	@echo "Запуск Go-сервера..."
	cd $(BACKEND_DIR) && go run main.go

## Фронтенд
run-frontend:
	@echo "Запуск Vue.js сервера..."
	cd $(FRONTEND_DIR) && npm run dev

## Сборка
build: build-backend build-frontend

build-backend:
	@echo "Сборка Go-бинарника..."
	cd $(BACKEND_DIR) && GOOS=linux GOARCH=amd64 go build -o ../bin/$(PROJECT_NAME)-api

build-frontend:
	@echo "Сборка Vue.js приложения..."
	cd $(FRONTEND_DIR) && npm run build

## Тестирование
test: test-backend test-frontend

test-backend:
	@echo "Запуск Go-тестов..."
	cd $(BACKEND_DIR) && go test -v -cover ./...

test-frontend:
	@echo "Запуск Vue.js тестов..."
	cd $(FRONTEND_DIR) && npm run test:unit

## Очистка
clean:
	@echo "Очистка проекта..."
	rm -rf bin/
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/node_modules
	cd $(BACKEND_DIR) && go clean -cache -testcache

## Docker
docker-up:
	@echo "Запуск Docker-контейнеров..."
	docker-compose -f $(DOCKER_COMPOSE) up -d --build

docker-down:
	@echo "Остановка Docker-контейнеров..."
	docker-compose -f $(DOCKER_COMPOSE) down -v

docker-build:
	@echo "Сборка Docker-образов..."
	docker build -t $(PROJECT_NAME)-api:latest -f $(BACKEND_DIR)/Dockerfile $(BACKEND_DIR)
	docker build -t $(PROJECT_NAME)-frontend:latest -f $(FRONTEND_DIR)/Dockerfile $(FRONTEND_DIR)

## Деплой
deploy-staging: build docker-build
	@echo "Деплой на staging..."
    # Добавьте свои команды деплоя

## Помощь
help:
	@echo "Доступные команды:"
	@echo "  make init         - Инициализация зависимостей"
	@echo "  make run          - Запуск всего проекта (бэкенд + фронтенд)"
	@echo "  make build        - Сборка production-версии"
	@echo "  make test         - Запуск всех тестов"
	@echo "  make docker-up    - Запуск проекта в Docker"
	@echo "  make deploy       - Деплой на staging-окружение"
	@echo "  make clean        - Очистка артефактов сборки"
