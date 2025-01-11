MAKEFLAGS += --silent

.PHONY: *
SHELL=/bin/bash -o pipefail

COLOR="\033[32m%-25s\033[0m %s\n"

ifneq ($(wildcard ./docker/.env),)
	include ./docker/.env
endif

PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin
PROJECT_TMP = $(PROJECT_DIR)/tmp
MIGRATION_DIR = $(PROJECT_DIR)/schema
MIGRATION_DSN = "postgres://${RDMS_DB_USER}:${RDMS_DB_PASSWORD}@127.0.0.1:${RDMS_DB_PORT}/${RDMS_DB_NAME}?sslmode=disable"

.PHONY: help
help: ## Show this help
	@printf "\033[33m%s:\033[0m\n" 'Available commands'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "  "${COLOR}, $$1, $$2}' ${MAKEFILE_LIST}

install_deps: ## Setup install deps
	curl -o /tmp/migrate.linux-amd64.tar.gz -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz
	mkdir /tmp/golang-migrate-4.18.1 && tar fxvz /tmp/migrate.linux-amd64.tar.gz -C /tmp/golang-migrate-4.18.1
	mv /tmp/golang-migrate-4.18.1/migrate ${PROJECT_DIR}/bin && rm -rf /tmp/golang-migrate-4.18.1 && rm -rf /tmp/migrate.linux-amd64.tar.gz 

migration_create: ## Migration Create name=migration_name
	$(PROJECT_BIN)/migrate create -ext sql -dir ${MIGRATION_DIR} ${name}

migration_up: ## Migration Up
	$(PROJECT_BIN)/migrate -path ${MIGRATION_DIR} -database ${MIGRATION_DSN} up	
	
migration_down: ## Migration Down
	$(PROJECT_BIN)/migrate -path ${MIGRATION_DIR} -database ${MIGRATION_DSN} down

build_app: ## Build app
	docker-compose -f ./docker/docker-compose.yml --env-file ./docker/.env build

run_app: ## Run app
	docker-compose -f ./docker/docker-compose.yml --env-file ./docker/.env up -d

down_app: ## Down app
	docker-compose -f ./docker/docker-compose.yml --env-file ./docker/.env down -v

lint: ## Run linter
	golangci-lint run -c ../.golangci.yml

# Global
.DEFAULT_GOAL := help
