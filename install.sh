#!/bin/sh
# AMBIE CLI installer.
# Detects platform and architecture, downloads the latest release binary
# from https://github.com/ambie-ai/ambie-cli/releases/latest, and installs
# it to /usr/local/bin (or $AMBIE_INSTALL_DIR).
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/ambie-ai/ambie-cli/main/install.sh | sh
#
# Override install dir:
#   curl -fsSL ... | AMBIE_INSTALL_DIR=$HOME/.local/bin sh

set -eu

REPO="ambie-ai/ambie-cli"
INSTALL_DIR="${AMBIE_INSTALL_DIR:-/usr/local/bin}"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "unsupported arch: $ARCH" >&2; exit 1 ;;
esac
case "$OS" in
    linux|darwin) ;;
    *) echo "unsupported OS: $OS (use the pre-built Windows zip from https://github.com/$REPO/releases)" >&2; exit 1 ;;
esac

LATEST_URL="https://api.github.com/repos/$REPO/releases/latest"
VERSION=$(curl -fsSL "$LATEST_URL" | grep -m1 '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo "failed to determine latest version" >&2
    exit 1
fi

ASSET="ambie_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/v${VERSION}/${ASSET}"

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

echo "downloading $URL"
curl -fsSL -o "$TMP/$ASSET" "$URL"
tar -xzf "$TMP/$ASSET" -C "$TMP"

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP/ambie" "$INSTALL_DIR/ambie"
else
    echo "installing to $INSTALL_DIR (requires sudo)"
    sudo mv "$TMP/ambie" "$INSTALL_DIR/ambie"
fi

echo "installed: $($INSTALL_DIR/ambie version)"
echo
echo "Get an API key at https://ambie.ai/signup, then:"
echo "  export AMBIE_API_KEY=amb_live_..."
echo "  ambie transcribe meeting.mp3"
