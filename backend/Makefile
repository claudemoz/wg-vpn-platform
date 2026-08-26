APP_NAME=wg-vpn-platform
MAIN=./cmd/api
BIN=bin/$(APP_NAME)

DB_HOST?=localhost
DB_PORT?=5432
DB_USER?=postgres
DB_PASSWORD?=postgres
DB_NAME?=wg-vpn-platform
DB_SSLMODE?=disable

DB_URL?=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
SEED_DIR=seed
SEED_DB_URL=$(DB_URL)&x-migrations-table=seed_schema_migrations

.PHONY: migration migrate-up migrate-down migrate-version \
	seed-migration seed-up seed-down seed-version \
	run dev build test fmt tidy seed setup-db

run:
	go run $(MAIN)

dev:
	air -c .air.toml

build:
	go build -o $(BIN) $(MAIN)
	
test:
	go test ./... -cover

fmt:
	go fmt ./...

tidy:
	go mod tidy
	
# 	
migration:
	migrate create -ext sql -dir migrations -seq $(name)

seed-migration:
	migrate create -ext sql -dir $(SEED_DIR) -seq $(name)

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

seed-up:
	migrate -path $(SEED_DIR) -database "$(SEED_DB_URL)" up

seed-down:
	migrate -path $(SEED_DIR) -database "$(SEED_DB_URL)" down 1

setup-db: migrate-up seed-up