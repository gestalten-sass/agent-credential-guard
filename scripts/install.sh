#!/usr/bin/env bash
set -euo pipefail

REPO="gestalten-sass/agent-credential-guard"
BINARY_NAME="guard"
USER_INSTALL_DIR="${HOME}/.local/bin"
SYSTEM_INSTALL_DIR="/usr/local/bin"

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "Fehlt: $1" >&2
    exit 1
  }
}

append_path_if_needed() {
  local rc_file="$1"
  local line='export PATH="$HOME/.local/bin:$PATH"'

  if [[ ! -f "$rc_file" ]]; then
    printf '%s\n' "$line" > "$rc_file"
    return
  fi

  if ! grep -Fq '$HOME/.local/bin' "$rc_file"; then
    printf '\n%s\n' "$line" >> "$rc_file"
  fi
}

maybe_install_global_hook() {
  local guard_bin="$1"

  if [[ "${GUARD_AUTO_HOOK:-}" == "1" ]]; then
    "$guard_bin" hook install --global
    "$guard_bin" hook status --global || true
    return
  fi

  if [[ ! -t 0 ]]; then
    echo "Hinweis: Kein interaktives Terminal. Globalen Hook manuell aktivieren mit:"
    echo "  guard hook install --global"
    echo "Oder direkt automatisch bei Installation:"
    echo '  curl -fsSL https://raw.githubusercontent.com/gestalten-sass/agent-credential-guard/master/scripts/install.sh | GUARD_AUTO_HOOK=1 bash'
    return
  fi

  printf "Globalen Git-Pre-Commit-Hook jetzt aktivieren? [Y/n] "
  read -r answer
  case "${answer:-Y}" in
    Y|y|"")
      "$guard_bin" hook install --global
      "$guard_bin" hook status --global || true
      ;;
    *)
      echo "Uebersprungen. Du kannst spaeter aktivieren mit: guard hook install --global"
      ;;
  esac
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

src_bin="$tmpdir/guard-linux-${asset_arch}"

if [[ -w "$SYSTEM_INSTALL_DIR" ]]; then
  install_dir="$SYSTEM_INSTALL_DIR"
  install -m 0755 "$src_bin" "$install_dir/$BINARY_NAME"
  echo "Installiert: $install_dir/$BINARY_NAME"
  maybe_install_global_hook "$install_dir/$BINARY_NAME"
  exit 0
fi

if command -v sudo >/dev/null 2>&1; then
  install_dir="$SYSTEM_INSTALL_DIR"
  if sudo install -m 0755 "$src_bin" "$install_dir/$BINARY_NAME"; then
    echo "Installiert: $install_dir/$BINARY_NAME"
    maybe_install_global_hook "$install_dir/$BINARY_NAME"
    exit 0
  fi
fi

install_dir="$USER_INSTALL_DIR"
mkdir -p "$install_dir"
install -m 0755 "$src_bin" "$install_dir/$BINARY_NAME"

shell_name="$(basename "${SHELL:-}")"
case "$shell_name" in
  bash) append_path_if_needed "$HOME/.bashrc" ;;
  zsh) append_path_if_needed "$HOME/.zshrc" ;;
  *)
    append_path_if_needed "$HOME/.profile"
    ;;
esac

echo "Installiert: $install_dir/$BINARY_NAME"
echo "PATH wurde dauerhaft erweitert (falls noetig)."
echo "Bitte Terminal neu oeffnen oder einmalig ausfuehren:"
echo '  export PATH="$HOME/.local/bin:$PATH"'
maybe_install_global_hook "$install_dir/$BINARY_NAME"
