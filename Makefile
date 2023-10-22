export OSTYPE := $(shell uname -s)

include ./utils/depends.Makefile

# Exporting bin folder to the path for makefile
export PATH   := $(PWD)/bin:$(PATH)
export SHELL  := bash

up: dev
down: docker-stop
clean: docker-teardown docker-clean

# ==============================================================================
# Modules
tidy:
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

deps-upgrade:
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

# ==============================================================================
# Build commands
build:
	@ printf "Building application... "
	@ go build \
		-trimpath  \
		-o clean \
		./cmd/clean
	@ echo "done"

build-race:
	@ printf "Building application with race flag... "
	@ go build \
		-trimpath  \
		-race      \
		-o clean \
		./cmd/clean
	@ echo "done"

# ==============================================================================
# Test commands

TESTS_ARGS := --format testname --jsonfile tmp/tests/gotestsum.json.out
TESTS_ARGS += --max-fails 2
TESTS_ARGS += -- ./...
TESTS_ARGS += -test.parallel 2
TESTS_ARGS += -test.count    1
TESTS_ARGS += -test.failfast
TESTS_ARGS += -test.coverprofile   tmp/tests/coverage.out
TESTS_ARGS += -test.timeout        5s
TESTS_ARGS += -race

run-tests: $(GOTESTSUM)
	@ gotestsum $(TESTS_ARGS) -short

tests: run-tests $(TPARSE)
	@cat tmp/tests/gotestsum.json.out | $(TPARSE) -all -notests

coverage: run-tests
	@go tool cover -html=./tmp/tests/coverage.out

# ==============================================================================
# Docker commands
.ONESHELL:
image-build:
	@ echo "Docker Build"
	@ DOCKER_BUILDKIT=0 docker build \
 		--file docker/Dockerfile \
        --tag clean \
        	.

# ==============================================================================
# Docker compose commands
dev:
	docker-compose -f docker/docker-compose.dev.yml up -d --build

dev-env:
	@ docker-compose -f docker/docker-compose.dev.yml up -d --build postgresql redis otelcol jaeger prometheus node_exporter

docker-stop:
	@ docker-compose -f docker/docker-compose.dev.yml down

docker-teardown:
	@ docker-compose -f docker/docker-compose.dev.yml down --remove-orphans -v

docker-clean:
	@ docker image prune -f
	@ docker rmi clean-clean
	@ docker volume prune -f

# ==============================================================================
# Database migrations
PG_MIGRATIONS_DSN := "postgres://postgres:postgres@localhost:5432/clean_db?sslmode=disable"
PG_SEED_DSN := "postgres://postgres:postgres@localhost:5432/clean_db?sslmode=disable&x-migrations-table=seed"
MIGRATIONS_PATH := "db/migrations"
SEED_PATH := "db/seed"

migrate_version: $(MIGRATE)
	@ migrate -database $(PG_MIGRATIONS_DSN) -path $(MIGRATIONS_PATH) version

migrate_up: $(MIGRATE)
	@ read -p "How many migration you wants to perform (default value: [all]): " N; \
	migrate -database $(PG_MIGRATIONS_DSN) -path $(MIGRATIONS_PATH) up ${NN}

migrate_down: $(MIGRATE)
	@ read -p "How many migration you wants to perform (default value: [all]): " N; \
	migrate -database $(PG_MIGRATIONS_DSN) -path $(MIGRATIONS_PATH) down ${NN}

migrate_force: $(MIGRATE)
	@ migrate -database $(PG_MIGRATIONS_DSN) -path $(MIGRATIONS_PATH) force 1

seed_up: $(MIGRATE)
	@ migrate -database $(PG_SEED_DSN) -path $(SEED_PATH) up

seed_down: $(MIGRATE)
	@ migrate -database $(PG_SEED_DSN) -path $(SEED_PATH) down

# ==============================================================================
# Tools commands

lint: $(GOLANGCI)
	echo "Starting linters"
	golangci-lint version
	golangci-lint run ./...

mock: $(MOCKERY)
	mockery

clean:
	@ rm -rf tmp