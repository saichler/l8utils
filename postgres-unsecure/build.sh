#!/bin/bash
set -e 
echo "Building PostgreSQL Docker image..."
echo "  Output: saichler/unsecure-postgres:latest"

docker build \
    --no-cache --platform=linux/amd64 \
    --build-arg NAME="alpine:latest" \
    --build-arg HOST_UID="$(id -u)" \
    --build-arg HOST_GID="$(id -g)" \
    -t "saichler/unsecure-postgres:latest" \
    -f Dockerfile .

docker push saichler/unsecure-postgres:latest

