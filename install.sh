#!/bin/sh
# install.sh — Installer for the vkc CLI binary
#
# Usage:
#   curl -sSfL https://github.com/jeanmolossi/vibe-and-kalika-code/releases/latest/download/install.sh | sh
#   curl -sSfL https://github.com/jeanmolossi/vibe-and-kalika-code/releases/latest/download/install.sh | VERSION=1.2.3 sh
#   VERSION=1.2.3 sh install.sh
#
# Environment variables:
#   VERSION   — override the version to install (e.g. VERSION=1.2.3)
#   NO_COLOR  — set to any value to disable colored output

set -e

# ---------------------------------------------------------------------------
# Constants
# ---------------------------------------------------------------------------
REPO="jeanmolossi/vibe-and-kalika-code"
BINARY_NAME="vkc"
GITHUB_API="https://api.github.com/repos/${REPO}/releases/latest"
GITHUB_RELEASES="https://github.com/${REPO}/releases/download"

# ---------------------------------------------------------------------------
# Color helpers (respect NO_COLOR)
# ---------------------------------------------------------------------------
if [ -z "${NO_COLOR}" ] && [ -t 1 ]; then
  _BLUE='\033[0;34m'
  _GREEN='\033[0;32m'
  _RED='\033[0;31m'
  _RESET='\033[0m'
else
  _BLUE=''
  _GREEN=''
  _RED=''
  _RESET=''
fi

info()    { printf "${_BLUE}[info]${_RESET}  %s\n" "$*" >&2; }
success() { printf "${_GREEN}[ok]${_RESET}    %s\n" "$*" >&2; }
error()   { printf "${_RED}[error]${_RESET} %s\n" "$*" >&2; exit 1; }

# ---------------------------------------------------------------------------
# Detect OS
# ---------------------------------------------------------------------------
detect_os() {
  _os="$(uname -s 2>/dev/null | tr '[:upper:]' '[:lower:]')"
  case "${_os}" in
    linux*)
      echo "linux"
      ;;
    darwin*)
      echo "darwin"
      ;;
    mingw*|msys*|cygwin*|windows*)
      error "Windows is not supported by this installer. Please download the binary manually from: https://github.com/${REPO}/releases"
      ;;
    *)
      error "Unsupported OS: ${_os}. Please download the binary manually from: https://github.com/${REPO}/releases"
      ;;
  esac
}

# ---------------------------------------------------------------------------
# Detect architecture
# ---------------------------------------------------------------------------
detect_arch() {
  _arch="$(uname -m 2>/dev/null)"
  case "${_arch}" in
    x86_64|amd64)
      echo "amd64"
      ;;
    aarch64|arm64)
      echo "arm64"
      ;;
    *)
      error "Unsupported architecture: ${_arch}. Only amd64 and arm64 are supported."
      ;;
  esac
}

# ---------------------------------------------------------------------------
# Resolve version (env override or latest from GitHub API)
# ---------------------------------------------------------------------------
resolve_version() {
  if [ -n "${VERSION}" ]; then
    # Strip leading 'v' if provided (we add it ourselves when needed)
    echo "${VERSION#v}"
    return
  fi

  info "Fetching latest release version from GitHub API..."

  # Try curl first, then wget
  if command -v curl >/dev/null 2>&1; then
    _response="$(curl -sSfL "${GITHUB_API}" 2>/dev/null)"
  elif command -v wget >/dev/null 2>&1; then
    _response="$(wget -qO- "${GITHUB_API}" 2>/dev/null)"
  else
    error "Neither curl nor wget is available. Please install one and retry."
  fi

  # Parse tag_name from JSON (avoid jq dependency with grep/sed)
  _version="$(echo "${_response}" | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/' | sed 's/^v//')"

  if [ -z "${_version}" ]; then
    error "Could not determine latest version. Set VERSION env var to override (e.g. VERSION=1.2.3)."
  fi

  echo "${_version}"
}

