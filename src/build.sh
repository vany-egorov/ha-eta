#!/bin/bash
VERSION="0.0.1"
GO_BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%S)
GOOS="linux"
GOARCH="amd64"
GO_LDFLAGS="-s -w -X main.buildDate=${GO_BUILD_DATE} -X \"main.version=${VERSION}\" -extldflags -static"

GOOS=${GOOS} GOARCH=${GOARCH} \
go build \
  -v \
  -ldflags "${GO_LDFLAGS}" -o ./ha-eta ./common.go ./main.go
