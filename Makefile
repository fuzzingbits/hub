-include .env
export

GO_PATH := $(shell go env GOPATH 2> /dev/null)
MODULE := $(shell awk '/^module/ {print $$2}' go.mod)
NAMESPACE := $(shell awk -F "/" '/^module/ {print $$(NF-1)}' go.mod)
PROJECT_NAME := $(shell awk -F "/" '/^module/ {print $$(NF)}' go.mod)

help:
	@echo "Makefile targets:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' Makefile \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/\t\1 :: \3/p' \
	| column -t -c 1  -s '::'

full: clean full-ui full-go ## Do a full build

full-go: lint-go test-go build-go

full-ui: lint-ui build-ui

docker: clean ## Build the Docker Image
	docker build -t $(NAMESPACE)/$(PROJECT_NAME):latest .

publish: docker ## Build and publish the Docker Image
	docker push $(NAMESPACE)/$(PROJECT_NAME):latest

clean: ## Remove all git ignored file
	git clean -Xdf --exclude="!/.env"

dev-go: ## Start a dev instance of the Go Server
	clear
	@DEV=true go run main.go

dev-ui: ## Start a dev instance of the UI
	clear
	[ -f ./node_modules ] || npm install
	npm run dev

dev-docker-up: ## Start the docker containers used for development
	docker-compose -f ops/hub_dev/docker-compose.yml up -d

dev-docker-down: ## Remove the docker containers used for development
	docker-compose -f ops/hub_dev/docker-compose.yml down

build-go:
	@go install github.com/gobuffalo/packr/packr
	$(GO_PATH)/bin/packr build -o $(CURDIR)/var/$(PROJECT_NAME)
	@ln -sf $(CURDIR)/var/$(PROJECT_NAME) $(GO_PATH)/bin/$(PROJECT_NAME)

build-ui:
	npm run build

lint-go:
	@cd ; go get golang.org/x/lint/golint
	@cd ; go get golang.org/x/tools/cmd/goimports
	go get -d ./...
	gofmt -s -w .
	go vet ./...
	$(GO_PATH)/bin/golint -set_exit_status=1 ./...
	$(GO_PATH)/bin/goimports -w .

lint-ui:
	[ -f ./node_modules ] || npm install
	npm run fmt

test-go:
	@mkdir -p var/
	@go test -race -cover -coverprofile  var/coverage.txt ./...
	@go tool cover -func var/coverage.txt | awk '/^total/{print $$1 " " $$3}'

post-lint:
	@git diff --exit-code --quiet || (echo "There should not be any changes after the lint runs" && git status && exit 122;)

pipeline: full post-lint
