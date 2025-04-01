#!/usr/bin/env bash

# Fail on errors and don't open cover file
set -e

git checkout go.mod
rm -rf go.sum
rm -rf vendor

GOPROXY=direct GOPRIVATE=github.com go mod tidy
go mod vendor

cd ./share/shallow_security
./build.sh
mv security.so ../../.
cd ../../

# Run unit tests with coverage
go test -v -coverpkg=./share/... -coverprofile=cover-report.html ./... --failfast

# Open the coverage report in a browser
go tool cover -html=cover-report.html
