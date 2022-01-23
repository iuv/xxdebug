#!/bin/bash
# linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s'
upx registry-hub
mv registry-hub build/linux/registry-hub
# mac
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags '-w -s'
upx registry-hub
mv registry-hub build/mac/registry-hub