.PHONY: all sort test build install

.EXPORT_ALL_VARIABLES:
HTTPS_PROXY=socks5://127.0.0.1:1080
GO111MODULE=on
CGO_ENABLED=0


all: sort build


sort:  ## refine code and packages
	go fmt -x ./...
	@echo
	go mod vendor -v


test:
	go test -v ./...


build: ## make all tools
	go build -o whonet

install:

