.PHONY: run build test lint vet swagger fmt tidy docker-up docker-down

MODULE := github.com/yourorg/go-fiber-template
BINARY := api

run:
	go run ./cmd/api

build:
	go build -o bin/$(BINARY) ./cmd/api

test:
	go test ./... -race -cover

vet:
	go vet ./...

fmt:
	gofmt -w .

tidy:
	go mod tidy

# Regenerate the OpenAPI docs from annotations.
swagger:
	swag init -g cmd/api/main.go -o docs

lint:
	go vet ./...

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down