# ---------------------------------------------------------------------------
# Download helper (curl or wget)
# ---------------------------------------------------------------------------
download() {
  _url="$1"
  _dest="$2"

  if command -v curl >/dev/null 2>&1; then
    curl -sSfL -o "${_dest}" "${_url}"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "${_dest}" "${_url}"
  else
    error "Neither curl nor wget is available. Please install one and retry."
  fi
}

# ---------------------------------------------------------------------------
# Verify SHA256 checksum
# ---------------------------------------------------------------------------
verify_checksum() {
  _archive="$1"        # path to archive file
  _checksums="$2"      # path to checksums file
  _filename="$3"       # expected filename inside checksums file

  _expected="$(grep "${_filename}" "${_checksums}" | awk '{print $1}')"

  if [ -z "${_expected}" ]; then
    error "Checksum entry for '${_filename}' not found in checksums file."
  fi

  if command -v sha256sum >/dev/null 2>&1; then
    # Linux: sha256sum
    _actual="$(sha256sum "${_archive}" | awk '{print $1}')"
  elif command -v shasum >/dev/null 2>&1; then
    # macOS: shasum -a 256
    _actual="$(shasum -a 256 "${_archive}" | awk '{print $1}')"
  else
    error "No SHA256 utility found (sha256sum or shasum). Cannot verify checksum."
  fi

  if [ "${_actual}" != "${_expected}" ]; then
    error "Checksum mismatch for '${_filename}'!\n  expected: ${_expected}\n  actual:   ${_actual}"
  fi
}

# ---------------------------------------------------------------------------
# Install shell completions automatically
# ---------------------------------------------------------------------------
install_completions() {
  # Resolve the installed binary path for generating completions.
  if command -v vkc >/dev/null 2>&1; then
    _vkc_bin="vkc"
  elif [ -x "/usr/local/bin/${BINARY_NAME}" ]; then
    _vkc_bin="/usr/local/bin/${BINARY_NAME}"
  elif [ -x "${HOME}/.local/bin/${BINARY_NAME}" ]; then
    _vkc_bin="${HOME}/.local/bin/${BINARY_NAME}"
  else
    info "Could not locate ${BINARY_NAME} binary to generate completions. Skipping."
    return
  fi

  _shell="${SHELL##*/}"

  case "${_shell}" in
    bash)
      _bash_comp_dir="${HOME}/.local/share/bash-completion/completions"
      mkdir -p "${_bash_comp_dir}"
      "${_vkc_bin}" completion bash > "${_bash_comp_dir}/${BINARY_NAME}" 2>/dev/null && \
        success "Bash completion installed to ${_bash_comp_dir}/${BINARY_NAME}" || \
        info "Could not install bash completion."
      ;;
    zsh)
      _zsh_comp_dir="${HOME}/.zsh/completions"
      mkdir -p "${_zsh_comp_dir}"
      "${_vkc_bin}" completion zsh > "${_zsh_comp_dir}/_${BINARY_NAME}" 2>/dev/null && \
        success "Zsh completion installed to ${_zsh_comp_dir}/_${BINARY_NAME}" || \
        info "Could not install zsh completion."

      # Ensure the completion dir is in fpath and compinit is called.
      _zshrc="${ZDOTDIR:-${HOME}}/.zshrc"
      if [ -f "${_zshrc}" ] && ! grep -q "/.zsh/completions" "${_zshrc}" 2>/dev/null; then
        printf '\n# vkc shell completions\nfpath=(%s $fpath)\nautoload -Uz compinit && compinit\n' "${_zsh_comp_dir}" >> "${_zshrc}"
        info "Added fpath entry to ${_zshrc}. Restart your shell or run: source ${_zshrc}"
      fi
      ;;
    fish)
      _fish_comp_dir="${HOME}/.config/fish/completions"
      mkdir -p "${_fish_comp_dir}"
      "${_vkc_bin}" completion fish > "${_fish_comp_dir}/${BINARY_NAME}.fish" 2>/dev/null && \
        success "Fish completion installed to ${_fish_comp_dir}/${BINARY_NAME}.fish" || \
        info "Could not install fish completion."
      ;;
    *)
      info "Shell '${_shell}' not recognised. Skipping automatic completion install."
      info "To enable completions manually, run: ${_vkc_bin} completion <shell>"
      ;;
  esac
}

