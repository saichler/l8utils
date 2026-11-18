#!/bin/bash
set -e 
echo "Building PostgreSQL Docker image..."
echo "  Output: registry.layer-8.dev:44443/library/unsecure-postgres:latest"

docker build \
    --no-cache --platform=linux/amd64 \
    --build-arg NAME="alpine:latest" \
    --build-arg HOST_UID="$(id -u)" \
    --build-arg HOST_GID="$(id -g)" \
    -t "registry.layer-8.dev:44443/library/unsecure-postgres:latest" \
    -f Dockerfile .

docker push registry.layer-8.dev:44443/library/unsecure-postgres:latest

