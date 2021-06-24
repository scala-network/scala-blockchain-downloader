.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

APP_NAME := scala-blockchain-downloader

default: build

clean:
	rm -rf bin/ && mkdir bin/
	mkdir bin/win-x64
	mkdir bin/linux-x64
	mkdir bin/macos-x64

build_windows:
	GOOS=windows \
	GOARCH=amd64 \
	go build -o bin/win-x64/${APP_NAME}.exe src/main.go

build_linux:
	GOOS=linux \
	GOARCH=amd64 \
	go build -o bin/linux-x64/${APP_NAME} src/main.go

build_macos:
	GOOS=darwin \
	GOARCH=amd64 \
	go build -o bin/macos-x64/${APP_NAME} src/main.go

build: build_linux \
	build_windows \
	build_macos
