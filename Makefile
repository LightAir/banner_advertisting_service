BIN := "./bin/banner"
BIN_MIGRATE := "./bin/migrate"
DOCKER_IMG="bas:develop"
DOCKER_IMG="bas-migrate:develop"

test:
	go test -race -count=100 ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

up:
	docker-compose -f deployments/docker-compose.yaml up -d

down:
	docker-compose -f deployments/docker-compose.yaml down

docs:
	swag init --parseDependency --parseInternal=true --dir ./cmd/banner

build:
	go build -v -o $(BIN) ./cmd/banner
	go build -v -o $(BIN_MIGRATE) ./cmd/goose

run-local: build
	$(BIN_MIGRATE) --config=configs/config-migrate.yaml
	$(BIN) --config=configs/config.yaml

run:
	docker-compose -f deployments/docker-compose.yaml up

build-img:
	docker build \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

migrate:
	go run cmd/goose/*.go --config=configs/config-migrate.yaml

build-migrate-img:
	docker build \
		-t $(DOCKER_MIGRATE_IMG) \
		-f build/Migrate.Dockerfile .

run-migrate-img: build-img
	docker run $(DOCKER_MIGRATE_IMG)

integration-tests:
	docker-compose -f deployments/docker-compose.test.yaml up --abort-on-container-exit --exit-code-from bas-integration-tests bas-integration-tests
	docker-compose -f deployments/docker-compose.test.yaml down

.PHONY: test install-lint-deps lint up down docs run migrate build build-img run-img integration-tests