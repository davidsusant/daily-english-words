.PHONY: run build migrate seed clean

ifneq (,$(wildcard ./.env))
include .env
export
endif

DATABASE_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

run:
	go run ./cmd/server/main.go

build:
	go build -o bin/server ./cmd/server/main.go

migrate:
	psql "$(DATABASE_URL)" -f migrations/001_create_words.sql

seed:
	psql "$(DATABASE_URL)" -f seed.sql

clean:
	rm -rf bin/