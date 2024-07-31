.PHONY: migrate up down

export DATABASE_URL=$(shell grep DATABASE_URL .env | cut -d '=' -f2)

# Path to migrations folder
MIGRATIONS_PATH := migrations

create_migration:
	goose create -dir $(MIGRATIONS_PATH)/$(shell date +%Y%m%d%H%M%S)_add_new_table sql

# Target to run all migrations
up:
	goose -dir $(MIGRATIONS_PATH) postgres $(DATABASE_URL) up

# Target to Drop all migrations
down:
	goose -dir $(MIGRATIONS_PATH) postgres $(DATABASE_URL) down


# build go 
build:
	go build -o bin/scan-attendance cmd/main.go