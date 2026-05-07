#!/usr/bin/env bash
set -euo pipefail

REPO="gestalten-sass/agent-credential-guard"
INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="guard"

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "Fehlt: $1" >&2
    exit 1
  }
}

need_cmd curl
need_cmd tar

arch="$(uname -m)"
case "$arch" in
  x86_64) asset_arch="amd64" ;;
  aarch64|arm64) asset_arch="arm64" ;;
  *)
    echo "Nicht unterstuetzte Architektur: $arch" >&2
    exit 1
    ;;
esac

api_url="https://api.github.com/repos/${REPO}/releases/latest"
release_json="$(curl -fsSL "$api_url")"
tag="$(printf '%s' "$release_json" | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n1)"

if [[ -z "$tag" ]]; then
  echo "Konnte kein Release-Tag finden. Bitte zuerst ein GitHub Release erstellen." >&2
  exit 1
fi

asset="guard-linux-${asset_arch}.tar.gz"
url="https://github.com/${REPO}/releases/download/${tag}/${asset}"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

curl -fL "$url" -o "$tmpdir/$asset"
tar -xzf "$tmpdir/$asset" -C "$tmpdir"

mkdir -p "$INSTALL_DIR"
install -m 0755 "$tmpdir/guard-linux-${asset_arch}" "$INSTALL_DIR/$BINARY_NAME"

echo "Installiert: $INSTALL_DIR/$BINARY_NAME"
if ! command -v guard >/dev/null 2>&1; then
  echo "Hinweis: $INSTALL_DIR ist evtl. nicht im PATH."
  echo 'Temporär: export PATH="$HOME/.local/bin:$PATH"'
fi
