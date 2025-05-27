BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

GOOSE_DRIVER="postgres"
GOOSE_MIGRATION_DIR="./migrations"
GOOSE_DBSTRING="PGUSER=$(PGUSER) PGPASSWORD=$(PGPASSWORD) PGHOST=$(PGHOST) PGPORT=$(PGPORT) PGDATABASE=$(PGDATABASE)"

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) --config="./configs/config.yaml"

migration:
	goose create change_my_name sql -dir $GOOSE_MIGRATION_DIR

migrate:
	goose up GOOSE_DRIVER=$GOOSE_DRIVER GOOSE_MIGRATION_DIR=$GOOSE_MIGRATION_DIR GOOSE_DBSTRING=$GOOSE_DBSTRING

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
