.PONY: all build deps image lint proto test

all: test build ## Run the tests and build the binary.

build: ## Build the binary.
	@go build .

deps: ## Install dependencies.
	@go get -u github.com/golang/lint/golint
	@go get -u github.com/Masterminds/glide && glide install

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

image: ## Build the Docker image.
	@docker build .

lint: ## Run golint to ensure the code follows Go styleguide.
	@golint -set_exit_status account broker certificates configuration cryptopolis domain rpc/api secrets storage

proto: ## Generate the Go definitions for the protocol buffer schemas.
	@protoc -I rpc rpc/rpc.proto --go_out=plugins=grpc:rpc

test: lint proto ## Run tests.
	@go test -v `go list ./... | grep -v /vendor/`
