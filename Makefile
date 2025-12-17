.PHONY: help build run test clean docker-up docker-down

APP_NAME=smartstocks-api
DOCKER_COMPOSE=docker-compose

help:
	@echo "Smart Stocks API - Comandos disponibles:"
	@echo ""
	@echo "  install         - Instalar dependencias"
	@echo "  build           - Compilar aplicación"
	@echo "  run             - Ejecutar servidor"
	@echo "  test            - Ejecutar tests"
	@echo "  docker-up       - Levantar con Docker"
	@echo "  docker-down     - Detener Docker"
	@echo "  docker-logs     - Ver logs"
	@echo "  clean           - Limpiar archivos"

install:
	go mod download
	go mod tidy

build:
	go build -o bin/$(APP_NAME) ./cmd/api

run:
	go run ./cmd/api/main.go

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

docker-build:
	docker build -t $(APP_NAME):latest .

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f

docker-restart: docker-down docker-up

fmt:
	go fmt ./...

vet:
	go vet ./...

.DEFAULT_GOAL := help