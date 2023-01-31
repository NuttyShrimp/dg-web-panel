##
# DeGrens panel
#
# @file
# @version 0.1

PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

run:
	go run cmds/panel/main.go -config ./config.yml

test:
	go test ./... -v -coverprofile=coverage.out

view-cov:
	go tool cover -func=coverage.out

view-cov-report:
	go tool cover -html=coverage.out

lint: ## Lint the files
	go vet ./...
	go fmt ./...

dep: ## Get the dependencies
	go mod download

dep-vendor: dep
	go mod vendor

build: dep ## Build the binary file
	go build -v cmds/panel/main.go


# end