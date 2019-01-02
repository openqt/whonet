.PHONY: all clean build install

.EXPORT_ALL_VARIABLES:

all: build
	@echo $$GOPATH


clean:  ## clean everything not in git!!
	go fmt -x ./...


test:
	go test -v ./...


build: ## make all tools


install:

