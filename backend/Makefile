DATABASE_URL ?= postgres://todo:todo@localhost:5433/todos?sslmode=disable

.PHONY: run build expand test vet tidy migrate-up migrate-down sqlc-generate

## run: ejecuta la API en desarrollo
run:
	go run ./cmd/todo-api

## build: compila el binario en bin/
build:
	go build -o bin/todo-api.exe ./cmd/todo-api

## test: ejecuta los tests
test:
	go test ./...

## vet: linters nativos de Go
vet:
	go vet ./...

## migrate-up: aplica las migraciones pendientes
migrate-up:
	go tool goose -dir db/migrations postgres "$(DATABASE_URL)" up

## migrate-down: revierte la última migración
migrate-down:
	go tool goose -dir db/migrations postgres "$(DATABASE_URL)" down

## sqlc-generate: regenera el código de acceso a datos desde db/queries
sqlc-generate:
	go tool sqlc generate
