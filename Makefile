.PHONY: help full full-go full-npm docker build build-npm build-go lint lint-npm lint-go test test-npm test-go clean clean-full copy-config projectl git-change-check

SHELL=/bin/bash -o pipefail

.DEFAULT_GOAL := help
GO_PATH := $(shell go env GOPATH 2> /dev/null)
PATH := $(GO_PATH)/bin:$(PATH)

help: ## Display general help about this command
	@echo 'Makefile targets:'
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' Makefile \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/    \1 :: \3/p' \
	| column -t -c 1  -s '::'

full: lint test build

full-go: lint-go test-go build-go

full-npm: lint-npm test-npm build-npm

docker:
	docker build -t fuzzingbits/hub:latest .

build: build-npm build-go ## Build the application

build-npm:
	npm install
	npm run build

build-go:
	@go generate
	go build -ldflags='-s -w' -o $(CURDIR)/var/hub .
	@ln -sf $(CURDIR)/var/hub $(GO_PATH)/bin/hub

lint: lint-npm lint-go ## Lint the application

lint-npm:
	npm install
	npm run lint

lint-go:
	@go install golang.org/x/lint/golint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	go get -d ./...
	go mod tidy
	gofmt -s -w .
	go vet ./...
	golint -set_exit_status=1 ./...
	goimports -w .

test: test-npm test-go ## Test the application

test-npm:
	npm install
	npm run test

test-go:
	@mkdir -p var/
	@go test -race -cover -coverprofile  var/coverage.txt ./...
	@go tool cover -func var/coverage.txt | awk '/^total/{print $$1 " " $$3}'

clean: ## Remove files listed in .gitignore (possibly with some exceptions)
	@git init 2> /dev/null
	git clean -Xdff --exclude='!/.env'

clean-full:
	@git init 2> /dev/null
	git clean -Xdff

copy-config: ## Copy missing config files into place
	[ -f /.env ] || cp /.env.dist /.env

projectl:
	@go install github.com/aaronellington/projectl@latest
	$(shell go env GOPATH)/bin/projectl

git-change-check:
	@git diff --exit-code --quiet || (echo 'There should not be any changes at this point' && git status && exit 1;)
