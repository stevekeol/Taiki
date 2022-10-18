#!/usr/bin/env bash

AUTHOR=stevekeol
NAME=Taiki

Taiki: clean build

## build: Builds application binary and stores it in `./bin/Taiki`
build:
	@echo "  >  \033[32mBuilding binary...\033[0m "
	go build -o ./bin/Taiki ./cmd/taiki/main.go

## clean: Clean the binary file
clean:
	rm -rf ./bin

## clean: Clean the database data generated after testing
cleanDB:
	rm -rf .data	

## install: install the Taiki binary in $GOPATH/bin
install: build
	mv ./bin/Taiki $(GOPATH)/bin/Taiki

## deps: Install missing dependencies. Runs `go mod tidy` internally.
deps:
	@echo "  >  \033[32mInstalling dependencies...\033[0m "
	go mod tidy

## TODO: test
## TODO: docker