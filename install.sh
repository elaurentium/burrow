#!/usr/bin/env bash
set -euo pipefail

REPO="elaurentium/burrow"
BINARY_NAME="b"
TMP_DIR="$(mktemp -d)"
INSTALL_DIR="/usr/local/bin"

UNAME_WIN="${USERNAME:-}"
UNAME_UNIX="${USER:-}"
if [[ -n "$UNAME_WIN" ]]; then
  USER_HOME_WIN="C:/Users/$UNAME_WIN"
else
  USER_HOME_WIN="C:/Users/${UNAME_UNIX:-User}"
fi
INSTALL_DIR_WIN="${USER_HOME_WIN}/.burrow"

detect_platform() {
  OS="$(uname -s)"
  ARCH="$(uname -m)"
  case "$OS" in
    Linux)   GOOS="linux" ;;
    Darwin)  GOOS="darwin" ;;
    MINGW*|MSYS*|CYGWIN*) GOOS="windows" ;;
    *) echo "OS not supported: $OS" >&2; exit 1 ;;
  esac
  case "$ARCH" in
    x86_64|amd64) GOARCH="amd64" ;;
    arm64|aarch64) GOARCH="arm64" ;;
    *) echo "Arch not suported: $ARCH" >&2; exit 1 ;;
  esac
}

latest_release_tag() {
  curl -s https://api.github.com/repos/$REPO/releases/latest | grep -oE '"tag_name":\s*"[^"]+"' | cut -d'"' -f4
}

ensure_dir() {
  local d="$1"
  if [[ "$GOOS" = "windows" ]]; then
    mkdir -p "$d"
  else
    sudo mkdir -p "$d"
  fi
}

main() {
  detect_platform
  TAG="${TAG_OVERRIDE:-$(latest_release_tag)}"
  if [ -z "${TAG:-}" ]; then
    echo "Could not get latest release. Set TAG_OVERRIDE." >&2
    exit 1
  fi

  FILE="burrow-${GOOS}-${GOARCH}"
  if [ "$GOOS" = "windows" ]; then
    FILE="${FILE}.exe"
  fi
  URL="https://github.com/${REPO}/releases/download/${TAG}/${FILE}"

  echo "Downloading ${URL}..."
  curl -fsSL "$URL" -o "${TMP_DIR}/${FILE}"

  if [ "$GOOS" != "windows" ]; then
    chmod +x "${TMP_DIR}/${FILE}"
  fi

  if [ "$GOOS" = "windows" ]; then
    ensure_dir "$INSTALL_DIR_WIN"
    mv "${TMP_DIR}/${FILE}" "${INSTALL_DIR_WIN}/${BINARY_NAME}.exe"
    echo "Installing in: ${INSTALL_DIR_WIN}\\${BINARY_NAME}.exe"
    echo "Version: ${TAG}"
    "${INSTALL_DIR_WIN}/${BINARY_NAME}.exe" --help || true
    echo
    echo "OBSERVATION: add ${INSTALL_DIR_WIN} on PATH to use 'b' on anywhere."
    echo "On PowerShell (temp):"
    echo '  $env:Path += ";'"${INSTALL_DIR_WIN}"'"'
  else
    ensure_dir "$INSTALL_DIR"
    sudo mv "${TMP_DIR}/${FILE}" "${INSTALL_DIR}/${BINARY_NAME}"
    echo "Install in: ${INSTALL_DIR}/${BINARY_NAME}"
    echo "Version: ${TAG}"
    "${INSTALL_DIR}/${BINARY_NAME}" --help || true
  fi
}

main "$@"