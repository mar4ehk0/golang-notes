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
	# setup migrate 
	curl -o /tmp/migrate.linux-amd64.tar.gz -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz
	mkdir /tmp/golang-migrate-4.18.1 && tar fxvz /tmp/migrate.linux-amd64.tar.gz -C /tmp/golang-migrate-4.18.1
	mv /tmp/golang-migrate-4.18.1/migrate ${PROJECT_DIR}/bin && rm -rf /tmp/golang-migrate-4.18.1 && rm -rf /tmp/migrate.linux-amd64.tar.gz 
	# setup golangci-lint 
	GOBIN=$(PROJECT_BIN) GOTOOLCHAIN=go1.23.4 go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4

migration_create: ## Migration Create name=migration_name
	$(PROJECT_BIN)/migrate create -ext sql -dir ${MIGRATION_DIR} ${name}

migration_up: ## Migration Up
	$(PROJECT_BIN)/migrate -path ${MIGRATION_DIR} -database ${MIGRATION_DSN} up	
	
migration_down: ## Migration Down
	$(PROJECT_BIN)/migrate -path ${MIGRATION_DIR} -database ${MIGRATION_DSN} down

migration_force: ## Migration Force version=20250120130543
	$(PROJECT_BIN)/migrate -path ${MIGRATION_DIR} -database ${MIGRATION_DSN} force ${version}

migration_version: ## Migration Version
	$(PROJECT_BIN)/migrate -path ${MIGRATION_DIR} -database ${MIGRATION_DSN} version

build_app: ## Build app
	docker-compose -f ./docker/docker-compose.yml --env-file ./docker/.env build

run_app: ## Run app
	docker-compose -f ./docker/docker-compose.yml --env-file ./docker/.env up -d

down_app: ## Down app
	docker-compose -f ./docker/docker-compose.yml --env-file ./docker/.env down -v

lint: ## Run linter
	clear && echo "Start lint"
	$(PROJECT_BIN)/golangci-lint run -c .golangci.yml
	echo "Finish lint"

# Global
.DEFAULT_GOAL := help
