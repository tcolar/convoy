#!/bin/sh

# Build linux x64 binary, to be used by docker container
mkdir -p ./build/linux_amd64/
binpath=./build/linux_amd64/convoy
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $binpath || exit 1
chmod +x $binpath
