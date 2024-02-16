.PHONY: build
build:
	go build -o main.exe -v ./cmd/main.go 

.DEFAULT_GOAL := build