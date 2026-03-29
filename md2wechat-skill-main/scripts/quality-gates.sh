#!/usr/bin/env bash

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "🔍 Running local/CI quality gates..."

echo "1) Check formatting"
unformatted="$(gofmt -l .)"
if [[ -n "$unformatted" ]]; then
  echo "Unformatted files:"
  printf '%s\n' "$unformatted"
  exit 1
fi

echo "2) Run go vet"
go vet ./...

echo "3) Run golangci-lint"
bash scripts/run-golangci-lint.sh

echo "4) Run tests"
CGO_ENABLED=1 go test -count=1 ./...

echo "5) Run release checks"
make release-check

echo "quality-gates: OK"
