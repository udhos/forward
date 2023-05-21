#!/bin/bash

version=$(go run ./cmd/forward -version | awk '{ print $2 }' | awk -F= '{ print $2 }')

echo version=$version

docker build \
    --no-cache \
    -t udhos/forward:latest \
    -t udhos/forward:$version \
    -f docker/Dockerfile .
