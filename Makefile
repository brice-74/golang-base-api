SHELL := /bin/bash

docker_postgres_container_name_test := "base-database-test"

docker_api_image_name := "base-api-dev"
docker_api_container_name := "base-api-dev"
docker_database_container_name := "base-database-dev"
current_dir := $(shell pwd)

dr_api := docker run -v $(current_dir):/api -w /api $(docker_api_image_name)
de_api := docker exec -it $(docker_api_container_name)
de_db := docker exec -it $(docker_database_container_name)
dr_golangci := docker run --rm -v $(current_dir):/api -w /api golangci/golangci-lint:v1.41.1 golangci-lint

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

.DEFAULT_GOAL := help
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/dev/air: run api using air reload inside Docker container in dev mode
.PHONY: run/dev/air
run/dev/air:
	@docker-compose -f docker-compose-api-air.dev.yml -f docker-compose-db.dev.yml --env-file .env.dev up --build 

## run/dev/reflex: run api using reflex reload inside Docker container in dev mode
.PHONY: run/dev/reflex
run/dev/reflex:
	@docker-compose -f docker-compose-api-reflex.dev.yml -f docker-compose-db.dev.yml --env-file .env.dev up --build 

# ==================================================================================== #
# DATABASE
# ==================================================================================== #

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	@$(dr_api) migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	@$(de_api) sh -c "migrate -path ./migrations -database \$$DATABASE_URL up"

## db/migrations/down: revert database migrations
.PHONY: db/migrations/down
db/migrations/down:
	@echo 'Running down migrations...'
	@$(de_api) sh -c "migrate -path ./migrations -database \$$DATABASE_URL down"

# ==================================================================================== #
# DOCKER
# ==================================================================================== #

## docker/sh: connect to the api container while running in parallel
.PHONY: docker/sh
docker/sh:
	@$(de_api) /bin/sh

## docker/psql: connect to the database container while running in parallel
.PHONY: docker/psql
docker/psql:
	@$(de_db) sh -c "exec psql -U \$$POSTGRES_USER -d \$$POSTGRES_DB"

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## qa/test optional(func=$1 pkg=$2): run automated tests
.PHONY: qa/test
qa/test:
	@echo 'Running tests...'
	-@docker-compose -f docker-compose.test.yml run --rm apitest sh -c "sleep 10 && migrate -path ./migrations -database \$$DATABASE_URL up && go test -p 1 -v -vet=off -run \"$(func)\" ./.../$(pkg)"
	@echo 'Stop & Remove db services...'
	@docker stop $(docker_postgres_container_name_test) && docker rm $(docker_postgres_container_name_test)

## qa/coverage: run automated tests and create coverage report
.PHONY: qa/coverage
qa/coverage:
	@echo 'Running tests and creating coverage report...'
	-@docker-compose -f docker-compose.test.yml run --rm apitest sh -c "sleep 10 && migrate -path ./migrations -database \$$DATABASE_URL up && go test -p 1 -coverprofile=coverage.txt -covermode=atomic -v -vet=off ./..."
	@echo 'Stop & Remove db services...'
	@docker stop $(docker_postgres_container_name_test) && docker rm $(docker_postgres_container_name_test)

## watch/coverage: watch coverage report on navigator
.PHONY: qa/coverage
watch/coverage:
	@go tool cover -html=coverage.txt

## qa/lint: run golangci linters
.PHONY: qa/lint
qa/lint:
	@echo 'Running linters...'
	@$(dr_golangci) run -v