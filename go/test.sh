#!/usr/bin/env bash
go mod vendor
# Fail on errors and don't open cover file
set -e

cd ./share/shallow_security
go build -buildmode=plugin -o security.so plugin.go ShallowSecurityProvider.go
mv security.so ../../.
cd ../../

# Run unit tests with coverage
go test -v -coverpkg=./share/... -coverprofile=cover-report.html ./... --failfast

# Open the coverage report in a browser
go tool cover -html=cover-report.html