# ---------------------------------------------------------------------------
# Install binary (try /usr/local/bin, fall back to ~/.local/bin)
# ---------------------------------------------------------------------------
install_binary() {
  _src="$1"
  _system_dir="/usr/local/bin"
  _user_dir="${HOME}/.local/bin"

  # Try system-wide install
  if [ -w "${_system_dir}" ]; then
    cp "${_src}" "${_system_dir}/${BINARY_NAME}"
    chmod +x "${_system_dir}/${BINARY_NAME}"
    success "Installed ${BINARY_NAME} to ${_system_dir}/${BINARY_NAME}"
    return
  fi

  # Try sudo if available
  if command -v sudo >/dev/null 2>&1 && sudo -n true 2>/dev/null; then
    sudo cp "${_src}" "${_system_dir}/${BINARY_NAME}"
    sudo chmod +x "${_system_dir}/${BINARY_NAME}"
    success "Installed ${BINARY_NAME} to ${_system_dir}/${BINARY_NAME} (via sudo)"
    return
  fi

  # Fall back to user-local bin
  mkdir -p "${_user_dir}"
  cp "${_src}" "${_user_dir}/${BINARY_NAME}"
  chmod +x "${_user_dir}/${BINARY_NAME}"
  success "Installed ${BINARY_NAME} to ${_user_dir}/${BINARY_NAME}"

  # Check if the user dir is in PATH and remind if not
  case ":${PATH}:" in
    *":${_user_dir}:"*)
      ;;
    *)
      printf "\n"
      info "NOTE: ${_user_dir} is not in your PATH."
      info "Add the following line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
      printf "  export PATH=\"\$HOME/.local/bin:\$PATH\"\n"
      printf "\n"
      ;;
  esac
}

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------
main() {
  OS="$(detect_os)"
  ARCH="$(detect_arch)"
  VERSION="$(resolve_version)"

  # Build filenames following goreleaser convention
  ARCHIVE_NAME="${BINARY_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
  CHECKSUM_NAME="${BINARY_NAME}_${VERSION}_checksums.txt"

  ARCHIVE_URL="${GITHUB_RELEASES}/v${VERSION}/${ARCHIVE_NAME}"
  CHECKSUM_URL="${GITHUB_RELEASES}/v${VERSION}/${CHECKSUM_NAME}"

  info "Installing ${BINARY_NAME} v${VERSION} (${OS}/${ARCH})..."

  # Create a temporary working directory; clean up on any exit
  TMP_DIR="$(mktemp -d)"
  trap 'rm -rf "${TMP_DIR}"' EXIT INT TERM

  # Download archive
  info "Downloading ${ARCHIVE_NAME}..."
  download "${ARCHIVE_URL}" "${TMP_DIR}/${ARCHIVE_NAME}"

  # Download checksums
  info "Downloading ${CHECKSUM_NAME}..."
  download "${CHECKSUM_URL}" "${TMP_DIR}/${CHECKSUM_NAME}"

  # Verify checksum
  info "Verifying SHA256 checksum..."
  verify_checksum \
    "${TMP_DIR}/${ARCHIVE_NAME}" \
    "${TMP_DIR}/${CHECKSUM_NAME}" \
    "${ARCHIVE_NAME}"
  success "Checksum verified."

  # Extract binary from archive
  info "Extracting binary..."
  tar -xzf "${TMP_DIR}/${ARCHIVE_NAME}" -C "${TMP_DIR}" "${BINARY_NAME}"

  if [ ! -f "${TMP_DIR}/${BINARY_NAME}" ]; then
    error "Binary '${BINARY_NAME}' not found in archive. The archive layout may have changed."
  fi

  # Install
  install_binary "${TMP_DIR}/${BINARY_NAME}"

  # Register shell completions automatically
  install_completions

  success "Done! Run '${BINARY_NAME} --help' to get started."
}

main
