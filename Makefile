-include .env
export

SHELL := /bin/bash -O globstar
GO_PATH := $(shell go env GOPATH 2> /dev/null)
MODULE := $(shell awk '/^module/ {print $$2}' go.mod)
NAMESPACE := $(shell awk -F "/" '/^module/ {print $$(NF-1)}' go.mod)
PROJECT_NAME := $(shell awk -F "/" '/^module/ {print $$(NF)}' go.mod)
PATH := $(GO_PATH)/bin:$(PATH)

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

dev-go: install-hooks dev-docker-up ## Start a dev instance of the Go Server
	clear
	@go generate
	@DEV_PROXY_TO_NUXT=true DEV_CLEAR_EXISTING_DATA=true go run main.go

dev-ui: install-hooks ## Start a dev instance of the UI
	clear
	[ -d ./node_modules ] || npm install
	npm run dev

dev-docker-up: install-hooks ## Start the docker containers used for development
	docker-compose -f ops/hub_dev/docker-compose.yml up -d

dev-docker-down: install-hooks ## Remove the docker containers used for development
	docker-compose -f ops/hub_dev/docker-compose.yml down

build-go:
	@go install github.com/gobuffalo/packr/packr
	@go generate
	packr build -o $(CURDIR)/var/$(PROJECT_NAME)
	@ln -sf $(CURDIR)/var/$(PROJECT_NAME) $(GO_PATH)/bin/$(PROJECT_NAME)

build-ui:
	npm run build

lint-go:
	@cd ; go get golang.org/x/lint/golint
	@cd ; go get golang.org/x/tools/cmd/goimports
	go get -d ./...
	gofmt -s -w .
	go vet ./...
	golint -set_exit_status=1 ./...
	goimports -w .

lint-ui:
	[ -d ./node_modules ] || npm install
	npm run fmt

test-go:
	@mkdir -p var/
	@go test -race -cover -coverprofile  var/coverage.txt ./...
	@go tool cover -func var/coverage.txt | awk '/^total/{print $$1 " " $$3}'

post-lint:
	@git diff --exit-code --quiet || (echo "There should not be any changes after the lint runs" && git status && exit 122;)

install-hooks:
	echo $$PATH
	@cp ops/hooks/* .git/hooks/

pipeline: full post-lint

loc:
	@echo -n "Go:"
	@wc -l **/*.go | tail -n 1 | awk '{print " " $$1}'
	@echo -n "TypeScript:"
	@wc -l ui/**/*.ts | tail -n 1 | awk '{print " " $$1}'
	@echo -n "Vue:"
	@wc -l ui/**/*.vue | tail -n 1 | awk '{print " " $$1}'
	@echo -n "CSS:"
	@wc -l ui/**/*.css | tail -n 1 | awk '{print " " $$1}'
	@echo -n "Total:"
	@wc -l **/*.go ui/**/*.ts ui/**/*.vue ui/**/*.css | tail -n 1 | awk '{print " " $$1}'
