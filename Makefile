.PHONY: build run-input run-orchestrator test tidy env
up:
	@docker-compose up -d --build
down:
	@docker-compose down
build:
	@go build -o ./bin/input ./cmd/input/main.go
	@go build -o ./bin/orchestrator ./cmd/orchestrator/main.go

run-input:
	@go run ./cmd/input/main.go

run-orchestrator:
	@go run ./cmd/orchestrator/main.go

test:
	@./scripts/test.sh

tidy:
	@go mod tidy

env:
	@cp .env.example .env
	@cp .env.docker.example .env.docker
