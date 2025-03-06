#!/usr/bin/env bash
go mod vendor
# Fail on errors and don't open cover file
set -e

# Run unit tests with coverage
go test -v -coverpkg=./share/... -coverprofile=cover-report.html ./... --failfast

# Open the coverage report in a browser
go tool cover -html=cover-report.html
