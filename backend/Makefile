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