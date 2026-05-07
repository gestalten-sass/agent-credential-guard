#!/usr/bin/env bash
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$REPO_DIR"
go install ./cmd/guard

GOBIN="$(go env GOPATH)/bin"

echo "guard wurde installiert."
echo "Binary erwartet unter: ${GOBIN}/guard"
echo "Wenn noch nicht gesetzt, PATH erweitern:"
echo '  export PATH="$PATH:$(go env GOPATH)/bin"'
