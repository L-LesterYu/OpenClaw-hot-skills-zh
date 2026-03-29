#!/usr/bin/env bash

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

LINT_VERSION="2.5.0"
LINT_PACKAGE="github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v${LINT_VERSION}"

if command -v golangci-lint >/dev/null 2>&1; then
  version_output="$(golangci-lint version 2>/dev/null || true)"
  if [[ "$version_output" == *"version ${LINT_VERSION} "* ]]; then
    exec golangci-lint run ./...
  fi
  if [[ -n "$version_output" ]]; then
    echo "golangci-lint version mismatch, want ${LINT_VERSION}, got: ${version_output}" >&2
    echo "Falling back to pinned ${LINT_PACKAGE} via go run..." >&2
  fi
fi

exec go run "${LINT_PACKAGE}" run ./...
