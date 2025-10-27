#!/usr/bin/env bash
set -euo pipefail

REPO="elaurentium/burrow"
BINARY_NAME="b"
INSTALL_DIR="/usr/local/bin"
TMP_DIR="$(mktemp -d)"

ua_hdr=(-H "User-Agent: burrow-installer")
auth_hdr=()
if [ -n "${GITHUB_TOKEN:-}" ]; then
  auth_hdr=(-H "Authorization: token ${GITHUB_TOKEN}")
fi

detect_platform() {
  OS="$(uname -s)"
  ARCH="$(uname -m)"
  case "$OS" in
    Linux)  GOOS="linux" ;;
    Darwin) GOOS="darwin" ;;
    *) echo "Sistema operacional não suportado: $OS" >&2; exit 1 ;;
  esac
  case "$ARCH" in
    x86_64|amd64) GOARCH="amd64" ;;
    arm64|aarch64) GOARCH="arm64" ;;
    *) echo "Arquitetura não suportada: $ARCH" >&2; exit 1 ;;
  esac
}

latest_release_tag() {
  curl -fsSL "${ua_hdr[@]}" "${auth_hdr[@]}" \
    "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep -oE '"tag_name":\s*"[^"]+"' | cut -d'"' -f4
}

download_asset() {
  local url="$1" out="$2"
  echo "Baixando ${url}..."
  curl -fsSL "${ua_hdr[@]}" "${auth_hdr[@]}" -o "$out" "$url"
}

is_elf_or_macho() {
  if command -v file >/dev/null 2>&1; then
    file "$1" | grep -Eiq 'ELF|Mach-O'
  else
    # fallback básico
    local sig
    sig="$(head -c 4 "$1" | hexdump -v -e '/1 "%02X"' )"
    [[ "$sig" == "7F454C46" || "$sig" == "FEEDFACE" || "$sig" == "FEEDFACF" || "$sig" == "CAFEBABE" ]]
  fi
}

main() {
  detect_platform
  TAG="${TAG_OVERRIDE:-${VERSION:-$(latest_release_tag)}}"
  if [ -z "${TAG:-}" ]; then
    echo "Não foi possível obter a última release. Defina TAG_OVERRIDE ou VERSION." >&2
    exit 1
  fi

  FILE="burrow-${GOOS}-${GOARCH}"
  URL="https://github.com/${REPO}/releases/download/${TAG}/${FILE}"
  OUT="${TMP_DIR}/${FILE}"

  download_asset "$URL" "$OUT"

  if ! is_elf_or_macho "$OUT"; then
    echo "ERRO: asset baixado não é executável nativo (esperado ELF/Mach-O)." >&2
    echo "Dica: verifique se o release publica binário 'puro' gerado por 'go build'." >&2
    echo "file(1) diz: $(file "$OUT")" >&2 || true
    exit 1
  fi

  chmod +x "$OUT"
  echo "Instalando em ${INSTALL_DIR}/${BINARY_NAME} (sudo pode ser solicitado)..."
  sudo mv "$OUT" "${INSTALL_DIR}/${BINARY_NAME}"
  echo "Instalado em ${INSTALL_DIR}/${BINARY_NAME}"
  "${INSTALL_DIR}/${BINARY_NAME}" --help || true
}

main "$@"