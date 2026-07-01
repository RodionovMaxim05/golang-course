SHELL := /bin/sh

CONTAINER_RUNTIME ?= docker
GOBIN := $(shell go env GOPATH)/bin

PROTOLINT := $(GOBIN)/protolint
GOIMPORTS := $(GOBIN)/goimports
GOLANGCI_LINT := $(GOBIN)/golangci-lint
PROTOC_GEN_GO := $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(GOBIN)/protoc-gen-go-grpc

PROTOC_GEN_GO_VERSION := v1.32.0
PROTOC_GEN_GO_GRPC_VERSION := v1.6.1
GOLANGCI_LINT_VERSION := v2.11.2

.PHONY: check-container-runtime up down down-volumes \
	run-tests integration-test unit-test \
	lint fix protobuf protolint gofmtcheck golint swagger tools help

## Docker / environment

up: check-container-runtime down ## Up containers with compose
	$(CONTAINER_RUNTIME) compose up --build -d

down: check-container-runtime ## Stop and remove containers
	$(CONTAINER_RUNTIME) compose down

down-volumes: check-container-runtime ## Stop and remove containers and all volumes
	$(CONTAINER_RUNTIME) compose down -v

run-tests: ## Run tests container
	$(CONTAINER_RUNTIME) run --rm --network=host tests:latest

integration-test: check-container-runtime ## Up containers and run integration tests
	@$(MAKE) down-volumes
	@$(MAKE) up
	@echo "Waiting for cluster to start" \
		&& for i in $$(seq 15); do \
			curl -sf http://localhost:28080 >/dev/null 2>&1 && break \
				|| true; \
			sleep 1; \
		done || { echo "Error: timeout"; exit 1; }
	@$(MAKE) run-tests; status=$$?; $(MAKE) down-volumes; exit $$status

check-container-runtime: ## Check container runtime is available
ifeq (0,$(MAKELEVEL))
	@$(if $(strip $(CONTAINER_RUNTIME)),\
		$(info Using $(CONTAINER_RUNTIME) as container runtime),\
		$(error No container runtime found. Install Podman or Docker))
endif

## Go tooling

lint: protolint gofmtcheck golint ## Run all linters

fix: ## Auto-fix proto/go formatting and lint issues
	$(PROTOLINT) lint -fix .
	find . -type f -name '*.go' ! -name '*.pb.go' -print0 | xargs -0 $(GOIMPORTS) -w
	$(GOLANGCI_LINT) run --fix --timeout=2m -E gocritic -v ./...

protobuf: ## Compile protobuf files
	PATH="$(GOBIN):$$PATH" protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/subscriber/subscriber.proto
	PATH="$(GOBIN):$$PATH" protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/collector/collector.proto
	PATH="$(GOBIN):$$PATH" protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/processor/processor.proto
	PATH="$(GOBIN):$$PATH" protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/common/events.proto

protolint: ## Lint .proto files
	$(PROTOLINT) .

gofmtcheck: ## Check Go formatting/imports
	find . -type f -name '*.go' ! -name '*.pb.go' -print0 | xargs -0 $(GOIMPORTS) -w
	git diff --exit-code

golint: ## Run golangci-lint
	$(GOLANGCI_LINT) run --timeout=2m -E gocritic -v ./...

swagger: ## Generate Swagger docs for API Gateway
	swag init -g cmd/app/main.go -o api/docs -d ./api --parseDependency --parseInternal

unit-test: ## Run unit tests and generate coverage report
	go test -race -coverprofile cover.out \
		$(shell go list ./... | grep -E -v '/proto')
	go tool cover -html=cover.out -o cover.html

tools: ## Install dev tools
	go install github.com/yoheimuta/protolint/cmd/protolint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
