#!/usr/bin/env bash
set -euo pipefail

REPO="elaurentium/burrow"
BINARY_NAME="b"
TMP_DIR="$(mktemp -d)"
INSTALL_DIR="/usr/local/bin"

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
    *) echo "Arch not supported: $ARCH" >&2; exit 1 ;;
  esac
}

latest_release_tag() {
  curl -s https://api.github.com/repos/$REPO/releases/latest | grep -oE '"tag_name":\s*"[^"]+"' | cut -d'"' -f4
}

main() {
  detect_platform
  TAG="${TAG_OVERRIDE:-$(latest_release_tag)}"
  if [ -z "${TAG:-}" ]; then
    echo "Could not get latest release. Set TAG_OVERRIDE" >&2
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
    mv "${TMP_DIR}/${FILE}" "${INSTALL_DIR}/${BINARY_NAME}.exe"
  else
    sudo mv "${TMP_DIR}/${FILE}" "${INSTALL_DIR}/${BINARY_NAME}"
  fi

  if [ "$GOOS" = "windows" ]; then
    echo "Installing ${INSTALL_DIR}/${BINARY_NAME}.exe"
  else
    echo "Installing ${INSTALL_DIR}/${BINARY_NAME}"
  fi
  echo "Versão: ${TAG}"
  if [ "$GOOS" = "windows" ]; then
    "${INSTALL_DIR}/${BINARY_NAME}.exe" --help || true
  else
    "${INSTALL_DIR}/${BINARY_NAME}" --help || true
  fi
}

main "$@"