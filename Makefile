include .env

.SILENT:
.DEFAULT_GOAL = run
.PHONY: build run stop migrate-create migrate-down migrate-up clean

CMD = docker-compose up --remove-orphans

MIGRATION_DIR = ./scripts/migrations/
POSTGRES_URL = postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5436/$(DB_NAME)?sslmode=$(DB_SSLMODE)

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: build
	$(CMD)

rebuild: build
	$(CMD) --build

stop:
	docker-compose down

migrate-create:
	migrate create -ext sql -dir $(MIGRATION_DIR) 'balance_management'

migrate-up:
	migrate -path $(MIGRATION_DIR) -database $(POSTGRES_URL) up

migrate-down:
	migrate -path $(MIGRATION_DIR) -database $(POSTGRES_URL) down

lint:
	golangci-lint run

clean:
	rm -rf .bin/