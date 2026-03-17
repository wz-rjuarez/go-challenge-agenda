.PHONY: proto mock swagger build run-agenda run-api docker-up docker-down docker-up-postgres docker-down-postgres test

proto:
	buf generate

mock:
	mockery

swagger:
	cd services/api && $(shell go env GOPATH)/bin/swag init -g cmd/main.go -o docs

build:
	go build ./...

run-agenda:
	go run ./services/agenda/cmd

run-api:
	go run ./services/api/cmd

docker-up:
	docker compose up --build

docker-up-postgres:
	DB_DRIVER=postgres DB_SOURCE="host=postgres user=agenda password=agenda_pass dbname=agenda port=5432 sslmode=disable" docker compose --profile postgres up --build

docker-down:
	docker compose down -v

docker-down-postgres:
	docker compose --profile postgres down -v

test:
	go test ./...
