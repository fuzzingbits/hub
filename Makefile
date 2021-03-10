.PHONY: help docker build build-npm build-go lint lint-npm lint-go test test-go test-npm clean clean-full copy-config post-lint

SHELL=/bin/bash -o pipefail

.DEFAULT_GOAL := help
GO_PATH := $(shell go env GOPATH 2> /dev/null)
PATH := $(GO_PATH)/bin:$(PATH)

help: ## Display general help about this command
	@echo 'Makefile targets:'
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' Makefile \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/    \1 :: \3/p' \
	| column -t -c 1  -s '::'

docker:
	docker build -t fuzzingbits/hub:latest .

build: build-npm build-go ## Build the application

build-npm:
	npm install
	npm run build

build-go:
	@cd ; go get github.com/gobuffalo/packr/packr
	@go generate
	packr build -ldflags='-s -w' -o $(CURDIR)/var/hub .
	@ln -sf $(CURDIR)/var/hub $(GO_PATH)/bin/hub

lint: lint-npm lint-go ## Lint the application

lint-npm:
	npm install
	npm run lint

lint-go:
	@cd ; go get golang.org/x/lint/golint
	@cd ; go get golang.org/x/tools/cmd/goimports
	go get -d ./...
	go mod tidy
	gofmt -s -w .
	go vet ./...
	golint -set_exit_status=1 ./...
	goimports -w .

test: test-go test-npm ## Test the application

test-go:
	@mkdir -p var/
	@go test -race -cover -coverprofile  var/coverage.txt ./...
	@go tool cover -func var/coverage.txt | awk '/^total/{print $$1 " " $$3}'

test-npm:
	npm install
	npm run test

clean: ## Remove files listed in .gitignore (possibly with some exceptions)
	git clean -Xdff --exclude='!/.env'

clean-full:
	git clean -Xdff

copy-config: ## Copy missing config files into place
	[ -f /.env ] || cp /.env.dist /.env

post-lint:
	@git diff --exit-code --quiet || (echo 'There should not be any changes after the lint runs' && git status && exit 1;)
