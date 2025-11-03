#!/usr/bin/env bash

# Fail on errors and don't open cover file
set -e

rm -rf go.mod
rm -rf go.sum
rm -rf vendor

go mod init
GOPROXY=direct GOPRIVATE=github.com go mod tidy
go mod vendor

echo "Starting Tests"
# Run unit tests with coverage
go test -v -coverpkg=./utils/... -coverprofile=cover-report.html ./... --failfast

# Open the coverage report in a browser
go tool cover -html=cover-report.html
