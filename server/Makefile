##
# DeGrens panel
#
# @file
# @version 0.1

PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

run:
	go run cmds/panel/main.go -config ./config.yml

run-prod:
	go run cmds/panel/main.go -config ./config.prodconn.yml

run-redis:
	go run cmds/panel/main.go -config ./config.redis.yml

lint: ## Lint the files
	go vet ./...
	go fmt ./...

dep: ## Get the dependencies
	go mod download

dep-vendor: dep
	go mod vendor

build: dep ## Build the binary file
	go build -v cmds/panel/main.go

watch:
	air

# end
