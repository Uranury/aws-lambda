#!/usr/bin/env bash
set -euo pipefail

go mod tidy
go test ./...

# provided.al2023 expects a Linux binary literally named "bootstrap".
# lambda.norpc drops the legacy RPC code path, which shrinks the binary.
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go
zip -j function.zip bootstrap
echo "Built function.zip (runtime: provided.al2023, architecture: arm64)"
