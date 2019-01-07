.PHONY: all sort test build install

.EXPORT_ALL_VARIABLES:
HTTPS_PROXY=socks5://127.0.0.1:1080
GO111MODULE=on
CGO_ENABLED=0

APPVER=v0.1.0
GITVER=$(shell git rev-parse --short HEAD)
GOVER=$(shell go version)
BUILDTIME=$(shell date +%FT%T%z)

all: build


sort:  ## refine code and packages
	go fmt -x ./...
	@echo
	go mod vendor -v


test:
	go test -v ./...


build: ## make all tools
	go build -v -ldflags "-X 'github.com/openqt/whonet/cmd.AppVersion=${APPVER}' -X 'github.com/openqt/whonet/cmd.GoVersion=${GOVER}' -X 'github.com/openqt/whonet/cmd.GitVersion=${GITVER}' -X 'github.com/openqt/whonet/cmd.BuildTime=${BUILDTIME}'"

install:

