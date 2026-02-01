.PHONY: build run all

build:
	@go build -v -ldflags="-X main.Version=$(VERSION) -s -w" -o ./bin/current-state ./cmd

run:
	@./bin/current-state

all: build run