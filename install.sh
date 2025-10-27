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
    *) echo "Sistema operacional não suportado: $OS" >&2; exit 1 ;;
  esac
  case "$ARCH" in
    x86_64|amd64) GOARCH="amd64" ;;
    arm64|aarch64) GOARCH="arm64" ;;
    *) echo "Arquitetura não suportada: $ARCH" >&2; exit 1 ;;
  esac
}

latest_release_tag() {
  curl -s https://api.github.com/repos/$REPO/releases/latest | grep -oE '"tag_name":\s*"[^"]+"' | cut -d'"' -f4
}

main() {
  detect_platform
  TAG="${TAG_OVERRIDE:-$(latest_release_tag)}"
  if [ -z "${TAG:-}" ]; then
    echo "Não foi possível obter a última release. Defina TAG_OVERRIDE" >&2
    exit 1
  fi

  FILE="burrow-${GOOS}-${GOARCH}"
  URL="https://github.com/${REPO}/releases/download/${TAG}/${FILE}"

  echo "Baixando ${URL}..."
  curl -fsSL "$URL" -o "${TMP_DIR}/${FILE}"

  chmod +x "${TMP_DIR}/${FILE}"
  sudo mv "${TMP_DIR}/${FILE}" "${INSTALL_DIR}/${BINARY_NAME}"

  echo "Instalado em ${INSTALL_DIR}/${BINARY_NAME}"
  echo "Versão: ${TAG}"
  "${INSTALL_DIR}/${BINARY_NAME}" --help || true
}

main "$@"