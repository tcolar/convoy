#!/bin/bash

# Build linux x64 binary, to be used by docker container
rm -rf build/
mkdir -p ./build/linux_amd64/
binpath=./build/linux_amd64/convoy
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s" \
  -a -installsuffix cgo -o $binpath ./cmd/convoy/ || exit 1
chmod +x $binpath
