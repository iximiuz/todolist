.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: test-e2e
test-e2e:
	go test -v -tags e2e -count 1 ./tests/e2e

.PHONY: build
build:
	CGO_ENABLED=0 go build -o server ./cmd/server

.PHONY: docker-build
docker-build:
	docker build -t ghcr.io/iximiuz/labs/todolist .

.PHONY: up
up:
	docker compose up --build
