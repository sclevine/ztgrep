#!/bin/bash

version=${1:-0.0.0}

GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$version" -o ztgrep-macos-amd64 ./cmd/ztgrep
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=$version" -o ztgrep-macos-arm64 ./cmd/ztgrep
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$version" -o ztgrep-linux-amd64 ./cmd/ztgrep
GOOS=linux GOARCH=arm64 go build -ldflags "-X main.Version=$version" -o ztgrep-linux-arm64 ./cmd/ztgrep
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=$version" -o ztgrep.exe ./cmd/ztgrep

docker build . --build-arg "version=$version" -t "sclevine/ztgrep:$version"